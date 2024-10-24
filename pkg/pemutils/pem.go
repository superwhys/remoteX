package pemutils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
	
	"github.com/pkg/errors"
)

func pemBlockForKey(priv interface{}) (*pem.Block, error) {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}, nil
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, err
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}, nil
	default:
		return nil, errors.New("unknown key type")
	}
}

// GenerateCertificate generates a PEM formatted key pair and self-signed certificate in memory.
func GenerateCertificate(commonName string, lifetimeDays int) (*pem.Block, *pem.Block, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("generate key: %w", err)
	}
	
	notBefore := time.Now().Truncate(24 * time.Hour)
	notAfter := notBefore.Add(time.Duration(lifetimeDays*24) * time.Hour)
	
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "generating serial number")
	}
	
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:         commonName,
			Organization:       []string{"RemoteX"},
			OrganizationalUnit: []string{"Automatically Generated"},
		},
		DNSNames:              []string{commonName},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}
	
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv)
	if err != nil {
		return nil, nil, fmt.Errorf("create certutils: %w", err)
	}
	
	certBlock := &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}
	keyBlock, err := pemBlockForKey(priv)
	if err != nil {
		return nil, nil, fmt.Errorf("save key: %w", err)
	}
	return certBlock, keyBlock, nil
}

func SaveCertificate(certFile, keyFile string, certBlock, keyBlock *pem.Block) error {
	certOut, err := os.Create(certFile)
	if err != nil {
		return fmt.Errorf("save certutils: %w", err)
	}
	if err = pem.Encode(certOut, certBlock); err != nil {
		return fmt.Errorf("save certutils: %w", err)
	}
	if err = certOut.Close(); err != nil {
		return fmt.Errorf("save certutils: %w", err)
	}
	
	keyOut, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("save key: %w", err)
	}
	if err = pem.Encode(keyOut, keyBlock); err != nil {
		return fmt.Errorf("save key: %w", err)
	}
	if err = keyOut.Close(); err != nil {
		return fmt.Errorf("save key: %w", err)
	}
	
	return nil
}
