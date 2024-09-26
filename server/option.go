package server

import (
	"crypto/tls"
	"net/url"
	
	"github.com/superwhys/remoteX/config"
	"github.com/superwhys/remoteX/domain/node"
	"github.com/superwhys/remoteX/pkg/certutils"
	"github.com/superwhys/remoteX/pkg/common"
)

type Option struct {
	Conf      *config.Config
	Cert      tls.Certificate
	TlsConfig *tls.Config
	Local     *node.Node
}

func InitOption(conf *config.Config) (opt *Option, err error) {
	opt = &Option{Conf: conf}
	opt.Cert, err = certutils.LoadOrGenerateCertificate(conf.Tls)
	if err != nil {
		return nil, err
	}
	
	opt.TlsConfig = &tls.Config{
		ServerName:             "remoteX",
		Certificates:           []tls.Certificate{opt.Cert},
		MinVersion:             tls.VersionTLS12,
		ClientAuth:             tls.RequestClientCert,
		InsecureSkipVerify:     true,
		SessionTicketsDisabled: true,
	}
	
	opt.Local = conf.LocalNode
	opt.Local.NodeId = common.NewNodeID(opt.Cert.Certificate[0])
	opt.Local.Status = node.NodeStatusOnline
	return
}

func (o *Option) LocalUrl() *url.URL {
	return o.Local.URL()
}
