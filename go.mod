module github.com/ToucanSoftware/cloudship-operator

go 1.15

require (
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/go-logr/logr v0.3.0
	github.com/gofrs/flock v0.8.0
	github.com/google/martian v2.1.0+incompatible
	github.com/onsi/ginkgo v1.15.0
	github.com/onsi/gomega v1.10.5
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/operator-framework/operator-lib v0.4.0
	github.com/operator-framework/operator-sdk v1.5.0
	github.com/pkg/errors v0.9.1
	gopkg.in/yaml.v2 v2.3.0
	helm.sh/helm/v3 v3.5.0
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/cli-runtime v0.20.2
	k8s.io/client-go v0.20.2
	k8s.io/helm v2.17.0+incompatible
	sigs.k8s.io/controller-runtime v0.8.2
)
