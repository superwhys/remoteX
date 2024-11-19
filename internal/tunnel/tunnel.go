// File:		tunnel.go
// Created by:	Hoven
// Created on:	2024-11-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package tunnel

import (
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/superwhys/sshtunnel"
)

type SshConfig struct {
	HostName string
	User     string
}

type tunnel struct {
	cancel func()
	t      *sshtunnel.SshTunnel
}

type SshTunnel struct {
	lock         sync.RWMutex
	sshConfig    *SshConfig
	identityFile *embed.FS
	manager      map[string]*tunnel
}

func NewSshTunnel(sshConfig *SshConfig, identityFile *embed.FS) *SshTunnel {
	return &SshTunnel{sshConfig: sshConfig, identityFile: identityFile}
}

func (t *SshTunnel) registerTunnel(key string, st *sshtunnel.SshTunnel, cancel func()) {
	t.lock.Lock()
	t.manager[key] = &tunnel{t: st, cancel: cancel}
	t.lock.Unlock()
}

func (t *SshTunnel) CloseTunnel(key string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if tunnel, ok := t.manager[key]; ok {
		tunnel.cancel()
		tunnel.t.Close()
		delete(t.manager, key)
	} else {
		return
	}
}

type doFunc func(ctx context.Context, localAddr, remoteAddr string) error

func (t *SshTunnel) doTunnel(ctx context.Context, tunnel *sshtunnel.SshTunnel, key string, localAddr, remoteAddr string, fn doFunc) error {
	t.lock.RLock()
	if _, exists := t.manager[key]; exists {
		t.lock.RUnlock()
		return errors.New("tunnel already exists")
	}
	t.lock.RUnlock()

	ctx, cancel := context.WithCancel(ctx)
	if err := fn(ctx, localAddr, remoteAddr); err != nil {
		return errors.Wrap(err, "doTunnel")
	}

	t.registerTunnel(key, tunnel, cancel)
	return nil
}

func (t *SshTunnel) Forward(ctx context.Context, localAddr, remoteAddr string) (string, error) {
	key := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("forward:%s->%s", localAddr, remoteAddr)))

	tunnel := sshtunnel.NewTunnelWithEmbed(t.identityFile, &sshtunnel.SshConfig{
		User:     t.sshConfig.User,
		HostName: t.sshConfig.HostName,
	})

	if err := t.doTunnel(ctx, tunnel, key, localAddr, remoteAddr, tunnel.Forward); err != nil {
		return "", errors.Wrap(err, "forward tunnel")
	}

	return key, nil
}

func (t *SshTunnel) Reverse(ctx context.Context, localAddr, remoteAddr string) (string, error) {
	key := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("reverse:%s->%s", localAddr, remoteAddr)))

	tunnel := sshtunnel.NewTunnelWithEmbed(t.identityFile, &sshtunnel.SshConfig{
		User:     t.sshConfig.User,
		HostName: t.sshConfig.HostName,
	})

	if err := t.doTunnel(ctx, tunnel, key, localAddr, remoteAddr, tunnel.Reverse); err != nil {
		return "", errors.Wrap(err, "forward tunnel")
	}

	return key, nil

}
