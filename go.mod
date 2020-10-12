module github.com/jenkins-x/jx-test-collector

go 1.15

require (
	github.com/Azure/go-autorest/autorest v0.9.8 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.8.3 // indirect
	github.com/fatih/color v1.9.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/google/uuid v1.1.2 // indirect
	github.com/googleapis/gnostic v0.4.2 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/jenkins-x/jx-gitops v0.0.382
	github.com/jenkins-x/jx-helpers/v3 v3.0.6
	github.com/jenkins-x/jx-secret v0.0.170
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-envconfig v0.1.2
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.6.1
	google.golang.org/appengine v1.6.6 // indirect
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	sigs.k8s.io/yaml v1.2.0
)

replace k8s.io/client-go => k8s.io/client-go v0.19.2
