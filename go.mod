module github.com/deislabs/ratify

go 1.16

require (
	github.com/Azure/azure-sdk-for-go v64.1.0+incompatible
	github.com/AzureAD/microsoft-authentication-library-for-go v0.4.0
	github.com/docker/cli v20.10.16+incompatible
	github.com/docker/distribution v2.8.1+incompatible
	github.com/google/go-containerregistry v0.8.1-0.20220125170349-50dfc2733d10
	github.com/gorilla/mux v1.8.0
	github.com/notaryproject/notation-go-lib v1.0.0-alpha-1
	github.com/open-policy-agent/frameworks/constraint v0.0.0-20220527234808-13b0f3dbe9f0
	github.com/opencontainers/go-digest v1.0.0
	github.com/opencontainers/image-spec v1.0.2
	github.com/oras-project/artifacts-spec v1.0.0-rc.1
	github.com/pkg/errors v0.9.1
	github.com/sigstore/cosign v1.5.2
	github.com/sigstore/sigstore v1.1.1-0.20220130134424-bae9b66b8442
	github.com/sirupsen/logrus v1.8.1
	github.com/spdx/tools-golang v0.2.0
	github.com/spf13/cobra v1.4.0
	github.com/xlab/treeprint v1.1.0
	k8s.io/api v0.24.1
	k8s.io/apimachinery v0.24.1
	k8s.io/client-go v0.24.1
	oras.land/oras-go/v2 v2.0.0-20220630033939-f37492936f3e
)

replace github.com/open-policy-agent/opa => github.com/open-policy-agent/opa v0.40.0
