package mock

import (
	"errors"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/bketelsen/crypt/backend"
)

var (
	mockedStore map[string][]byte

	once = sync.Once{}
	lock = sync.RWMutex{}
)

type Client struct{}

func New(machines []string) (*Client, error) {
	once.Do(func() {
		mockedStore = make(map[string][]byte, 2)
	})
	return &Client{}, nil
}

func (c *Client) Get(key string) ([]byte, error) {
	lock.RLock()
	defer lock.RUnlock()

	if v, ok := mockedStore[key]; ok {
		return v, nil
	}
	err := errors.New("Could not find key: " + key)
	return nil, err
}

func (c *Client) List(key string) (backend.KVPairs, error) {
	lock.RLock()
	defer lock.RUnlock()

	var list backend.KVPairs
	dir := path.Clean(key) + "/"
	for k, v := range mockedStore {
		if strings.HasPrefix(k, dir) {
			list = append(list, &backend.KVPair{Key: k, Value: v})
		}
	}
	return list, nil
}

func (c *Client) Set(key string, value []byte) error {
	lock.Lock()
	defer lock.Unlock()

	mockedStore[key] = value
	return nil
}

func (c *Client) Watch(key string, stop chan bool) <-chan *backend.Response {
	lock.RLock()
	defer lock.RUnlock()

	respChan := make(chan *backend.Response, 0)
	go func() {
		for {
			b, err := c.Get(key)
			if err != nil {
				respChan <- &backend.Response{nil, err}
				time.Sleep(time.Second * 5)
				continue
			}
			respChan <- &backend.Response{b, nil}
		}
	}()
	return respChan
}
