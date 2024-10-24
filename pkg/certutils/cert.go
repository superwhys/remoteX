package certutils

import (
	"crypto/tls"
	"encoding/pem"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/superwhys/remoteX/config"
	"github.com/superwhys/remoteX/pkg/pemutils"
)

var (
	tlsDefaultCommonName   = "remoteX"
	deviceCertLifetimeDays = 20 * 365
)

func LoadOrGenerateCertificate(certConf *config.TlsConfig) (tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certConf.CertFile, certConf.KeyFile)
	if err != nil {
		plog.Debugf("create Certificate: cert: %s. key: %s.", certConf.CertFile, certConf.KeyFile)
		return GenerateCertificate(certConf.CertFile, certConf.KeyFile)
	}
	return cert, nil
}

func GenerateCertificate(certFile, keyFile string) (tls.Certificate, error) {
	return NewCertificate(certFile, keyFile, tlsDefaultCommonName, deviceCertLifetimeDays)
}

// NewCertificate generates and returns a new TLS certificate, saved to the given PEM files.
func NewCertificate(certFile, keyFile string, commonName string, lifetimeDays int) (tls.Certificate, error) {
	certBlock, keyBlock, err := pemutils.GenerateCertificate(commonName, lifetimeDays)
	if err != nil {
		return tls.Certificate{}, err
	}
	
	if err := pemutils.SaveCertificate(certFile, keyFile, certBlock, keyBlock); err != nil {
		return tls.Certificate{}, err
	}
	
	return tls.X509KeyPair(pem.EncodeToMemory(certBlock), pem.EncodeToMemory(keyBlock))
}
