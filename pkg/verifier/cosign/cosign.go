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
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/deislabs/ratify/internal/logger"
	"github.com/deislabs/ratify/pkg/common"
	"github.com/deislabs/ratify/pkg/keymanagementprovider"
	"github.com/deislabs/ratify/pkg/keymanagementprovider/azurekeyvault"
	"github.com/deislabs/ratify/pkg/ocispecs"
	"github.com/deislabs/ratify/pkg/referrerstore"
	"github.com/deislabs/ratify/pkg/utils"
	"github.com/deislabs/ratify/pkg/verifier"
	"github.com/deislabs/ratify/pkg/verifier/config"
	"github.com/deislabs/ratify/pkg/verifier/factory"
	"github.com/deislabs/ratify/pkg/verifier/types"
	"golang.org/x/crypto/cryptobyte"
	"golang.org/x/crypto/cryptobyte/asn1"

	re "github.com/deislabs/ratify/errors"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/opencontainers/go-digest"
	imgspec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/fulcio"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/rekor"
	"github.com/sigstore/cosign/v2/pkg/cosign"
	"github.com/sigstore/cosign/v2/pkg/cosign/bundle"
	"github.com/sigstore/cosign/v2/pkg/oci"
	"github.com/sigstore/cosign/v2/pkg/oci/static"
	"github.com/sigstore/sigstore/pkg/signature"
)

type PluginConfig struct {
	Name             string              `json:"name"`
	Type             string              `json:"type,omitempty"`
	ArtifactTypes    string              `json:"artifactTypes"`
	KeyRef           string              `json:"key,omitempty"`
	RekorURL         string              `json:"rekorURL,omitempty"`
	NestedReferences []string            `json:"nestedArtifactTypes,omitempty"`
	TrustPolicies    []TrustPolicyConfig `json:"trustPolicies,omitempty"`
}

// LegacyExtension is the structure for the verifier result extensions
// used for backwards compatibility with the legacy cosign verifier
type LegacyExtension struct {
	SignatureExtension []cosignExtension `json:"signatures,omitempty"`
}

// Extension is the structure for the verifier result extensions
// contains a list of signature verification results
// where each entry corresponds to a single signature verified
type Extension struct {
	SignatureExtension []cosignExtensionList `json:"signatures,omitempty"`
}

// cosignExtensionList is the structure verifications performed
// per signature found in the image manifest
type cosignExtensionList struct {
	Signature     string            `json:"signature"`
	Verifications []cosignExtension `json:"verifications"`
}

// cosignExtension is the structure for the verification result
// of a single signature found in the image manifest for a
// single public key
type cosignExtension struct {
	SignatureDigest digest.Digest `json:"signatureDigest,omitempty"`
	IsSuccess       bool          `json:"isSuccess"`
	BundleVerified  bool          `json:"bundleVerified"`
	Err             string        `json:"error,omitempty"`
	KeyInformation  PKKey         `json:"keyInformation,omitempty"`
}

type cosignVerifier struct {
	name             string
	verifierType     string
	artifactTypes    []string
	nestedReferences []string
	config           *PluginConfig
	isLegacy         bool
	trustPolicies    *TrustPolicies
	namespace        string
}

type cosignVerifierFactory struct{}

var logOpt = logger.Option{
	ComponentType: logger.Verifier,
}

// used for mocking purposes
var getKeysMaps = getKeysMapsDefault

const verifierType string = "cosign"

// init() registers the cosign verifier with the factory
func init() {
	factory.Register(verifierType, &cosignVerifierFactory{})
}

// Create creates a new cosign verifier
func (f *cosignVerifierFactory) Create(_ string, verifierConfig config.VerifierConfig, _ string, namespace string) (verifier.ReferenceVerifier, error) {
	logger.GetLogger(context.Background(), logOpt).Debugf("creating cosign verifier with config %v, namespace '%v'", verifierConfig, namespace)
	verifierName, hasName := verifierConfig[types.Name].(string)
	if !hasName {
		return nil, re.ErrorCodeConfigInvalid.WithComponentType(re.Verifier).WithPluginName(verifierName).WithDetail("missing name in verifier config")
	}

	config, err := parseVerifierConfig(verifierConfig)
	if err != nil {
		return nil, re.ErrorCodeConfigInvalid.WithComponentType(re.Verifier).WithPluginName(verifierName)
	}

	// if key or rekorURL is provided, trustPolicies should not be provided
	if (config.KeyRef != "" || config.RekorURL != "") && len(config.TrustPolicies) > 0 {
		return nil, re.ErrorCodeConfigInvalid.WithComponentType(re.Verifier).WithPluginName(verifierName).WithDetail("'key' and 'rekorURL' are part of cosign legacy configuration and cannot be used with `trustPolicies`")
	}

	var trustPolicies *TrustPolicies
	legacy := true
	if config.KeyRef == "" && config.RekorURL == "" {
		trustPolicies, err = CreateTrustPolicies(config.TrustPolicies, verifierName)
		if err != nil {
			return nil, err
		}
		legacy = false
	}

	return &cosignVerifier{
		name:             verifierName,
		verifierType:     config.Type,
		artifactTypes:    strings.Split(config.ArtifactTypes, ","),
		nestedReferences: config.NestedReferences,
		config:           config,
		isLegacy:         legacy,
		trustPolicies:    trustPolicies,
		namespace:        namespace,
	}, nil
}

// Name returns the name of the cosign verifier
func (v *cosignVerifier) Name() string {
	return v.name
}

// Type returns 'cosign' as the type of the verifier
func (v *cosignVerifier) Type() string {
	return verifierType
}

// CanVerify returns true if the referenceDescriptor's artifact type is in the list of artifact types supported by the verifier
func (v *cosignVerifier) CanVerify(_ context.Context, referenceDescriptor ocispecs.ReferenceDescriptor) bool {
	for _, at := range v.artifactTypes {
		if at == "*" || at == referenceDescriptor.ArtifactType {
			return true
		}
	}
	return false
}

func (v *cosignVerifier) Verify(ctx context.Context, subjectReference common.Reference, referenceDescriptor ocispecs.ReferenceDescriptor, referrerStore referrerstore.ReferrerStore) (verifier.VerifierResult, error) {
	if v.isLegacy {
		return v.verifyLegacy(ctx, subjectReference, referenceDescriptor, referrerStore)
	}
	return v.verifyInternal(ctx, subjectReference, referenceDescriptor, referrerStore)
}

func (v *cosignVerifier) verifyInternal(ctx context.Context, subjectReference common.Reference, referenceDescriptor ocispecs.ReferenceDescriptor, referrerStore referrerstore.ReferrerStore) (verifier.VerifierResult, error) {
	// get the map of keys and relevant cosign options for that reference
	keysMap, cosignOpts, err := getKeysMaps(ctx, v.trustPolicies, subjectReference.Original, v.namespace)
	if err != nil {
		return errorToVerifyResult(v.name, v.verifierType, err), nil
	}

	// get the reference manifest (cosign oci image)
	referenceManifest, err := referrerStore.GetReferenceManifest(ctx, subjectReference, referenceDescriptor)
	if err != nil {
		return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("failed to get reference manifest: %w", err)), nil
	}

	// manifest must be an OCI Image
	if referenceManifest.MediaType != imgspec.MediaTypeImageManifest {
		return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("reference manifest is not an image")), nil
	}

	// get the subject image descriptor
	subjectDesc, err := referrerStore.GetSubjectDescriptor(ctx, subjectReference)
	if err != nil {
		return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("failed to create subject hash: %w", err)), nil
	}

	// create the hash of the subject image descriptor (used as the hashed payload)
	subjectDescHash := v1.Hash{
		Algorithm: subjectDesc.Digest.Algorithm().String(),
		Hex:       subjectDesc.Digest.Hex(),
	}

	sigExtensions := make([]cosignExtensionList, 0)
	hasValidSignature := false
	// check each signature found
	for _, blob := range referenceManifest.Blobs {
		extensionListEntry := cosignExtensionList{
			Signature:     blob.Annotations[static.SignatureAnnotationKey],
			Verifications: make([]cosignExtension, 0),
		}
		// fetch the blob content of the signature from the referrer store
		blobBytes, err := referrerStore.GetBlobContent(ctx, subjectReference, blob.Digest)
		if err != nil {
			return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("failed to get blob content: %w", err)), nil
		}
		// convert the blob to a static signature
		staticOpts, err := staticLayerOpts(blob)
		if err != nil {
			return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("failed to parse static signature opts: %w", err)), nil
		}
		sig, err := static.NewSignature(blobBytes, blob.Annotations[static.SignatureAnnotationKey], staticOpts...)
		if err != nil {
			return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("failed to generate static signature: %w", err)), nil
		}
		// check each key in the map of keys returned by the trust policy
		for mapKey, pubKey := range keysMap {
			hashType := crypto.SHA256
			// default hash type is SHA256 but for AKV scenarios, the hash type is determined by the key size
			// TODO: investigate if it's possible to extract hash type from sig directly. This is a workaround for now
			if pubKey.ProviderType == azurekeyvault.ProviderName {
				hashType, sig, err = processAKVSignature(blob.Annotations[static.SignatureAnnotationKey], sig, pubKey.Key, blobBytes, staticOpts)
				if err != nil {
					return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("failed to process AKV signature: %w", err)), nil
				}
			}

			// return the correct verifier based on public key type and bytes
			verifier, err := signature.LoadVerifier(pubKey.Key, hashType)
			if err != nil {
				return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("failed to load public key from provider [%s] name [%s] version [%s]: %w", mapKey.Provider, mapKey.Name, mapKey.Version, err)), nil
			}
			cosignOpts.SigVerifier = verifier
			// verify signature with cosign options + perform bundle verification
			bundleVerified, err := cosign.VerifyImageSignature(ctx, sig, subjectDescHash, &cosignOpts)
			extension := cosignExtension{
				IsSuccess:      true,
				BundleVerified: bundleVerified,
				KeyInformation: mapKey,
			}
			if err != nil {
				extension.IsSuccess = false
				extension.Err = err.Error()
			} else {
				hasValidSignature = true
			}
			extensionListEntry.Verifications = append(extensionListEntry.Verifications, extension)
		}

		// TODO: perform keyless verification instead if no keys are found
		sigExtensions = append(sigExtensions, extensionListEntry)
	}

	if hasValidSignature {
		return verifier.VerifierResult{
			Name:       v.name,
			Type:       v.verifierType,
			IsSuccess:  true,
			Message:    "cosign verification success. valid signatures found. please refer to extensions field for verifications performed.",
			Extensions: Extension{SignatureExtension: sigExtensions},
		}, nil
	}

	errorResult := errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("no valid signatures found"))
	errorResult.Extensions = Extension{SignatureExtension: sigExtensions}
	return errorResult, nil
}

// **LEGACY** This implementation will be removed in Ratify v2.0.0. Verify verifies the subject reference using the cosign verifier.
func (v *cosignVerifier) verifyLegacy(ctx context.Context, subjectReference common.Reference, referenceDescriptor ocispecs.ReferenceDescriptor, referrerStore referrerstore.ReferrerStore) (verifier.VerifierResult, error) {
	cosignOpts := &cosign.CheckOpts{
		ClaimVerifier: cosign.SimpleClaimVerifier,
	}

	var ecdsaVerifier signature.Verifier
	var roots *x509.CertPool
	var err error
	if v.config.KeyRef != "" {
		ecdsaVerifier, err = loadPublicKey(v.config.KeyRef)
		if err != nil {
			return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("failed to load public key: %w", err)), nil
		}
		cosignOpts.SigVerifier = ecdsaVerifier
	} else {
		roots, err = fulcio.GetRoots()
		if err != nil {
			return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("failed to get fulcio roots: %w", err)), nil
		}
		cosignOpts.RootCerts = roots
		if cosignOpts.RootCerts == nil {
			return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("failed to initialize root certificates")), nil
		}
	}

	if v.config.RekorURL != "" {
		cosignOpts.RekorClient, err = rekor.NewClient(v.config.RekorURL)
		if err != nil {
			return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("failed to create Rekor client from URL %s: %w", v.config.RekorURL, err)), nil
		}
		cosignOpts.CTLogPubKeys, err = cosign.GetCTLogPubs(ctx)
		if err != nil {
			return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("failed to set Certificate Transparency Log public keys: %w", err)), nil
		}
		// Fetches the Rekor public keys from the Rekor server
		cosignOpts.RekorPubKeys, err = cosign.GetRekorPubs(ctx)
		if err != nil {
			return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("failed to set Rekor public keys: %w", err)), nil
		}
	} else {
		// if no rekor url is provided, turn off transparency log verification and ignore SCTs
		cosignOpts.IgnoreTlog = true
		cosignOpts.IgnoreSCT = true
	}

	referenceManifest, err := referrerStore.GetReferenceManifest(ctx, subjectReference, referenceDescriptor)
	if err != nil {
		return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("failed to get reference manifest: %w", err)), nil
	}

	// manifest must be an OCI Image
	if referenceManifest.MediaType != imgspec.MediaTypeImageManifest {
		return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("reference manifest is not an image")), nil
	}

	subjectDesc, err := referrerStore.GetSubjectDescriptor(ctx, subjectReference)
	if err != nil {
		return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("failed to create subject hash: %w", err)), nil
	}
	subjectDescHash := v1.Hash{
		Algorithm: subjectDesc.Digest.Algorithm().String(),
		Hex:       subjectDesc.Digest.Hex(),
	}

	sigExtensions := make([]cosignExtension, 0)
	signatures := []oci.Signature{}
	for _, blob := range referenceManifest.Blobs {
		blobBytes, err := referrerStore.GetBlobContent(ctx, subjectReference, blob.Digest)
		if err != nil {
			return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("failed to get blob content: %w", err)), nil
		}
		staticOpts, err := staticLayerOpts(blob)
		if err != nil {
			return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("failed to parse static signature opts: %w", err)), nil
		}
		sig, err := static.NewSignature(blobBytes, blob.Annotations[static.SignatureAnnotationKey], staticOpts...)
		if err != nil {
			return errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("failed to generate static signature: %w", err)), nil
		}
		// The verification will return an error if the signature is not valid.
		bundleVerified, err := cosign.VerifyImageSignature(ctx, sig, subjectDescHash, cosignOpts)
		extension := cosignExtension{
			SignatureDigest: blob.Digest,
			IsSuccess:       true,
			BundleVerified:  bundleVerified,
		}
		if err != nil {
			extension.IsSuccess = false
			extension.Err = err.Error()
		} else {
			signatures = append(signatures, sig)
		}
		sigExtensions = append(sigExtensions, extension)
	}

	if len(signatures) > 0 {
		return verifier.VerifierResult{
			Name:       v.name,
			Type:       v.verifierType,
			IsSuccess:  true,
			Message:    "cosign verification success. valid signatures found",
			Extensions: LegacyExtension{SignatureExtension: sigExtensions},
		}, nil
	}

	errorResult := errorToVerifyResult(v.name, v.verifierType, fmt.Errorf("no valid signatures found"))
	errorResult.Extensions = LegacyExtension{SignatureExtension: sigExtensions}
	return errorResult, nil
}

// GetNestedReferences returns the nested reference artifact types configured
func (v *cosignVerifier) GetNestedReferences() []string {
	return v.nestedReferences
}

// ParseVerifierConfig parses the verifier config and returns a PluginConfig
func parseVerifierConfig(verifierConfig config.VerifierConfig) (*PluginConfig, error) {
	verifierName, hasName := verifierConfig[types.Name].(string)
	if !hasName {
		return nil, re.ErrorCodeConfigInvalid.NewError(re.Verifier, "", re.EmptyLink, nil, "missing name in verifier config", re.HideStackTrace)
	}
	conf := PluginConfig{}

	verifierConfigBytes, err := json.Marshal(verifierConfig)
	if err != nil {
		return nil, re.ErrorCodeConfigInvalid.NewError(re.Verifier, verifierName, re.EmptyLink, err, nil, re.HideStackTrace)
	}

	if err := json.Unmarshal(verifierConfigBytes, &conf); err != nil {
		return nil, re.ErrorCodeConfigInvalid.NewError(re.Verifier, verifierName, re.EmptyLink, err, fmt.Sprintf("failed to unmarshal to cosign verifier config from: %+v.", verifierConfig), re.HideStackTrace)
	}

	// if Type is not provided, use the Name as the Type (backwards compatibility)
	if conf.Type == "" {
		conf.Type = conf.Name
	}

	if conf.ArtifactTypes == "" {
		return nil, re.ErrorCodeConfigInvalid.NewError(re.Verifier, verifierName, re.EmptyLink, nil, "missing artifactTypes in verifier config", re.HideStackTrace)
	}

	return &conf, nil
}

// LoadPublicKey loads the public key from the keyRef
func loadPublicKey(keyRef string) (verifier signature.Verifier, err error) {
	keyPath := filepath.Clean(utils.ReplaceHomeShortcut(keyRef))
	raw, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	// PEM encoded file.
	ed, err := cosign.PemToECDSAKey(raw)
	if err != nil {
		return nil, errors.Wrap(err, "pem to ecdsa")
	}
	return signature.LoadECDSAVerifier(ed, crypto.SHA256)
}

// StaticLayerOpts builds the cosign options for static layer signatures
func staticLayerOpts(desc imgspec.Descriptor) ([]static.Option, error) {
	options := []static.Option{}
	options = append(options, static.WithAnnotations(desc.Annotations))
	cert := desc.Annotations[static.CertificateAnnotationKey]
	chain := desc.Annotations[static.ChainAnnotationKey]
	if cert != "" && chain != "" {
		options = append(options, static.WithCertChain([]byte(cert), []byte(chain)))
	}
	var rekorBundle bundle.RekorBundle
	if val, ok := desc.Annotations[static.BundleAnnotationKey]; ok {
		if err := json.Unmarshal([]byte(val), &rekorBundle); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal bundle from blob payload")
		}
		options = append(options, static.WithBundle(&rekorBundle))
	}

	return options, nil
}

// ErrorToVerifyResult returns a verifier result with the error message and isSuccess set to false
func errorToVerifyResult(name string, verifierType string, err error) verifier.VerifierResult {
	return verifier.VerifierResult{
		IsSuccess: false,
		Name:      name,
		Type:      verifierType,
		Message:   errors.Wrap(err, "cosign verification failed").Error(),
	}
}

// decodeASN1Signature decodes the ASN.1 signature to raw signature bytes
func decodeASN1Signature(sig []byte) ([]byte, error) {
	// Convert the ASN.1 Sequence to a concatenated r||s byte string
	// This logic is based from https://cs.opensource.google/go/go/+/refs/tags/go1.17.3:src/crypto/ecdsa/ecdsa.go;l=339
	var (
		r, s  = &big.Int{}, &big.Int{}
		inner cryptobyte.String
	)

	rawSigBytes := sig
	input := cryptobyte.String(sig)
	if input.ReadASN1(&inner, asn1.SEQUENCE) {
		// if ASN.1 sequence is found, parse r and s
		if !inner.ReadASN1Integer(r) {
			return nil, fmt.Errorf("failed parsing r")
		}
		if !inner.ReadASN1Integer(s) {
			return nil, fmt.Errorf("failed parsing s")
		}
		if !inner.Empty() {
			return nil, fmt.Errorf("failed parsing signature")
		}
		rawSigBytes = []byte{}
		rawSigBytes = append(rawSigBytes, r.Bytes()...)
		rawSigBytes = append(rawSigBytes, s.Bytes()...)
	}

	return rawSigBytes, nil
}

// getKeysMapsDefault returns the map of keys and cosign options for the reference
func getKeysMapsDefault(ctx context.Context, trustPolicies *TrustPolicies, reference string, namespace string) (map[PKKey]keymanagementprovider.PublicKey, cosign.CheckOpts, error) {
	// get the trust policy for the reference
	trustPolicy, err := trustPolicies.GetScopedPolicy(reference)
	if err != nil {
		return nil, cosign.CheckOpts{}, err
	}
	logger.GetLogger(ctx, logOpt).Debugf("selected trust policy %s for reference %s", trustPolicy.GetName(), reference)

	// get the map of keys for that reference
	keysMap, err := trustPolicy.GetKeys(namespace)
	if err != nil {
		return nil, cosign.CheckOpts{}, err
	}

	// get the cosign options for that trust policy
	cosignOpts, err := trustPolicy.GetCosignOpts(ctx)
	if err != nil {
		return nil, cosign.CheckOpts{}, err
	}

	return keysMap, cosignOpts, nil
}

// processAKVSignature processes the AKV signature and returns the hash type, signature and error
func processAKVSignature(sigEncoded string, staticSig oci.Signature, publicKey crypto.PublicKey, payloadBytes []byte, staticOpts []static.Option) (crypto.Hash, oci.Signature, error) {
	var hashType crypto.Hash
	switch keyType := publicKey.(type) {
	case *rsa.PublicKey:
		switch keyType.Size() {
		case 256:
			hashType = crypto.SHA256
		case 384:
			hashType = crypto.SHA384
		case 512:
			hashType = crypto.SHA512
		default:
			return crypto.SHA256, nil, fmt.Errorf("RSA key check: unsupported key size: %d", keyType.Size())
		}

		// TODO: remove section after fix for bug in cosign azure key vault implementation
		// tracking issue: https://github.com/sigstore/sigstore/issues/1384
		// summary: azure keyvault implementation ASN.1 encodes sig after online signing with keyvault
		// EC verifiers in cosign have built in ASN.1 decoding, but RSA verifiers do not
		base64DecodedBytes, err := base64.StdEncoding.DecodeString(sigEncoded)
		if err != nil {
			return crypto.SHA256, nil, fmt.Errorf("RSA key check: failed to decode base64 signature: %w", err)
		}
		// decode ASN.1 signature to raw signature if it is ASN.1 encoded
		decodedSigBytes, err := decodeASN1Signature(base64DecodedBytes)
		if err != nil {
			return crypto.SHA256, nil, fmt.Errorf("RSA key check: failed to decode ASN.1 signature: %w", err)
		}
		encodedBase64SigBytes := base64.StdEncoding.EncodeToString(decodedSigBytes)
		staticSig, err = static.NewSignature(payloadBytes, encodedBase64SigBytes, staticOpts...)
		if err != nil {
			return crypto.SHA256, nil, fmt.Errorf("RSA key check: failed to generate static signature: %w", err)
		}
	case *ecdsa.PublicKey:
		switch keyType.Curve {
		case elliptic.P256():
			hashType = crypto.SHA256
		case elliptic.P384():
			hashType = crypto.SHA384
		case elliptic.P521():
			hashType = crypto.SHA512
		default:
			return crypto.SHA256, nil, fmt.Errorf("ECDSA key check: unsupported key curve: %s", keyType.Params().Name)
		}
	default:
		return crypto.SHA256, nil, fmt.Errorf("unsupported public key type: %T", publicKey)
	}
	return hashType, staticSig, nil
}
