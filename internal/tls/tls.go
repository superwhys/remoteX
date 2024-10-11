package tls

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
	
	"github.com/pkg/errors"
)

var (
	tlsHandshakeTimeout = 10 * time.Second
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

// generateCertificate generates a PEM formatted key pair and self-signed certificate in memory.
func generateCertificate(commonName string, lifetimeDays int) (*pem.Block, *pem.Block, error) {
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
			Organization:       []string{"Syncthing"},
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

// NewCertificate generates and returns a new TLS certificate, saved to the given PEM files.
func NewCertificate(certFile, keyFile string, commonName string, lifetimeDays int) (tls.Certificate, error) {
	certBlock, keyBlock, err := generateCertificate(commonName, lifetimeDays)
	if err != nil {
		return tls.Certificate{}, err
	}
	
	certOut, err := os.Create(certFile)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("save certutils: %w", err)
	}
	if err = pem.Encode(certOut, certBlock); err != nil {
		return tls.Certificate{}, fmt.Errorf("save certutils: %w", err)
	}
	if err = certOut.Close(); err != nil {
		return tls.Certificate{}, fmt.Errorf("save certutils: %w", err)
	}
	
	keyOut, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("save key: %w", err)
	}
	if err = pem.Encode(keyOut, keyBlock); err != nil {
		return tls.Certificate{}, fmt.Errorf("save key: %w", err)
	}
	if err = keyOut.Close(); err != nil {
		return tls.Certificate{}, fmt.Errorf("save key: %w", err)
	}
	
	return tls.X509KeyPair(pem.EncodeToMemory(certBlock), pem.EncodeToMemory(keyBlock))
}

func TlsTimedHandshake(tc *tls.Conn) error {
	tc.SetDeadline(time.Now().Add(tlsHandshakeTimeout))
	defer tc.SetDeadline(time.Time{})
	
	return tc.Handshake()
}

func LocalHost(tc *tls.Conn) string {
	return fmt.Sprintf("%s://%s", tc.LocalAddr().Network(), tc.LocalAddr().String())
}

func RemoteHost(tc *tls.Conn) string {
	return fmt.Sprintf("%s://%s", tc.RemoteAddr().Network(), tc.RemoteAddr().String())
}
