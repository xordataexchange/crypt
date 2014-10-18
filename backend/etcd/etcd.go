package etcd

import (
	goetcd "github.com/coreos/go-etcd/etcd"
)

type Client struct {
	client *goetcd.Client
}

func New(machines []string) *Client {
	return &Client{goetcd.NewClient(machines)}
}

func (c *Client) Get(key string) ([]byte, error) {
	resp, err := c.client.Get(key, false, false)
	if err != nil {
		return nil, err
	}
	return []byte(resp.Node.Value), nil
}

func (c *Client) Set(key string, value []byte) error {
	_, err := c.client.Set(key, string(value), 0)
	return err
}
