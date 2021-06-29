module github.com/jenkins-x/jx-test-collector

go 1.15

replace (
	// helm dependencies
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible

	github.com/sethvargo/go-envconfig => github.com/sethvargo/go-envconfig v0.1.2

	k8s.io/api => k8s.io/api v0.20.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.6
	k8s.io/client-go => k8s.io/client-go v0.20.6
)

require (
	github.com/fatih/color v1.12.0
	github.com/jenkins-x-plugins/jx-gitops v0.2.105
	github.com/jenkins-x-plugins/jx-secret v0.1.40
	github.com/jenkins-x/jx-helpers/v3 v3.0.119
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-envconfig v0.3.2
	github.com/sirupsen/logrus v1.8.1
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	sigs.k8s.io/yaml v1.2.0
)
