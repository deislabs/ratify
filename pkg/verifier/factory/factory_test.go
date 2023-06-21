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

package factory

import (
	"context"
	"testing"

	"github.com/deislabs/ratify/pkg/common"
	"github.com/deislabs/ratify/pkg/ocispecs"
	"github.com/deislabs/ratify/pkg/referrerstore"

	"github.com/deislabs/ratify/pkg/verifier"
	"github.com/deislabs/ratify/pkg/verifier/config"
	"github.com/deislabs/ratify/pkg/verifier/plugin"
)

type TestVerifier struct {
	verifierDirectory string
}
type TestVerifierFactory struct{}

func (s *TestVerifier) Name() string {
	return "test-verifier"
}

func (s *TestVerifier) CanVerify(_ context.Context, _ ocispecs.ReferenceDescriptor) bool {
	return true
}

func (s *TestVerifier) Verify(_ context.Context,
	_ common.Reference,
	_ ocispecs.ReferenceDescriptor,
	_ referrerstore.ReferrerStore) (verifier.VerifierResult, error) {
	return verifier.VerifierResult{IsSuccess: false}, nil
}

func (s *TestVerifier) GetNestedReferences() []string {
	return []string{}
}

func (f *TestVerifierFactory) Create(_ string, _ config.VerifierConfig, pluginDirectory string) (verifier.ReferenceVerifier, error) {
	return &TestVerifier{verifierDirectory: pluginDirectory}, nil
}

func TestCreateVerifiersFromConfig_BuiltInVerifiers_ReturnsExpected(t *testing.T) {
	builtInVerifiers = map[string]VerifierFactory{
		"test-verifier": &TestVerifierFactory{},
	}

	verifierConfig := map[string]interface{}{
		"name": "test-verifier",
	}
	verifiersConfig := config.VerifiersConfig{
		Verifiers: []config.VerifierConfig{verifierConfig},
	}

	verifiers, err := CreateVerifiersFromConfig(verifiersConfig, "test/dir")

	if err != nil {
		t.Fatalf("create verifiers failed with err %v", err)
	}

	if len(verifiers) != 1 {
		t.Fatalf("expected to have %d verifiers, actual count %d", 1, len(verifiers))
	}

	if verifiers[0].Name() != "test-verifier" {
		t.Fatalf("expected to create test verifier")
	}

	if _, ok := verifiers[0].(*plugin.VerifierPlugin); ok {
		t.Fatalf("type assertion failed expected a built in verifier")
	}

	if verifierTest, ok := verifiers[0].(*TestVerifier); !ok {
		t.Fatalf("type assertion failed expected a test verifier")
	} else {
		if verifierTest.verifierDirectory != "test/dir" {
			t.Fatalf("expected verifier directory to be empty")
		}
	}
}

func TestCreateVerifiersFromConfig_PluginVerifiers_ReturnsExpected(t *testing.T) {
	verifierConfig := map[string]interface{}{
		"name": "plugin-verifier",
	}
	verifiersConfig := config.VerifiersConfig{
		Verifiers: []config.VerifierConfig{verifierConfig},
	}

	verifiers, err := CreateVerifiersFromConfig(verifiersConfig, "")

	if err != nil {
		t.Fatalf("create verifiers failed with err %v", err)
	}

	if len(verifiers) != 1 {
		t.Fatalf("expected to have %d verifiers, actual count %d", 1, len(verifiers))
	}

	if verifiers[0].Name() != "plugin-verifier" {
		t.Fatalf("expected to create plugin verifier")
	}

	if _, ok := verifiers[0].(*plugin.VerifierPlugin); !ok {
		t.Fatalf("type assertion failed expected a plugin in verifier")
	}
}
