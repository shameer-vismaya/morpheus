module github.com/nginxinc/kubernetes-ingress

go 1.16

require (
	github.com/aws/aws-sdk-go-v2/config v1.3.0
	github.com/aws/aws-sdk-go-v2/service/marketplacemetering v1.3.1
	github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/emicklei/go-restful v2.15.0+incompatible // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golangci/golangci-lint v1.40.1
	github.com/google/go-cmp v0.5.8
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/nginxinc/nginx-plus-go-client v0.8.0
	github.com/nginxinc/nginx-prometheus-exporter v0.9.0
	github.com/prometheus/client_golang v1.10.0
	github.com/spiffe/go-spiffe v1.1.0
	k8s.io/api v0.25.5
	k8s.io/apimachinery v0.25.5
	k8s.io/client-go v0.21.1
	k8s.io/code-generator v0.21.1
	sigs.k8s.io/controller-tools v0.5.0
)
