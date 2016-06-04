package zookeeper

import (
	"errors"
	"fmt"
	zk "github.com/samuel/go-zookeeper/zk"
	"github.com/xordataexchange/crypt/backend"
	"strings"
	"time"
)

type Client struct {
	client    *zk.Conn
	waitIndex uint64
}

func New(machines []string) (*Client, error) {
	zkclient, _, err := zk.Connect(machines, time.Second)
	if err != nil {
		return nil, err
	}
	return &Client{zkclient, 0}, nil
}

func (c *Client) Get(key string) ([]byte, error) {
	resp, _, err := c.client.Get(key)
	if err != nil {
		return nil, err
	}
	return []byte(resp), nil
}

func nodeWalk(prefix string, c *Client, vars map[string]string) error {
	l, stat, err := c.client.Children(prefix)
	if err != nil {
		return err
	}

	if stat.NumChildren == 0 {
		b, _, err := c.client.Get(prefix)
		if err != nil {
			return err
		}
		vars[prefix] = string(b)

	} else {
		for _, key := range l {
			s := prefix + "/" + key
			_, stat, err := c.client.Exists(s)
			if err != nil {
				return err
			}
			if stat.NumChildren == 0 {
				b, _, err := c.client.Get(s)
				if err != nil {
					return err
				}
				vars[s] = string(b)
			} else {
				nodeWalk(s, c, vars)
			}
		}
	}
	return nil
}

func (c *Client) GetValues(key string, keys []string) (map[string]string, error) {
	vars := make(map[string]string)
	for _, v := range keys {
		v = fmt.Sprintf("%s/%s", key, v)
		v = strings.Replace(v, "/*", "", -1)
		_, _, err := c.client.Exists(v)
		if err != nil {
			return vars, err
		}
		if v == "/" {
			v = ""
		}
		err = nodeWalk(v, c, vars)
		if err != nil {
			return vars, err
		}
	}
	return vars, nil
}

func (c *Client) List(key string) (backend.KVPairs, error) {
	var list backend.KVPairs
	resp, stat, err := c.client.Children(key)
	if err != nil {
		return nil, err
	}

	if stat.NumChildren == 0 {
		return list, nil
	}

	entries, err := c.GetValues(key, resp)
	if err != nil {
		return nil, err
	}

	for k, v := range entries {
		list = append(list, &backend.KVPair{Key: k, Value: []byte(v)})
	}
	return list, nil
}

func (c *Client) createParents(key string) error {
	flags := int32(0)
	acl := zk.WorldACL(zk.PermAll)

	if key[0] != '/' {
		return errors.New("Invalid path")
	}

	payload := []byte("")
	pathString := ""
	pathNodes := strings.Split(key, "/")
	for i := 1; i < len(pathNodes); i++ {
		pathString += "/" + pathNodes[i]
		_, err := c.client.Create(pathString, payload, flags, acl)
		// not being able to create the node because it exists or not having
		// sufficient rights is not an issue. It is ok for the node to already
		// exist and/or us to only have read rights
		if err != nil && err != zk.ErrNodeExists && err != zk.ErrNoAuth {
			return err
		}
	}
	return nil
}

func (c *Client) Set(key string, value []byte) error {
	err := c.createParents(key)
	if err != nil {
		return err
	}
	_, err = c.client.Set(key, []byte(value), -1)
	return err
}

func (c *Client) Watch(key string, stop chan bool) <-chan *backend.Response {
	respChan := make(chan *backend.Response, 0)
	go func() {
		for {
			resp, _, watch, err := c.client.GetW(key)
			if err != nil {
				respChan <- &backend.Response{nil, err}
				time.Sleep(time.Second * 5)
			}

			select {
			case e := <-watch:
				if e.Type == zk.EventNodeDataChanged {
					resp, _, err = c.client.Get(key)
					if err != nil {
						respChan <- &backend.Response{nil, err}
					}
					c.waitIndex = 0
					respChan <- &backend.Response{[]byte(resp), nil}
				}
			}
		}
	}()
	return respChan
}
