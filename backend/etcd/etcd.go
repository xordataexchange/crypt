package etcd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/xordataexchange/crypt/backend"

	goetcd "github.com/coreos/etcd/client"
)

type Client struct {
	client    goetcd.Client
	keysAPI   goetcd.KeysAPI
	waitIndex uint64
}

func New(machines []string) (*Client, error) {
	newClient, err := goetcd.New(goetcd.Config{
		Endpoints: machines,
	})
	if err != nil {
		return nil, fmt.Errorf("creating new etcd client for crypt.backend.Client: %v", err)
	}
	keysAPI := goetcd.NewKeysAPI(newClient)
	return &Client{client: newClient, keysAPI: keysAPI, waitIndex: 0}, nil
}

func (c *Client) Get(key string) ([]byte, error) {
	resp, err := c.keysAPI.Get(context.TODO(), key, nil)
	if err != nil {
		return nil, err
	}
	return []byte(resp.Node.Value), nil
}

func addKVPairs(node *goetcd.Node, list backend.KVPairs) backend.KVPairs {
	if node.Dir {
		for _, n := range node.Nodes {
			list = addKVPairs(n, list)
		}
		return list
	}
	return append(list, &backend.KVPair{Key: node.Key, Value: []byte(node.Value)})
}

func (c *Client) List(key string) (backend.KVPairs, error) {
	resp, err := c.keysAPI.Get(context.TODO(), key, nil)
	if err != nil {
		return nil, err
	}
	if !resp.Node.Dir {
		return nil, errors.New("key is not a directory")
	}
	list := addKVPairs(resp.Node, nil)
	return list, nil
}

func (c *Client) Set(key string, value []byte) error {
	_, err := c.keysAPI.Set(context.TODO(), key, string(value), nil)
	return err
}

func (c *Client) Watch(key string, stop chan bool) <-chan *backend.Response {
	respChan := make(chan *backend.Response, 0)
	go func() {
		watcher := c.keysAPI.Watcher(key, nil)
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			<-stop
			cancel()
		}()
		for {
			var resp *goetcd.Response
			var err error
			// if c.waitIndex == 0 {
			// 	resp, err = c.client.Get(key, false, false)
			// 	if err != nil {
			// 		respChan <- &backend.Response{nil, err}
			// 		time.Sleep(time.Second * 5)
			// 		continue
			// 	}
			// 	c.waitIndex = resp.EtcdIndex
			// 	respChan <- &backend.Response{[]byte(resp.Node.Value), nil}
			// }
			// resp, err = c.client.Watch(key, c.waitIndex+1, false, nil, stop)
			resp, err = watcher.Next(ctx)
			if err != nil {
				respChan <- &backend.Response{nil, err}
				time.Sleep(time.Second * 5)
				continue
			}
			c.waitIndex = resp.Node.ModifiedIndex
			respChan <- &backend.Response{[]byte(resp.Node.Value), nil}
		}
	}()
	return respChan
}
