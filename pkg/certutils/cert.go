package certutils

import (
	"crypto/tls"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/superwhys/remoteX/config"
	tlsutils "github.com/superwhys/remoteX/internal/tls"
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
	return tlsutils.NewCertificate(certFile, keyFile, tlsDefaultCommonName, deviceCertLifetimeDays)
}
