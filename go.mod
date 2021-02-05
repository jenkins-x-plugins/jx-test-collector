module github.com/jenkins-x/jx-test-collector

go 1.15

require (
	github.com/Azure/go-autorest/autorest v0.11.17 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.11 // indirect
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.17+incompatible // indirect
	github.com/docker/cli v0.0.0-20200210162036-a4bedce16568 // indirect
	github.com/fatih/color v1.10.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-errors/errors v1.1.1 // indirect
	github.com/go-openapi/spec v0.20.2 // indirect
	github.com/go-openapi/swag v0.19.13 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/lint v0.0.0-20180702182130-06c8688daad7 // indirect
	github.com/google/go-github/v29 v29.0.3 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/jenkins-x/jx-api/v4 v4.0.24 // indirect
	github.com/jenkins-x/jx-gitops v0.1.2
	github.com/jenkins-x/jx-helpers/v3 v3.0.75
	github.com/jenkins-x/jx-secret v0.0.228
	github.com/jenkins-x/lighthouse v0.0.923 // indirect
	github.com/klauspost/cpuid v1.2.2 // indirect
	github.com/knative/build v0.1.2 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/nats-io/gnatsd v1.4.1 // indirect
	github.com/nats-io/go-nats v1.7.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-envconfig v0.3.2 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.7.0
	github.com/tsenart/vegeta v12.7.1-0.20190725001342-b5f4fca92137+incompatible // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/kube-openapi v0.0.0-20210113233702-8566a335510f // indirect
	knative.dev/test-infra v0.0.0-20200630141629-15f40fe97047 // indirect
	sigs.k8s.io/kustomize/kyaml v0.10.6 // indirect
	sigs.k8s.io/structured-merge-diff v1.0.1 // indirect
	sigs.k8s.io/testing_frameworks v0.1.1 // indirect
	sigs.k8s.io/yaml v1.2.0
)

replace k8s.io/client-go => k8s.io/client-go v0.19.2
