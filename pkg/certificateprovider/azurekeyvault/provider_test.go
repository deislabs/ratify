package azurekeyvault

import (
	"context"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/deislabs/ratify/pkg/certificateprovider/azurekeyvault/types"
	"github.com/stretchr/testify/assert"
)

func TestGetVaultURL(t *testing.T) {
	testEnvs := []string{"", "AZUREPUBLICCLOUD", "AZURECHINACLOUD", "AZUREGERMANCLOUD", "AZUREUSGOVERNMENTCLOUD"}
	vaultDNSSuffix := []string{"vault.azure.net", "vault.azure.net", "vault.azure.cn", "vault.microsoftazure.de", "vault.usgovcloudapi.net"}

	cases := []struct {
		desc        string
		vaultName   string
		expectedErr bool
	}{
		{
			desc:        "vault name > 24",
			vaultName:   "longkeyvaultnamewhichisnotvalid",
			expectedErr: true,
		},
		{
			desc:        "vault name < 3",
			vaultName:   "kv",
			expectedErr: true,
		},
		{
			desc:        "vault name contains non alpha-numeric chars",
			vaultName:   "kv_test",
			expectedErr: true,
		},
		{
			desc:        "valid vault name in public cloud",
			vaultName:   "testkv",
			expectedErr: false,
		},
	}

	for i, tc := range cases {
		t.Log(i, tc.desc)
		for idx := range testEnvs {
			azCloudEnv, err := ParseAzureEnvironment(testEnvs[idx])
			if err != nil {
				t.Fatalf("Error parsing cloud environment %v", err)
			}

			vaultURL, err := getVaultURL(tc.vaultName, azCloudEnv.KeyVaultDNSSuffix)
			if tc.expectedErr && err == nil || !tc.expectedErr && err != nil {
				t.Fatalf("expected error: %v, got error: %v", tc.expectedErr, err)
			}
			expectedURL := "https://" + tc.vaultName + "." + vaultDNSSuffix[idx] + "/"
			if !tc.expectedErr && expectedURL != *vaultURL {
				t.Fatalf("expected vault url: %s, got: %s", expectedURL, *vaultURL)
			}
		}
	}
}

func TestParseAzureEnvironment(t *testing.T) {
	envNamesArray := []string{"AZURECHINACLOUD", "AZUREGERMANCLOUD", "AZUREPUBLICCLOUD", "AZUREUSGOVERNMENTCLOUD", ""}
	for _, envName := range envNamesArray {
		azureEnv, err := ParseAzureEnvironment(envName)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if strings.EqualFold(envName, "") && !strings.EqualFold(azureEnv.Name, "AZUREPUBLICCLOUD") {
			t.Fatalf("string doesn't match, expected AZUREPUBLICCLOUD, got %s", azureEnv.Name)
		} else if !strings.EqualFold(envName, "") && !strings.EqualFold(envName, azureEnv.Name) {
			t.Fatalf("string doesn't match, expected %s, got %s", envName, azureEnv.Name)
		}
	}

	wrongEnvName := "AZUREWRONGCLOUD"
	_, err := ParseAzureEnvironment(wrongEnvName)
	if err == nil {
		t.Fatalf("expected error for wrong azure environment name")
	}
}

func TestGetLatestNKeyVaultObjects(t *testing.T) {
	now := time.Now()

	cases := []struct {
		desc            string
		kvObject        types.KeyVaultCertificate
		versions        types.KeyVaultObjectVersionList
		expectedObjects []types.KeyVaultCertificate
	}{
		{
			desc: "filename is name/index when no alias provided",
			kvObject: types.KeyVaultCertificate{
				CertificateName:           "cert1",
				CertificateVersion:        "latest",
				CertificateVersionHistory: 5,
			},
			versions: types.KeyVaultObjectVersionList{
				types.KeyVaultObjectVersion{
					Version: "a",
					Created: now.Add(time.Hour * 10),
				},
				types.KeyVaultObjectVersion{
					Version: "b",
					Created: now.Add(time.Hour * 9),
				},
				types.KeyVaultObjectVersion{
					Version: "c",
					Created: now.Add(time.Hour * 8),
				},
				types.KeyVaultObjectVersion{
					Version: "d",
					Created: now.Add(time.Hour * 7),
				},
				types.KeyVaultObjectVersion{
					Version: "e",
					Created: now.Add(time.Hour * 6),
				},
			},
			expectedObjects: []types.KeyVaultCertificate{
				{
					CertificateName:           "cert1",
					CertificateAlias:          filepath.Join("cert1", "0"),
					CertificateVersion:        "a",
					CertificateVersionHistory: 5,
				},
				{
					CertificateName:           "cert1",
					CertificateAlias:          filepath.Join("cert1", "1"),
					CertificateVersion:        "b",
					CertificateVersionHistory: 5,
				},
				{
					CertificateName:           "cert1",
					CertificateAlias:          filepath.Join("cert1", "2"),
					CertificateVersion:        "c",
					CertificateVersionHistory: 5,
				},
				{
					CertificateName:           "cert1",
					CertificateAlias:          filepath.Join("cert1", "3"),
					CertificateVersion:        "d",
					CertificateVersionHistory: 5,
				},
				{
					CertificateName:           "cert1",
					CertificateAlias:          filepath.Join("cert1", "4"),
					CertificateVersion:        "e",
					CertificateVersionHistory: 5,
				},
			},
		},
		{
			desc: "sorts versions by descending created date",
			kvObject: types.KeyVaultCertificate{
				CertificateName:           "cert1",
				CertificateVersion:        "latest",
				CertificateVersionHistory: 5,
			},
			versions: types.KeyVaultObjectVersionList{
				types.KeyVaultObjectVersion{
					Version: "c",
					Created: now.Add(time.Hour * 8),
				},
				types.KeyVaultObjectVersion{
					Version: "e",
					Created: now.Add(time.Hour * 6),
				},
				types.KeyVaultObjectVersion{
					Version: "b",
					Created: now.Add(time.Hour * 9),
				},
				types.KeyVaultObjectVersion{
					Version: "a",
					Created: now.Add(time.Hour * 10),
				},
				types.KeyVaultObjectVersion{
					Version: "d",
					Created: now.Add(time.Hour * 7),
				},
			},
			expectedObjects: []types.KeyVaultCertificate{
				{
					CertificateName:           "cert1",
					CertificateAlias:          filepath.Join("cert1", "0"),
					CertificateVersion:        "a",
					CertificateVersionHistory: 5,
				},
				{
					CertificateName:           "cert1",
					CertificateAlias:          filepath.Join("cert1", "1"),
					CertificateVersion:        "b",
					CertificateVersionHistory: 5,
				},
				{
					CertificateName:           "cert1",
					CertificateAlias:          filepath.Join("cert1", "2"),
					CertificateVersion:        "c",
					CertificateVersionHistory: 5,
				},
				{
					CertificateName:           "cert1",
					CertificateAlias:          filepath.Join("cert1", "3"),
					CertificateVersion:        "d",
					CertificateVersionHistory: 5,
				},
				{
					CertificateName:           "cert1",
					CertificateAlias:          filepath.Join("cert1", "4"),
					CertificateVersion:        "e",
					CertificateVersionHistory: 5,
				},
			},
		},
		{
			desc: "starts with latest version when no version specified",
			kvObject: types.KeyVaultCertificate{
				CertificateName:           "cert1",
				CertificateVersionHistory: 2,
			},
			versions: types.KeyVaultObjectVersionList{
				types.KeyVaultObjectVersion{
					Version: "a",
					Created: now.Add(time.Hour * 10),
				},
				types.KeyVaultObjectVersion{
					Version: "b",
					Created: now.Add(time.Hour * 9),
				},
			},
			expectedObjects: []types.KeyVaultCertificate{
				{
					CertificateName:           "cert1",
					CertificateAlias:          filepath.Join("cert1", "0"),
					CertificateVersion:        "a",
					CertificateVersionHistory: 2,
				},
				{
					CertificateName:           "cert1",
					CertificateAlias:          filepath.Join("cert1", "1"),
					CertificateVersion:        "b",
					CertificateVersionHistory: 2,
				},
			},
		},
		{
			desc: "fewer than CertificateVersionHistory results returns all versions",
			kvObject: types.KeyVaultCertificate{
				CertificateName:           "cert1",
				CertificateVersionHistory: 200,
			},
			versions: types.KeyVaultObjectVersionList{
				types.KeyVaultObjectVersion{
					Version: "a",
					Created: now.Add(time.Hour * 10),
				},
				types.KeyVaultObjectVersion{
					Version: "b",
					Created: now.Add(time.Hour * 9),
				},
			},
			expectedObjects: []types.KeyVaultCertificate{
				{
					CertificateName:           "cert1",
					CertificateAlias:          filepath.Join("cert1", "0"),
					CertificateVersion:        "a",
					CertificateVersionHistory: 200,
				},
				{
					CertificateName:           "cert1",
					CertificateAlias:          filepath.Join("cert1", "1"),
					CertificateVersion:        "b",
					CertificateVersionHistory: 200,
				},
			},
		},
		{
			desc: "starts at CertificateVersion when specified",
			kvObject: types.KeyVaultCertificate{
				CertificateName:           "cert1",
				CertificateVersion:        "c",
				CertificateVersionHistory: 5,
			},
			versions: types.KeyVaultObjectVersionList{
				types.KeyVaultObjectVersion{
					Version: "c",
					Created: now.Add(time.Hour * 8),
				},
				types.KeyVaultObjectVersion{
					Version: "e",
					Created: now.Add(time.Hour * 6),
				},
				types.KeyVaultObjectVersion{
					Version: "b",
					Created: now.Add(time.Hour * 9),
				},
				types.KeyVaultObjectVersion{
					Version: "a",
					Created: now.Add(time.Hour * 10),
				},
				types.KeyVaultObjectVersion{
					Version: "d",
					Created: now.Add(time.Hour * 7),
				},
			},
			expectedObjects: []types.KeyVaultCertificate{
				{
					CertificateName:           "cert1",
					CertificateAlias:          filepath.Join("cert1", "0"),
					CertificateVersion:        "c",
					CertificateVersionHistory: 5,
				},
				{
					CertificateName:           "cert1",
					CertificateAlias:          filepath.Join("cert1", "1"),
					CertificateVersion:        "d",
					CertificateVersionHistory: 5,
				},
				{
					CertificateName:           "cert1",
					CertificateAlias:          filepath.Join("cert1", "2"),
					CertificateVersion:        "e",
					CertificateVersionHistory: 5,
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			actualObjects := getLatestNKeyVaultObjects(tc.kvObject, tc.versions)

			if !reflect.DeepEqual(actualObjects, tc.expectedObjects) {
				t.Fatalf("expected: %+v, but got: %+v", tc.expectedObjects, actualObjects)
			}
		})
	}
}

func TestFormatKeyVaultCertificate(t *testing.T) {
	cases := []struct {
		desc                   string
		keyVaultObject         types.KeyVaultCertificate
		expectedKeyVaultObject types.KeyVaultCertificate
	}{
		{
			desc: "leading and trailing whitespace trimmed from all fields",
			keyVaultObject: types.KeyVaultCertificate{
				CertificateName:    "cert1     ",
				CertificateVersion: "",
				CertificateAlias:   "",
			},
			expectedKeyVaultObject: types.KeyVaultCertificate{
				CertificateName:    "cert1",
				CertificateVersion: "",
				CertificateAlias:   "",
			},
		},
		{
			desc: "no data loss for already sanitized object",
			keyVaultObject: types.KeyVaultCertificate{
				CertificateName:    "cert1",
				CertificateVersion: "version1",
				CertificateAlias:   "alias",
			},
			expectedKeyVaultObject: types.KeyVaultCertificate{
				CertificateName:    "cert1",
				CertificateVersion: "version1",
				CertificateAlias:   "alias",
			},
		},
		{
			desc: "no data loss for int properties",
			keyVaultObject: types.KeyVaultCertificate{
				CertificateName:           "cert1",
				CertificateVersion:        "latest",
				CertificateAlias:          "alias",
				CertificateVersionHistory: 12,
			},
			expectedKeyVaultObject: types.KeyVaultCertificate{
				CertificateName:           "cert1",
				CertificateVersion:        "latest",
				CertificateAlias:          "alias",
				CertificateVersionHistory: 12,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			formatKeyVaultCertificate(&tc.keyVaultObject)
			if !reflect.DeepEqual(tc.keyVaultObject, tc.expectedKeyVaultObject) {
				t.Fatalf("expected: %+v, but got: %+v", tc.expectedKeyVaultObject, tc.keyVaultObject)
			}
		})
	}
}

/*func SkipTestInitializeKVClient(t *testing.T) {
	testEnvs := []azure.Environment{
		azure.PublicCloud,
		azure.GermanCloud,
		azure.ChinaCloud,
		azure.USGovernmentCloud,
	}

	for i := range testEnvs {

		kvBaseClient, err := initializeKvClient(context.TODO(), testEnvs[i].KeyVaultEndpoint, testEnvs[i].ActiveDirectoryEndpoint, "", "")
		assert.NoError(t, err)
		assert.NotNil(t, kvBaseClient)
		assert.NotNil(t, kvBaseClient.Authorizer)
		assert.Contains(t, kvBaseClient.UserAgent, "ratify")
	}
}*/

func TestGetCertificatesContent(t *testing.T) {
	cases := []struct {
		desc        string
		parameters  map[string]string
		secrets     map[string]string
		expectedErr bool
	}{
		{
			desc:        "keyvault name not provided",
			parameters:  map[string]string{},
			expectedErr: true,
		},
		{
			desc: "tenantID not provided",
			parameters: map[string]string{
				"keyvaultName": "testKV",
			},
			expectedErr: true,
		},
		{
			desc: "usePodIdentity not a boolean as expected",
			parameters: map[string]string{
				"keyvaultName":   "testKV",
				"tenantId":       "tid",
				"usePodIdentity": "tru",
			},
			expectedErr: true,
		},
		{
			desc: "useVMManagedIdentity not a boolean as expected",
			parameters: map[string]string{
				"keyvaultName":         "testKV",
				"tenantId":             "tid",
				"usePodIdentity":       "false",
				"useVMManagedIdentity": "tru",
			},
			expectedErr: true,
		},
		{
			desc: "invalid cloud name",
			parameters: map[string]string{
				"keyvaultName": "testKV",
				"tenantId":     "tid",
				"cloudName":    "AzureCloud",
			},
			expectedErr: true,
		},
		{
			desc: "check azure cloud env file path is set",
			parameters: map[string]string{
				"keyvaultName":     "testKV",
				"tenantId":         "tid",
				"cloudName":        "AzureStackCloud",
				"cloudEnvFileName": "/etc/kubernetes/akscustom.json",
			},
			expectedErr: true,
		},
		{
			desc: "objects array not set",
			parameters: map[string]string{
				"keyvaultName":         "testKV",
				"tenantId":             "tid",
				"useVMManagedIdentity": "true",
			},
			expectedErr: true,
		},
		{
			desc: "objects not configured as an array",
			parameters: map[string]string{
				"keyvaultName":         "testKV",
				"tenantId":             "tid",
				"useVMManagedIdentity": "true",
				"objects": `
        - |
          CertificateName: cert1
          objectType: secret
          CertificateVersion: ""`,
			},
			expectedErr: true,
		},
		{
			desc: "objects array is empty",
			parameters: map[string]string{
				"keyvaultName":         "testKV",
				"tenantId":             "tid",
				"useVMManagedIdentity": "true",
				"objects": `
      array:`,
			},
			expectedErr: false,
		},
		{
			desc: "invalid object format",
			parameters: map[string]string{
				"keyvaultName":         "testKV",
				"tenantId":             "tid",
				"useVMManagedIdentity": "true",
				"objects": `
      array:
        - |
          CertificateName: cert1
          objectType: secret
          objectFormat: pkcs
          CertificateVersion: ""`,
			},
			expectedErr: true,
		},
		{
			desc: "invalid object encoding",
			parameters: map[string]string{
				"keyvaultName":         "testKV",
				"tenantId":             "tid",
				"useVMManagedIdentity": "true",
				"objects": `
      array:
        - |
          CertificateName: cert1
          objectType: secret
          objectEncoding: utf-16
          CertificateVersion: ""`,
			},
			expectedErr: true,
		},
		{
			desc: "error fetching from keyvault",
			parameters: map[string]string{
				"keyvaultName": "testKV",
				"tenantId":     "tid",
				"objects": `
      array:
        - |
          CertificateName: cert1
          objectType: secret
          CertificateVersion: ""`,
			},
			secrets: map[string]string{
				"clientid":     "AADClientID",
				"clientsecret": "AADClientSecret",
			},
			expectedErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {

			_, err := GetCertificatesContent(context.TODO(), tc.parameters)
			if tc.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetObjectVersion(t *testing.T) {
	id := "https://kindkv.vault.azure.net/secrets/cert1/c55925c29c6743dcb9bb4bf091be03b0"
	expectedVersion := "c55925c29c6743dcb9bb4bf091be03b0"
	actual := getObjectVersion(id)
	assert.Equal(t, expectedVersion, actual)
}
