package certificateprovider

import (
	"context"

	"github.com/deislabs/ratify/pkg/certificateprovider/akv"
	"github.com/deislabs/ratify/pkg/certificateprovider/akv/types"
	"github.com/sirupsen/logrus"
)

type CertificateProvider interface {
	GetSecretsStoreObjectContent(ctx context.Context, attrib map[string]string)
}

// CRD manager call this method fetch certificate in memory
func SetCert(ctx context.Context, attrib map[string]string) {
	// To implement
}

// Verifier call this method to get validation certificate
func GetCert(ctx context.Context) ([]types.SecretFile, error) {
	// TO implement
	// TODO: populate the map with keyvault info
	attrib := map[string]string{}
	attrib["keyvaultName"] = "notarycerts"
	attrib["clientID"] = "1c7ac023-5bf6-4916-83f2-96dd203e35a3"
	attrib["cloudName"] = "AzurePublicCloud"
	attrib["tenantID"] = "72f988bf-86f1-41af-91ab-2d7cd011db47"

	attrib["objects"] = "array:\n- |\n  certificateName: wabbit-networks-io	\n  certificateAlias: \"testCert\"\n  certificateVersion: 97a1545d893344079ce57699c8810590 \n  certificateVersionHistory: 0\n"
	files, _ := akv.GetSecretsStoreObjectContent(ctx, attrib)
	logrus.Infof(string(files[0].Content))
	return nil, nil
}
