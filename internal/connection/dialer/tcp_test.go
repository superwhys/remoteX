// File:		tcp_test.go
// Created by:	Hoven
// Created on:	2024-09-26
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package dialer

import (
	"context"
	"crypto/tls"
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/superwhys/remoteX/config"
	"github.com/superwhys/remoteX/domain/connection"
	"github.com/superwhys/remoteX/domain/node"
	"github.com/superwhys/remoteX/pkg/certutils"
	"github.com/superwhys/remoteX/pkg/common"
	"github.com/superwhys/remoteX/pkg/protocol"
)

func TestTcpDialer(t *testing.T) {
	initTcpDialer()
	t.Run("tcpDial", func(t *testing.T) {
		local := &node.Node{
			Name: "testDial",
			Address: protocol.Address{
				IpAddress: "127.0.0.1",
				Port:      21011,
				Schema:    "tcp",
			},
		}
		
		cert, err := certutils.LoadOrGenerateCertificate(&config.TlsConfig{
			CertFile: "/Users/yong/.ssh/cert.pem",
			KeyFile:  "/Users/yong/.ssh/key.pem",
		})
		
		assert.Nil(t, err)
		tlsConf := &tls.Config{
			ServerName:             "remoteX",
			Certificates:           []tls.Certificate{cert},
			MinVersion:             tls.VersionTLS12,
			ClientAuth:             tls.RequestClientCert,
			InsecureSkipVerify:     true,
			SessionTicketsDisabled: true,
		}
		local.NodeId = common.NewNodeID(cert.Certificate[0])
		
		target := &node.Node{
			Name: "testDialTarget",
			Address: protocol.Address{
				IpAddress: "localhost",
				Port:      21012,
				Schema:    "tcp",
			},
		}
		fac, err := connection.GetDialerFactory(target.URL())
		assert.Nil(t, err)
		tc, err := fac.New(local.URL(), tlsConf).Dial(context.Background(), target.URL())
		assert.Nil(t, err)
		
		s, err := tc.OpenStream()
		assert.Nil(t, err)
		err = s.WriteMessage(local)
		assert.Nil(t, err)
		
		remote := new(node.Node)
		err = s.ReadMessage(remote)
		assert.Nil(t, err)
		t.Log(remote)
	})
}
