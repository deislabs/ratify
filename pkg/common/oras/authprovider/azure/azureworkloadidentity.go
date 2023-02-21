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
	"time"

	provider "github.com/deislabs/ratify/pkg/common/oras/authprovider"
	"github.com/deislabs/ratify/pkg/utils"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerregistry/runtime/2019-08-15-preview/containerregistry"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type AzureWIProviderFactory struct{}
type azureWIAuthProvider struct {
	aadToken confidential.AuthResult
	tenantID string
	clientID string
}

type azureWIAuthProviderConf struct {
	Name     string `json:"name"`
	ClientID string `json:"clientID"`
}

const (
	azureWIAuthProviderName string = "azureWorkloadIdentity"
)

// init calls Register for our Azure Workload Identity provider
func init() {
	provider.Register(azureWIAuthProviderName, &AzureWIProviderFactory{})
}

// Create returns an AzureWIAuthProvider
func (s *AzureWIProviderFactory) Create(authProviderConfig provider.AuthProviderConfig) (provider.AuthProvider, error) {
	conf := azureWIAuthProviderConf{}
	authProviderConfigBytes, err := json.Marshal(authProviderConfig)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(authProviderConfigBytes, &conf); err != nil {
		return nil, fmt.Errorf("failed to parse auth provider configuration: %w", err)
	}

	tenant := os.Getenv("AZURE_TENANT_ID")

	if tenant == "" {
		return nil, fmt.Errorf("azure tenant id environment variable is empty")
	}

	clientID := os.Getenv("AZURE_CLIENT_ID")
	if clientID == "" {
		clientID = conf.ClientID
		if clientID == "" {
			return nil, fmt.Errorf("AZURE_CLIENT_ID environment variable is empty")
		}
	}

	// retrieve an AAD Access token
	token, err := utils.GetAADAccessToken(context.Background(), tenant, clientID, AADResource)
	if err != nil {
		return nil, err
	}

	return &azureWIAuthProvider{
		aadToken: token,
		tenantID: tenant,
		clientID: clientID,
	}, nil
}

// Enabled checks for non empty tenant ID and AAD access token
func (d *azureWIAuthProvider) Enabled(ctx context.Context) bool {
	if d.tenantID == "" || d.clientID == "" {
		return false
	}

	if d.aadToken.AccessToken == "" {
		return false
	}

	return true
}

// Provide returns the credentials for a specified artifact.
// Uses Azure Workload Identity to retrieve an AAD access token which can be
// exchanged for a valid ACR refresh token for login.
func (d *azureWIAuthProvider) Provide(ctx context.Context, artifact string) (provider.AuthConfig, error) {
	if !d.Enabled(ctx) {
		return provider.AuthConfig{}, fmt.Errorf("azure workload identity auth provider is not properly enabled")
	}
	// parse the artifact reference string to extract the registry host name
	artifactHostName, err := provider.GetRegistryHostName(artifact)
	if err != nil {
		return provider.AuthConfig{}, err
	}

	// need to refresh AAD token if it's expired
	if time.Now().Add(time.Minute * 5).After(d.aadToken.ExpiresOn) {
		newToken, err := utils.GetAADAccessToken(ctx, d.tenantID, d.clientID, AADResource)
		if err != nil {
			return provider.AuthConfig{}, errors.Wrap(err, "could not refresh AAD token")
		}
		d.aadToken = newToken
		logrus.Info("successfully refreshed AAD token")
	}

	// add protocol to generate complete URI
	serverUrl := "https://" + artifactHostName

	// create registry client and exchange AAD token for registry refresh token
	refreshTokenClient := containerregistry.NewRefreshTokensClient(serverUrl)
	rt, err := refreshTokenClient.GetFromExchange(context.Background(), "access_token", artifactHostName, d.tenantID, "", d.aadToken.AccessToken)
	if err != nil {
		return provider.AuthConfig{}, fmt.Errorf("failed to get refresh token for container registry - %w", err)
	}

	authConfig := provider.AuthConfig{
		Username:  dockerTokenLoginUsernameGUID,
		Password:  *rt.RefreshToken,
		Provider:  d,
		ExpiresOn: d.aadToken.ExpiresOn,
	}

	return authConfig, nil
}
