module github.com/jenkins-x/jx-test-collector

go 1.15

require (
	github.com/Azure/go-autorest/autorest v0.11.17 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.10 // indirect
	github.com/fatih/color v1.10.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-errors/errors v1.0.2 // indirect
	github.com/go-openapi/spec v0.19.15 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.1.5 // indirect
	github.com/googleapis/gnostic v0.4.2 // indirect
	github.com/jenkins-x/jx-api/v4 v4.0.21 // indirect
	github.com/jenkins-x/jx-gitops v0.0.525
	github.com/jenkins-x/jx-helpers/v3 v3.0.62
	github.com/jenkins-x/jx-secret v0.0.208
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-envconfig v0.1.2
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	k8s.io/api v0.19.7
	k8s.io/apimachinery v0.19.7
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/kube-openapi v0.0.0-20200923155610-8b5066479488 // indirect
	sigs.k8s.io/yaml v1.2.0
)

replace k8s.io/client-go => k8s.io/client-go v0.19.2
