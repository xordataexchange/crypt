package consul

import (
	"strings"

	"github.com/xordataexchange/crypt/backend"

	"github.com/armon/consul-api"
)

type Client struct {
	client *consulapi.KV
}

func New(machines []string) (*Client, error) {
	conf := consulapi.DefaultConfig()
	if len(machines) > 0 {
		conf.Address = machines[0]
	}
	client, err := consulapi.NewClient(conf)
	if err != nil {
		return nil, err
	}
	return &Client{client.KV()}, nil
}

func (c *Client) Get(key string) ([]byte, error) {
	kv, _, err := c.client.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return kv.Value, nil
}

func (c *Client) Set(key string, value []byte) error {
	key = strings.TrimPrefix(key, "/")
	kv := &consulapi.KVPair{
		Key:   key,
		Value: value,
	}
	_, err := c.client.Put(kv, nil)
	return err
}

func (c *Client) Watch(key string, stop chan bool) <-chan *backend.Response {
	respChan := make(chan *backend.Response, 0)
	return respChan
}
