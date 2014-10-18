package config

import (
	"io"

	"github.com/xordataexchange/crypt/backend"
	"github.com/xordataexchange/crypt/backend/etcd"
	"github.com/xordataexchange/crypt/encoding/secconf"
)

type ConfigManager interface {
	Get(key string) ([]byte, error)
}

type etcdConfigManager struct {
	keystore io.Reader
	store    backend.Store
}

func NewEtcdConfigManager(machines []string, keystore io.Reader) ConfigManager {
	return etcdConfigManager{keystore, etcd.New(machines)}
}

func (c etcdConfigManager) Get(key string) ([]byte, error) {
	value, err := c.store.Get(key)
	if err != nil {
		return nil, err
	}
	return secconf.Decode(value, c.keystore)
}
