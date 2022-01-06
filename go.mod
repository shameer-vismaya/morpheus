module github.com/nginxinc/kubernetes-ingress

go 1.16

require (
	github.com/aws/aws-sdk-go-v2/config v1.3.0
	github.com/aws/aws-sdk-go-v2/service/marketplacemetering v1.3.1
	github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/emicklei/go-restful v2.15.0+incompatible // indirect
	github.com/golang/glog v1.0.0
	github.com/golangci/golangci-lint v1.40.1
	github.com/google/go-cmp v0.5.6
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/nginxinc/nginx-plus-go-client v0.8.0
	github.com/nginxinc/nginx-prometheus-exporter v0.9.0
	github.com/prometheus/client_golang v1.11.0
	github.com/spiffe/go-spiffe v1.1.0
	github.com/stretchr/objx v0.2.0 // indirect
	k8s.io/api v0.23.0
	k8s.io/apimachinery v0.23.0
	k8s.io/client-go v0.23.0
	k8s.io/code-generator v0.23.0
	sigs.k8s.io/controller-tools v0.8.0
)
