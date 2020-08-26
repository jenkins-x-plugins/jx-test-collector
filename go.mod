module github.com/jenkins-x/jx-test-collector

go 1.13

require (
	github.com/Azure/go-autorest/autorest v0.9.8 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.8.3 // indirect
	github.com/fatih/color v1.9.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/imdario/mergo v0.3.10 // indirect
	github.com/jenkins-x/jx-gitops v0.0.230
	github.com/jenkins-x/jx-helpers v1.0.45
	github.com/jenkins-x/jx-logging v0.0.11
	github.com/jenkins-x/jx-secret v0.0.92
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-envconfig v0.1.2
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/text v0.3.3 // indirect
	k8s.io/api v0.18.1
	k8s.io/apimachinery v0.18.1
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.1+incompatible

	github.com/banzaicloud/bank-vaults => github.com/banzaicloud/bank-vaults v0.0.0-20191212164220-b327d7f2b681

	github.com/banzaicloud/bank-vaults/pkg/sdk => github.com/banzaicloud/bank-vaults/pkg/sdk v0.0.0-20191212164220-b327d7f2b681

	k8s.io/api => k8s.io/api v0.17.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.2
	k8s.io/client-go => k8s.io/client-go v0.16.5

)
