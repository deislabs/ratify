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

package notation

import (
	"context"
	"encoding/json"
	"fmt"
	paths "path/filepath"
	"strings"

	ratifyconfig "github.com/deislabs/ratify/config"
	re "github.com/deislabs/ratify/errors"
	"github.com/deislabs/ratify/internal/constants"
	"github.com/deislabs/ratify/internal/logger"
	"github.com/deislabs/ratify/pkg/common"
	"github.com/deislabs/ratify/pkg/homedir"

	"github.com/deislabs/ratify/pkg/ocispecs"
	"github.com/deislabs/ratify/pkg/referrerstore"
	"github.com/deislabs/ratify/pkg/verifier"
	"github.com/deislabs/ratify/pkg/verifier/config"
	"github.com/deislabs/ratify/pkg/verifier/factory"
	"github.com/deislabs/ratify/pkg/verifier/types"
	"github.com/notaryproject/notation-go/log"

	_ "github.com/notaryproject/notation-core-go/signature/cose" // register COSE signature
	_ "github.com/notaryproject/notation-core-go/signature/jws"  // register JWS signature
	"github.com/notaryproject/notation-go"
	notationVerifier "github.com/notaryproject/notation-go/verifier"
	"github.com/notaryproject/notation-go/verifier/trustpolicy"
	"github.com/notaryproject/notation-go/verifier/truststore"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	verifierType                      = "notation"
	defaultCertPath                   = "ratify-certs/notation/truststore"
	trustStoreTypeCA                  = string(truststore.TypeCA)
	trustStoreTypeypeSigningAuthority = string(truststore.TypeSigningAuthority)
)

// NotationPluginVerifierConfig describes the configuration of notation verifier
type NotationPluginVerifierConfig struct { //nolint:revive // ignore linter to have unique type name
	Name          string `json:"name"`
	ArtifactTypes string `json:"artifactTypes"`

	// VerificationCerts is array of directories containing certificates.
	VerificationCerts []string `json:"verificationCerts"`
	// VerificationCertStores is map defining which keyvault certificates belong to which trust store name.
	// VerificationCertStores accepts new format map[string]map[string][]string
	// {
	// 	"ca": {
	// 		"certs": {"kv1", "kv2"},
	// 	},
	// 	"signingauthority": {
	// 		"certs": {"kv3"}
	// 	},
	// }
	// VerificationCertStores accepts legacy format map[string][]string as well.
	// {
	// 	"certs": {"kv1", "kv2"},
	// },
	VerificationCertStores map[string]interface{} `json:"verificationCertStores"`
	// TrustPolicyDoc represents a trustpolicy.json document. Reference: https://pkg.go.dev/github.com/notaryproject/notation-go@v0.12.0-beta.1.0.20221125022016-ab113ebd2a6c/verifier/trustpolicy#Document
	TrustPolicyDoc trustpolicy.Document `json:"trustPolicyDoc"`
}

type notationPluginVerifier struct {
	name             string
	verifierType     string
	artifactTypes    []string
	notationVerifier *notation.Verifier
}

type notationPluginVerifierFactory struct{}

func init() {
	factory.Register(verifierType, &notationPluginVerifierFactory{})
}

func (f *notationPluginVerifierFactory) Create(_ string, verifierConfig config.VerifierConfig, pluginDirectory string, namespace string) (verifier.ReferenceVerifier, error) {
	logger.GetLogger(context.Background(), logOpt).Debugf("creating notation with config %v, namespace '%v'", verifierConfig, namespace)
	verifierName := fmt.Sprintf("%s", verifierConfig[types.Name])
	verifierTypeStr := ""
	if _, ok := verifierConfig[types.Type]; ok {
		verifierTypeStr = fmt.Sprintf("%s", verifierConfig[types.Type])
	}
	conf, err := parseVerifierConfig(verifierConfig, namespace)
	if err != nil {
		return nil, re.ErrorCodeConfigInvalid.WithComponentType(re.Verifier).WithPluginName(verifierName)
	}

	verifyService, err := getVerifierService(conf, pluginDirectory)
	if err != nil {
		return nil, re.ErrorCodePluginInitFailure.WithComponentType(re.Verifier).WithPluginName(verifierName).WithError(err)
	}

	artifactTypes := strings.Split(conf.ArtifactTypes, ",")
	return &notationPluginVerifier{
		name:             verifierName,
		verifierType:     verifierTypeStr,
		artifactTypes:    artifactTypes,
		notationVerifier: &verifyService,
	}, nil
}

func (v *notationPluginVerifier) Name() string {
	return v.name
}

func (v *notationPluginVerifier) Type() string {
	return v.verifierType
}

func (v *notationPluginVerifier) CanVerify(_ context.Context, referenceDescriptor ocispecs.ReferenceDescriptor) bool {
	for _, at := range v.artifactTypes {
		if at == "*" || at == referenceDescriptor.ArtifactType {
			return true
		}
	}
	return false
}

func (v *notationPluginVerifier) Verify(ctx context.Context,
	subjectReference common.Reference,
	referenceDescriptor ocispecs.ReferenceDescriptor,
	store referrerstore.ReferrerStore) (verifier.VerifierResult, error) {
	extensions := make(map[string]string)

	subjectDesc, err := store.GetSubjectDescriptor(ctx, subjectReference)
	if err != nil {
		return verifier.VerifierResult{IsSuccess: false}, re.ErrorCodeGetSubjectDescriptorFailure.NewError(re.ReferrerStore, store.Name(), re.EmptyLink, err, fmt.Sprintf("failed to resolve subject: %+v", subjectReference), re.HideStackTrace)
	}

	referenceManifest, err := store.GetReferenceManifest(ctx, subjectReference, referenceDescriptor)
	if err != nil {
		return verifier.VerifierResult{IsSuccess: false}, re.ErrorCodeGetReferenceManifestFailure.NewError(re.ReferrerStore, store.Name(), re.EmptyLink, err, fmt.Sprintf("failed to resolve reference manifest: %+v", referenceDescriptor), re.HideStackTrace)
	}

	if len(referenceManifest.Blobs) == 0 {
		return verifier.VerifierResult{IsSuccess: false}, re.ErrorCodeSignatureNotFound.NewError(re.Verifier, v.name, re.EmptyLink, nil, fmt.Sprintf("no signature content found for referrer: %s@%s", subjectReference.Path, referenceDescriptor.Digest.String()), re.HideStackTrace)
	}

	for _, blobDesc := range referenceManifest.Blobs {
		refBlob, err := store.GetBlobContent(ctx, subjectReference, blobDesc.Digest)
		if err != nil {
			return verifier.VerifierResult{IsSuccess: false}, re.ErrorCodeGetBlobContentFailure.NewError(re.ReferrerStore, store.Name(), re.EmptyLink, err, fmt.Sprintf("failed to get blob content of digest: %s", blobDesc.Digest), re.HideStackTrace)
		}

		// TODO: notation verify API only accepts digested reference now.
		// Pass in tagged reference instead once notation-go supports it.
		subjectRef := fmt.Sprintf("%s@%s", subjectReference.Path, subjectReference.Digest.String())
		outcome, err := v.verifySignature(ctx, subjectRef, blobDesc.MediaType, subjectDesc.Descriptor, refBlob)
		if err != nil {
			return verifier.VerifierResult{IsSuccess: false, Extensions: extensions}, re.ErrorCodeVerifyPluginFailure.NewError(re.Verifier, v.name, re.NotationTsgLink, err, "failed to verify signature of digest", re.HideStackTrace)
		}

		// Note: notation verifier already validates certificate chain is not empty.
		cert := outcome.EnvelopeContent.SignerInfo.CertificateChain[0]
		extensions["Issuer"] = cert.Issuer.String()
		extensions["SN"] = cert.Subject.String()
	}

	return verifier.VerifierResult{
		Name:       v.name,
		Type:       v.verifierType,
		IsSuccess:  true,
		Message:    "signature verification success",
		Extensions: extensions,
	}, nil
}

func getVerifierService(conf *NotationPluginVerifierConfig, pluginDirectory string) (notation.Verifier, error) {
	certStoresByType, err := reformVerificationCertsStores(conf)
	if err != nil {
		return nil, err
	}
	store := &trustStore{
		certPaths:        conf.VerificationCerts,
		certStoresByType: certStoresByType,
	}

	return notationVerifier.New(&conf.TrustPolicyDoc, store, NewRatifyPluginManager(pluginDirectory))
}

// reformVerificationCertsStores assert VerificationCertStores of type map[string]interface{} to certStoresByType map[string]map[string][]string
func reformVerificationCertsStores(conf *NotationPluginVerifierConfig) (map[string]map[string][]string, error) {
	certStoresByType := make(map[string]map[string][]string)
	for certStoreType, certStores := range conf.VerificationCertStores {
		if reformedCertStores, ok := certStores.(map[string]interface{}); ok {
			certStoresByType[certStoreType] = make(map[string][]string)
			for certStore, certs := range reformedCertStores {
				var reformedCerts []string
				for _, cert := range certs.([]interface{}) {
					if reformedCert, ok := cert.(string); ok {
						reformedCerts = append(reformedCerts, reformedCert)
					}
				}
				certStoresByType[certStoreType][certStore] = reformedCerts
			}
		} else {
			return nil, re.ErrorCodeEnvNotSet.WithComponentType(re.Verifier).WithDetail("")
		}
	}
	return certStoresByType, nil
}

func (v *notationPluginVerifier) verifySignature(ctx context.Context, subjectRef, mediaType string, subjectDesc oci.Descriptor, refBlob []byte) (*notation.VerificationOutcome, error) {
	opts := notation.VerifierVerifyOptions{
		SignatureMediaType: mediaType,
		ArtifactReference:  subjectRef,
	}
	ctx = log.WithLogger(ctx, logger.GetLogger(ctx, logOpt))

	return (*v.notationVerifier).Verify(ctx, subjectDesc, refBlob, opts)
}

func parseVerifierConfig(verifierConfig config.VerifierConfig, namespace string) (*NotationPluginVerifierConfig, error) {
	verifierName := verifierConfig[types.Name].(string)
	conf := &NotationPluginVerifierConfig{}
	verifierConfigBytes, err := json.Marshal(verifierConfig)
	if err != nil {
		return nil, re.ErrorCodeConfigInvalid.NewError(re.Verifier, verifierName, re.EmptyLink, err, nil, re.HideStackTrace)
	}
	if err := json.Unmarshal(verifierConfigBytes, &conf); err != nil {
		return nil, re.ErrorCodeConfigInvalid.NewError(re.Verifier, verifierName, re.EmptyLink, err, fmt.Sprintf("failed to unmarshal to notationPluginVerifierConfig from: %+v.", verifierConfig), re.HideStackTrace)
	}
	if len(conf.VerificationCertStores) > 0 {

		err := normalizeVerificationCertsStores(conf, namespace)
		if err != nil {
			return nil, err
		}
	}

	defaultCertsDir := paths.Join(homedir.Get(), ratifyconfig.ConfigFileDir, defaultCertPath)
	conf.VerificationCerts = append(conf.VerificationCerts, defaultCertsDir)
	return conf, nil
}

// normalizeVerificationCertsStores normalize the structure does not match the latest spec
func normalizeVerificationCertsStores(conf *NotationPluginVerifierConfig, namespace string) error {
	isCertStoresForTSA, isLegacyCertStore := false, false
	for key, val := range conf.VerificationCertStores {
		if key != trustStoreTypeCA && key != trustStoreTypeypeSigningAuthority {
			isLegacyCertStore = true
			logger.GetLogger(context.Background(), logOpt).Debugf("VerificationCertStores is not empty, will append namespace %v to certificate store if resource does not already contain a namespace", namespace)
			if err := prependNamespaceToCertStore(val.([]interface{}), namespace); err != nil {
				return err
			}
		} else {
			isCertStoresForTSA = true
			for certStore, certs := range val.(map[string][]interface{}) {
				if err := prependNamespaceToCertStore(certs, namespace); err != nil {
					return err
				}
				conf.VerificationCertStores[key].(map[string][]interface{})[certStore] = certs
			}
		}
	}
	if isCertStoresForTSA && isLegacyCertStore {
		return re.ErrorCodeConfigInvalid.NewError(re.Verifier, conf.Name, re.EmptyLink, nil, "both  old VerificationCertStores and new VerificationCertStores are provided, please provide only one", re.HideStackTrace)
	} else if !isCertStoresForTSA && isLegacyCertStore {
		// normalize <store>:<certs> to ca:<store><certs> if no store type is provided
		conf.VerificationCertStores = map[string]interface{}{
			trustStoreTypeCA: conf.VerificationCertStores,
		}
	}
	return nil
}

// signatures should not have nested references
func (v *notationPluginVerifier) GetNestedReferences() []string {
	return []string{}
}

// append namespace to certStore so they are uniquely identifiable
func prependNamespaceToCertStore(certs []interface{}, namespace string) error {
	if namespace == "" {
		return re.ErrorCodeEnvNotSet.WithComponentType(re.Verifier).WithDetail("failure to parse VerificationCertStores, namespace for VerificationCertStores must be provided")
	}
	for i, certstore := range certs {
		if !isNamespacedNamed(certstore.(string)) {
			certs[i] = namespace + constants.NamespaceSeperator + certstore.(string)
		}
	}
	return nil
}

// return true if string looks like a K8s namespaced resource. e.g. namespace/name
func isNamespacedNamed(name string) bool {
	return strings.Contains(name, constants.NamespaceSeperator)
}
