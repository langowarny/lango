package security

import (
	"fmt"

	"github.com/langoai/lango/internal/config"
)

// KMSProviderName identifies a supported KMS backend.
type KMSProviderName string

const (
	KMSProviderAWS    KMSProviderName = "aws-kms"
	KMSProviderGCP    KMSProviderName = "gcp-kms"
	KMSProviderAzure  KMSProviderName = "azure-kv"
	KMSProviderPKCS11 KMSProviderName = "pkcs11"
)

// Valid reports whether n is a recognised KMS provider name.
func (n KMSProviderName) Valid() bool {
	switch n {
	case KMSProviderAWS, KMSProviderGCP, KMSProviderAzure, KMSProviderPKCS11:
		return true
	}
	return false
}

// NewKMSProvider creates a CryptoProvider for the named KMS backend.
// Supported providers: "aws-kms", "gcp-kms", "azure-kv", "pkcs11".
// Build tags control which providers are compiled in; uncompiled providers
// return a descriptive error.
func NewKMSProvider(providerName KMSProviderName, kmsConfig config.KMSConfig) (CryptoProvider, error) {
	switch providerName {
	case KMSProviderAWS:
		return newAWSKMSProvider(kmsConfig)
	case KMSProviderGCP:
		return newGCPKMSProvider(kmsConfig)
	case KMSProviderAzure:
		return newAzureKVProvider(kmsConfig)
	case KMSProviderPKCS11:
		return newPKCS11Provider(kmsConfig)
	default:
		return nil, fmt.Errorf("unknown KMS provider: %q (supported: aws-kms, gcp-kms, azure-kv, pkcs11)", providerName)
	}
}
