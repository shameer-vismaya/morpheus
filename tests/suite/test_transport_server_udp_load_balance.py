import pytest
import re
import socket

from suite.resources_utils import (
    wait_before_test,
    get_ts_nginx_template_conf,
    scale_deployment
)
from suite.custom_resources_utils import (
    patch_ts,
    read_ts,
    delete_ts,
    create_ts_from_yaml,
)
from settings import TEST_DATA

@pytest.mark.ts
@pytest.mark.parametrize(
    "crd_ingress_controller, transport_server_setup",
    [
        (
            {
                "type": "complete",
                "extra_args":
                    [
                        "-global-configuration=nginx-ingress/nginx-configuration",
                        "-enable-leader-election=false"
                    ]
            },
            {"example": "transport-server-udp-load-balance"},
        )
    ],
    indirect=True,
)
class TestTransportServerUdpLoadBalance:

    def restore_ts(self, kube_apis, transport_server_setup) -> None:
        """
        Function to revert a TransportServer resource to a valid state.
        """
        patch_src = f"{TEST_DATA}/transport-server-udp-load-balance/standard/transport-server.yaml"
        patch_ts(
            kube_apis.custom_objects,
            transport_server_setup.name,
            patch_src,
            transport_server_setup.namespace,
        )

    def test_number_of_replicas(
        self, kube_apis, crd_ingress_controller, transport_server_setup, ingress_controller_prerequisites
    ):
        """
        The load balancing of UDP should result in 4 servers to match the 4 replicas of a service.
        """
        original = scale_deployment(kube_apis.apps_v1_api, "udp-service", transport_server_setup.namespace, 4)
        wait_before_test()

        result_conf = get_ts_nginx_template_conf(
            kube_apis.v1,
            transport_server_setup.namespace,
            transport_server_setup.name,
            transport_server_setup.ingress_pod_name,
            ingress_controller_prerequisites.namespace
        )

        print(result_conf)

        pattern = 'server .*;'
        num_servers = len(re.findall(pattern, result_conf))

        assert num_servers is 4

        scale_deployment(kube_apis.apps_v1_api, "udp-service", transport_server_setup.namespace, original)
        wait_before_test()

    def test_udp_request_load_balanced(
            self, kube_apis, crd_ingress_controller, transport_server_setup, ingress_controller_prerequisites
    ):
        """
        Requests to the load balanced UDP service should result in responses from 3 different endpoints.
        """
        wait_before_test()
        port = transport_server_setup.public_endpoint.udp_server_port
        host = transport_server_setup.public_endpoint.public_ip

        print(f"sending udp requests to: {host}:{port}")

        endpoints = {}
        for i in range(20):
            client = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, 0)
            client.sendto("ping".encode('utf-8'), (host, port))
            data, address = client.recvfrom(4096)
            endpoint = data.decode()
            print(f' req number {i}; response: {endpoint}')
            if endpoint not in endpoints:
                endpoints[endpoint] = 1
            else:
                endpoints[endpoint] = endpoints[endpoint] + 1
            client.close()

        assert len(endpoints) is 3

        result_conf = get_ts_nginx_template_conf(
            kube_apis.v1,
            transport_server_setup.namespace,
            transport_server_setup.name,
            transport_server_setup.ingress_pod_name,
            ingress_controller_prerequisites.namespace
        )

        pattern = 'server .*;'
        servers = re.findall(pattern, result_conf)
        for key in endpoints.keys():
            found = False
            for server in servers:
                if key in server:
                    found = True
            assert found

    def test_udp_request_load_balanced_multiple(
            self, kube_apis, crd_ingress_controller, transport_server_setup
    ):
        """
        Requests to the load balanced UDP service should result in responses from 3 different endpoints.
        """
        port = transport_server_setup.public_endpoint.udp_server_port
        host = transport_server_setup.public_endpoint.public_ip

        # Step 1, confirm load balancing is working.
        print(f"sending udp requests to: {host}:{port}")
        client = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, 0)
        client.sendto("ping".encode('utf-8'), (host, port))
        data, address = client.recvfrom(4096)
        endpoint = data.decode()
        print(f'response: {endpoint}')
        client.close()

        # Step 2, add a second TransportServer with the same port and confirm te collision
        transport_server_file = f"{TEST_DATA}/transport-server-udp-load-balance/second-transport-server.yaml"
        ts_resource = create_ts_from_yaml(
            kube_apis.custom_objects, transport_server_file, transport_server_setup.namespace
        )
        wait_before_test()

        second_ts_name = ts_resource['metadata']['name']
        response = read_ts(
            kube_apis.custom_objects,
            transport_server_setup.namespace,
            second_ts_name,
        )
        assert (
                response["status"]
                and response["status"]["reason"] == "Rejected"
                and response["status"]["state"] == "Warning"
                and response["status"]["message"] == "Listener udp-server is taken by another resource"
        )

        # Step 3, remove the default TransportServer with the same port
        delete_ts(kube_apis.custom_objects, transport_server_setup.resource, transport_server_setup.namespace)

        wait_before_test()
        response = read_ts(
            kube_apis.custom_objects,
            transport_server_setup.namespace,
            second_ts_name,
        )
        assert (
                response["status"]
                and response["status"]["reason"] == "AddedOrUpdated"
                and response["status"]["state"] == "Valid"
        )

        # Step 4, confirm load balancing is still working.
        client = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, 0)
        client.sendto("ping".encode('utf-8'), (host, port))
        data, address = client.recvfrom(4096)
        endpoint = data.decode()
        print(f'response: {endpoint}')
        client.close()
        assert endpoint is not ""

        # cleanup
        delete_ts(kube_apis.custom_objects, ts_resource, transport_server_setup.namespace)
        transport_server_file = f"{TEST_DATA}/transport-server-udp-load-balance/standard/transport-server.yaml"
        create_ts_from_yaml(
            kube_apis.custom_objects, transport_server_file, transport_server_setup.namespace
        )
        wait_before_test()

    @pytest.mark.parametrize("file", ["wrong-port-transport-server.yaml", "missing-service-transport-server.yaml"])
    def test_udp_request_fails(
            self, kube_apis, crd_ingress_controller, transport_server_setup, file
    ):
        patch_src = f"{TEST_DATA}/transport-server-udp-load-balance/{file}"
        patch_ts(
            kube_apis.custom_objects,
            transport_server_setup.name,
            patch_src,
            transport_server_setup.namespace,
        )

        wait_before_test()

        port = transport_server_setup.public_endpoint.udp_server_port
        host = transport_server_setup.public_endpoint.public_ip

        print(f"sending udp requests to: {host}:{port}")
        for i in range(3):
            client = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, 0)
            client.settimeout(2)
            client.sendto("ping".encode('utf-8'), (host, port))
            try:
                client.recvfrom(4096)
                # it should timeout
                print(f"incorrect config from {file} should have resulted in an error")
                assert False
            except socket.timeout:
                print("successfully timed out")
            client.close()

        self.restore_ts(kube_apis, transport_server_setup)
