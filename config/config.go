package config

import (
	"io"

	"github.com/xordataexchange/crypt/backend"
	"github.com/xordataexchange/crypt/backend/consul"
	"github.com/xordataexchange/crypt/backend/etcd"
	"github.com/xordataexchange/crypt/encoding/secconf"
)

type ConfigManager interface {
	Get(key string) ([]byte, error)
}

type configManager struct {
	keystore io.Reader
	store    backend.Store
}

func NewEtcdConfigManager(machines []string, keystore io.Reader) (ConfigManager, error) {
	store, err := etcd.New(machines)
	if err != nil {
		return nil, err
	}
	return configManager{keystore, store}, nil
}

func NewConsulConfigManager(machines []string, keystore io.Reader) (ConfigManager, error) {
	store, err := consul.New(machines)
	if err != nil {
		return nil, err
	}
	return configManager{keystore, store}, nil
}

func (c configManager) Get(key string) ([]byte, error) {
	value, err := c.store.Get(key)
	if err != nil {
		return nil, err
	}
	return secconf.Decode(value, c.keystore)
}
