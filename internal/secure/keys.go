package secure

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// LoadPrivateKey reads a PEM file and returns an ECDSA Private Key
func LoadPrivateKey(path string) (*ecdsa.PrivateKey, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read private key file: %w", err)
	}

	block, _ := pem.Decode(content)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from private key")
	}

	return x509.ParseECPrivateKey(block.Bytes)
}

// LoadPublicKey reads a PEM file and returns an ECDSA Public Key
func LoadPublicKey(path string) (*ecdsa.PublicKey, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read public key file: %w", err)
	}

	block, _ := pem.Decode(content)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	ecdsaPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an ECDSA public key")
	}

	return ecdsaPub, nil
}
