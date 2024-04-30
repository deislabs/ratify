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

package cosign

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"testing"

	"github.com/deislabs/ratify/pkg/keymanagementprovider"
)

func TestCreateTrustPolicy(t *testing.T) {
	tc := []struct {
		name    string
		cfg     TrustPolicyConfig
		wantErr bool
	}{
		{
			name:    "invalid config",
			cfg:     TrustPolicyConfig{},
			wantErr: true,
		},
		{
			name: "invalid local key path",
			cfg: TrustPolicyConfig{
				Name:   "test",
				Scopes: []string{"*"},
				Keys: []KeyConfig{
					{
						File: "invalid",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid local key path",
			cfg: TrustPolicyConfig{
				Name:   "test",
				Scopes: []string{"*"},
				Keys: []KeyConfig{
					{
						File: "../../../test/testdata/cosign.pub",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid keyless config with rekor specified",
			cfg: TrustPolicyConfig{
				Name:   "test",
				Scopes: []string{"*"},
				Keyless: KeylessConfig{
					RekorURL: DefaultRekorURL,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid config version",
			cfg: TrustPolicyConfig{
				Version: "0.0.0",
				Name:    "test",
				Scopes:  []string{"*"},
				Keyless: KeylessConfig{
					RekorURL: DefaultRekorURL,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateTrustPolicy(tt.cfg, "test-verifier")
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

// TestGetName tests the GetName function for Trust Policy
func TestGetName(t *testing.T) {
	trustPolicyConfig := TrustPolicyConfig{
		Name:    "test",
		Scopes:  []string{"*"},
		Keyless: KeylessConfig{RekorURL: DefaultRekorURL},
	}
	trustPolicy, err := CreateTrustPolicy(trustPolicyConfig, "test-verifier")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if trustPolicy.GetName() != trustPolicyConfig.Name {
		t.Fatalf("expected %s, got %s", trustPolicyConfig.Name, trustPolicy.GetName())
	}
}

// TestGetScopes tests the GetScopes function for Trust Policy
func TestGetScopes(t *testing.T) {
	trustPolicyConfig := TrustPolicyConfig{
		Name:    "test",
		Scopes:  []string{"*"},
		Keyless: KeylessConfig{RekorURL: DefaultRekorURL},
	}
	trustPolicy, err := CreateTrustPolicy(trustPolicyConfig, "test-verifier")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(trustPolicy.GetScopes()) != len(trustPolicyConfig.Scopes) {
		t.Fatalf("expected %v, got %v", trustPolicyConfig.Scopes, trustPolicy.GetScopes())
	}
	if trustPolicy.GetScopes()[0] != trustPolicyConfig.Scopes[0] {
		t.Fatalf("expected %s, got %s", trustPolicyConfig.Scopes[0], trustPolicy.GetScopes()[0])
	}
}

// TestGetKeys tests the GetKeys function for Trust Policy
func TestGetKeys(t *testing.T) {
	inputMap := map[keymanagementprovider.KMPMapKey]crypto.PublicKey{
		{Name: "key1"}: &ecdsa.PublicKey{},
	}
	keymanagementprovider.SetKeysInMap("ns/kmp", "", inputMap)
	tc := []struct {
		name    string
		cfg     TrustPolicyConfig
		wantErr bool
	}{
		{
			name: "only local keys",
			cfg: TrustPolicyConfig{
				Name:   "test",
				Scopes: []string{"*"},
				Keys: []KeyConfig{
					{
						File: "../../../test/testdata/cosign.pub",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nonexistent KMP",
			cfg: TrustPolicyConfig{
				Name:   "test",
				Scopes: []string{"*"},
				Keys: []KeyConfig{
					{
						Provider: "nonexistent",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid KMP",
			cfg: TrustPolicyConfig{
				Name:   "test",
				Scopes: []string{"*"},
				Keys: []KeyConfig{
					{
						Provider: "kmp",
						Name:     "key1",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			trustPolicy, err := CreateTrustPolicy(tt.cfg, "test-verifier")
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			keys, err := trustPolicy.GetKeys(context.Background(), "ns")
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
			if err == nil && len(keys) != len(tt.cfg.Keys) {
				t.Fatalf("expected %v, got %v", tt.cfg.Keys, keys)
			}
		})
	}
}

// TestValidate tests the validate function
func TestValidate(t *testing.T) {
	tc := []struct {
		name         string
		policyConfig TrustPolicyConfig
		wantErr      bool
	}{
		{
			name:         "no name",
			policyConfig: TrustPolicyConfig{},
			wantErr:      true,
		},
		{
			name: "no scopes",
			policyConfig: TrustPolicyConfig{
				Name: "test",
			},
			wantErr: true,
		},
		{
			name: "no keys or keyless defined",
			policyConfig: TrustPolicyConfig{
				Name:   "test",
				Scopes: []string{"*"},
			},
			wantErr: true,
		},
		{
			name: "keys and keyless defined",
			policyConfig: TrustPolicyConfig{
				Name:   "test",
				Scopes: []string{"*"},
				Keys: []KeyConfig{
					{
						Provider: "kmp",
					},
				},
				Keyless: KeylessConfig{RekorURL: DefaultRekorURL},
			},
			wantErr: true,
		},
		{
			name: "key provider and key path not defined",
			policyConfig: TrustPolicyConfig{
				Name:   "test",
				Scopes: []string{"*"},
				Keys:   []KeyConfig{{}},
			},
			wantErr: true,
		},
		{
			name: "key provider and key path both defined",
			policyConfig: TrustPolicyConfig{
				Name:   "test",
				Scopes: []string{"*"},
				Keys: []KeyConfig{
					{
						Provider: "kmp",
						File:     "path",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "key provider not defined but name defined",
			policyConfig: TrustPolicyConfig{
				Name:   "test",
				Scopes: []string{"*"},
				Keys: []KeyConfig{
					{
						Name: "key name",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "key provider name not defined but version defined",
			policyConfig: TrustPolicyConfig{
				Name:   "test",
				Scopes: []string{"*"},
				Keys: []KeyConfig{
					{
						Provider: "kmp",
						Version:  "key version",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid",
			policyConfig: TrustPolicyConfig{
				Version: "1.0.0",
				Name:    "test",
				Scopes:  []string{"*"},
				Keys: []KeyConfig{
					{
						Provider: "kmp",
						Name:     "key name",
						Version:  "key version",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			actual := validate(tt.policyConfig, "test-verifier")
			if (actual != nil) != tt.wantErr {
				t.Fatalf("expected %v, got %v", tt.wantErr, actual)
			}
		})
	}
}

// TestLoadKeyFromPath tests the loadKeyFromPath function
func TestLoadKeyFromPath(t *testing.T) {
	cosignValidPath := "../../../test/testdata/cosign.pub"
	key, err := loadKeyFromPath(cosignValidPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if key == nil {
		t.Fatalf("expected key, got nil")
	}
	switch keyType := key.(type) {
	case *ecdsa.PublicKey:
	default:
		t.Fatalf("expected ecdsa.PublicKey, got %v", keyType)
	}
}

// TestPrependNamespaceToKMPName tests the prependNamespaceToKMPName function
func TestPrependNamespaceToKMPName(t *testing.T) {
	tc := []struct {
		name     string
		kmpName  string
		ns       string
		expected string
	}{
		{
			name:     "empty namespace",
			kmpName:  "kmp",
			ns:       "",
			expected: "kmp",
		},
		{
			name:     "non-empty namespace",
			kmpName:  "kmp",
			ns:       "ns",
			expected: "ns/kmp",
		},
		{
			name:     "namespaced kmp",
			kmpName:  "ns/kmp",
			ns:       "ns",
			expected: "ns/kmp",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			actual := prependNamespaceToKMPName(tt.kmpName, tt.ns)
			if actual != tt.expected {
				t.Fatalf("expected %s, got %s", tt.expected, actual)
			}
		})
	}
}
