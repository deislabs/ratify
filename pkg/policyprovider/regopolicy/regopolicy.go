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

package regopolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/deislabs/ratify/pkg/common"
	"github.com/deislabs/ratify/pkg/executor/types"
	"github.com/deislabs/ratify/pkg/ocispecs"
	"github.com/deislabs/ratify/pkg/policyprovider"
	"github.com/deislabs/ratify/pkg/policyprovider/config"
	pf "github.com/deislabs/ratify/pkg/policyprovider/factory"
	"github.com/deislabs/ratify/pkg/policyprovider/policyengine"
	opa "github.com/deislabs/ratify/pkg/policyprovider/policyengine/opaengine"
	query "github.com/deislabs/ratify/pkg/policyprovider/policyquery/rego"
	policyTypes "github.com/deislabs/ratify/pkg/policyprovider/types"
	"github.com/sirupsen/logrus"
)

type policyEnforcer struct {
	Policy             string
	OpaEngine          policyengine.PolicyEngine
	passthroughEnabled bool
}

type policyEnforcerConf struct {
	Name               string `json:"name"`
	Policy             string `json:"policy"`
	PolicyPath         string `json:"policyPath"`
	PassthroughEnabled bool   `json:"passthroughEnabled"`
}

// Factory is a factory for creating rego policy enforcers.
type Factory struct{}

// init calls Register for our rego policy provider.
func init() {
	pf.Register(policyTypes.RegoPolicy, &Factory{})
}

// Create creates a new policy enforcer based on the policy provided in config.
func (f *Factory) Create(policyConfig config.PolicyPluginConfig) (policyprovider.PolicyProvider, error) {
	conf := policyEnforcerConf{}
	policyProviderConfigBytes, err := json.Marshal(policyConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal policy config: %w", err)
	}

	if err := json.Unmarshal(policyProviderConfigBytes, &conf); err != nil {
		return nil, fmt.Errorf("failed to parse policy provider configuration: %w", err)
	}
	if conf.Policy == "" {
		body, err := os.ReadFile(conf.PolicyPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read rego policy file at path: %s, err: %w", conf.PolicyPath, err)
		}
		conf.Policy = string(body)
	}
	if conf.Policy == "" {
		return nil, fmt.Errorf("policy is required for rego policy provider")
	}

	engine, err := policyengine.CreateEngineFromConfig(policyengine.Config{
		Name:          opa.OPA,
		QueryLanguage: query.RegoName,
		Policy:        conf.Policy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create OPA engine: %w", err)
	}

	policyEnforcer := &policyEnforcer{
		Policy:             conf.Policy,
		OpaEngine:          engine,
		passthroughEnabled: conf.PassthroughEnabled,
	}

	return policyEnforcer, nil
}

// VerifyNeeded determines if verification should be performed for a given artifact.
func (e *policyEnforcer) VerifyNeeded(_ context.Context, _ common.Reference, _ ocispecs.ReferenceDescriptor) bool {
	return true
}

// ContinueVerifyOnFailure determines if verification should continue if a previous verification failed.
func (e *policyEnforcer) ContinueVerifyOnFailure(_ context.Context, _ common.Reference, _ ocispecs.ReferenceDescriptor, _ types.VerifyResult) bool {
	return true
}

// ErrorToVerifyResult converts an error to a VerifyResult.
func (e *policyEnforcer) ErrorToVerifyResult(_ context.Context, _ string, _ error) types.VerifyResult {
	return types.VerifyResult{}
}

// OverallVerifyResult determines if the overall verification result should be a success or failure.
func (e *policyEnforcer) OverallVerifyResult(ctx context.Context, verifierReports []interface{}) bool {
	if e.passthroughEnabled {
		return false
	}

	nestedReports := map[string]interface{}{}
	nestedReports["verifierReports"] = verifierReports
	result, err := e.OpaEngine.Evaluate(ctx, nestedReports)
	if err != nil {
		logrus.Errorf("failed to evaluate policy: %v", err)
		return false
	}
	return result
}

// GetPolicyType returns the type of the policy.
func (e *policyEnforcer) GetPolicyType(_ context.Context) string {
	return policyTypes.RegoPolicy
}
