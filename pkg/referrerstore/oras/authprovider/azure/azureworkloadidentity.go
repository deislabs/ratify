/*
Copyright The Ratify Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/containerd/containerd/reference"
	provider "github.com/deislabs/ratify/pkg/referrerstore/oras/authprovider"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/klog/v2"
)

type AzureWIProviderFactory struct{}
type azureWIAuthProvider struct{}

type azureWIAuthProviderConf struct {
	Name string `json:"name"`
}

const (
	azureWIAuthProviderName      string = "azure-wi"
	dockerTokenLoginUsernameGUID string = "00000000-0000-0000-0000-000000000000"
)

// init calls Register for our default provider, which simply reads the .dockercfg file.
func init() {
	provider.Register(azureWIAuthProviderName, &AzureWIProviderFactory{})
	logrus.Info("Azure-WI provider registered")
}

// Create returns an empty defaultAuthProvider instance if the AuthProviderConfig is nil.
// Otherwise it returns the defaultAuthProvider with configPath set
func (s *AzureWIProviderFactory) Create(authProviderConfig provider.AuthProviderConfig) (provider.AuthProvider, error) {
	conf := azureWIAuthProviderConf{}
	authProviderConfigBytes, err := json.Marshal(authProviderConfig)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(authProviderConfigBytes, &conf); err != nil {
		return nil, fmt.Errorf("failed to parse auth provider configuration: %v", err)
	}
	logrus.Info("Azure-WI provider created succesfully")
	return &azureWIAuthProvider{}, nil
}

// Enabled always returns true for defaultAuthProvider
func (d *azureWIAuthProvider) Enabled() bool {
	return true
}

// Provide reads docker config file and returns corresponding credentials from file if exists
func (d *azureWIAuthProvider) Provide(artifact string) (provider.AuthConfig, error) {
	tenantID := os.Getenv("AZURE_TENANT_ID")

	parsedSpec, err := reference.Parse(artifact)
	if err != nil {
		return provider.AuthConfig{}, err
	}

	artifactHostName := parsedSpec.Hostname()

	aadToken, err := getAADAccessToken(tenantID, "https://management.azure.com/")
	if err != nil {
		return provider.AuthConfig{}, err
	}

	directive, err := receiveChallengeFromLoginServer(artifactHostName, "https")
	if err != nil {
		klog.Errorf("failed to receive challenge: %s", err)
		return provider.AuthConfig{}, err
	}

	refreshToken, err := performTokenExchange(artifactHostName, directive, tenantID, aadToken)
	if err != nil {
		return provider.AuthConfig{}, err
	}

	authConfig := provider.AuthConfig{
		Username: dockerTokenLoginUsernameGUID,
		Password: refreshToken,
		Provider: d,
	}
	return authConfig, nil
}

// Source: https://github.com/Azure/azure-workload-identity/blob/d126293e3c7c669378b225ad1b1f29cf6af4e56d/examples/msal-go/token_credential.go#L25
func getAADAccessToken(tenantID, resource string) (string, error) {
	// Azure AD Workload Identity webhook will inject the following env vars
	// 	AZURE_CLIENT_ID with the clientID set in the service account annotation
	// 	AZURE_TENANT_ID with the tenantID set in the service account annotation. If not defined, then
	// 	the tenantID provided via azure-wi-webhook-config for the webhook will be used.
	// 	AZURE_FEDERATED_TOKEN_FILE is the service account token path
	// 	AZURE_AUTHORITY_HOST is the AAD authority hostname
	clientID := os.Getenv("AZURE_CLIENT_ID")
	tokenFilePath := os.Getenv("AZURE_FEDERATED_TOKEN_FILE")
	authorityHost := os.Getenv("AZURE_AUTHORITY_HOST")

	// read the service account token from the filesystem
	signedAssertion, err := readJWTFromFS(tokenFilePath)
	if err != nil {
		klog.ErrorS(err, "failed to read the service account token from the filesystem")
		return "", errors.Wrap(err, "failed to read service account token")
	}
	cred, err := confidential.NewCredFromAssertion(signedAssertion)
	if err != nil {
		klog.ErrorS(err, "failed to create credential from signed assertion")
		return "", errors.Wrap(err, "failed to create confidential creds")
	}

	// create the confidential client to request an AAD token
	confidentialClientApp, err := confidential.New(
		clientID,
		cred,
		confidential.WithAuthority(fmt.Sprintf("%s%s/oauth2/token", authorityHost, tenantID)))
	if err != nil {
		klog.ErrorS(err, "failed to create confidential client")
		return "", errors.Wrap(err, "failed to create confidential client app")
	}

	// .default needs to be added to the scope
	if !strings.HasSuffix(resource, ".default") {
		resource += "/.default"
	}

	result, err := confidentialClientApp.AcquireTokenByCredential(context.Background(), []string{resource})
	if err != nil {
		klog.ErrorS(err, "failed to acquire token")
		return "", errors.Wrap(err, "failed to acquire token")
	}

	return result.AccessToken, nil
}

// readJWTFromFS reads the jwt from file system
// Source: https://github.com/Azure/azure-workload-identity/blob/d126293e3c7c669378b225ad1b1f29cf6af4e56d/examples/msal-go/token_credential.go#L88
func readJWTFromFS(tokenFilePath string) (string, error) {
	token, err := os.ReadFile(tokenFilePath)
	if err != nil {
		return "", err
	}
	return string(token), nil
}
