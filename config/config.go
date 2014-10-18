package config

import (
	"io"

	"github.com/xordataexchange/crypt/backend"
	"github.com/xordataexchange/crypt/backend/consul"
	"github.com/xordataexchange/crypt/backend/etcd"
	"github.com/xordataexchange/crypt/encoding/secconf"
)

// A ConfigManager retrieves and decrypts configuration from a key/value store. 
type ConfigManager interface {
	Get(key string) ([]byte, error)
}

type configManager struct {
	keystore io.Reader
	store    backend.Store
}

// NewEtcdConfigManager returns a new ConfigManager backed by etcd.
func NewEtcdConfigManager(machines []string, keystore io.Reader) (ConfigManager, error) {
	store, err := etcd.New(machines)
	if err != nil {
		return nil, err
	}
	return configManager{keystore, store}, nil
}

// NewConsulConfigManager returns a new ConfigManager backed by consul.
func NewConsulConfigManager(machines []string, keystore io.Reader) (ConfigManager, error) {
	store, err := consul.New(machines)
	if err != nil {
		return nil, err
	}
	return configManager{keystore, store}, nil
}

// Get retrieves and decodes a secconf value stored at key.
func (c configManager) Get(key string) ([]byte, error) {
	value, err := c.store.Get(key)
	if err != nil {
		return nil, err
	}
	return secconf.Decode(value, c.keystore)
}
