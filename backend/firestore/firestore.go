package firestore

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/bketelsen/crypt/backend"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

type Client struct {
	client *firestore.Client
}

type data struct {
	Data []byte `firestore:"data"`
}

func New(machines []string) (*Client, error) {
	if len(machines) == 0 {
		return nil, errors.New("project should be defined")
	}

	opts := []option.ClientOption{
		option.WithGRPCDialOption(grpc.WithBlock()),
	}
	c, err := firestore.NewClient(context.TODO(), machines[0], opts...)
	if err != nil {
		return nil, err
	}
	return &Client{c}, nil
}

func (c *Client) Get(path string) ([]byte, error) {
	return c.GetWithContext(context.TODO(), path)
}

func (c *Client) GetWithContext(ctx context.Context, path string) ([]byte, error) {
	snap, err := c.client.Doc(path).Get(ctx)
	if err != nil {
		return nil, err
	}

	d := &data{}
	err = snap.DataTo(&d)
	if err != nil {
		return nil, err
	}
	return d.Data, nil
}

func (c *Client) List(collection string) (backend.KVPairs, error) {
	return c.ListWithContext(context.TODO(), collection)
}

func (c *Client) ListWithContext(ctx context.Context, collection string) (backend.KVPairs, error) {
	res := backend.KVPairs{}
	it := c.client.Collection(collection).Documents(ctx)
	for {
		doc, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		d := &data{}
		err = doc.DataTo(&d)
		if err != nil {
			return nil, err
		}
		res = append(res, &backend.KVPair{
			Key:   doc.Ref.ID,
			Value: d.Data,
		})
	}
	return res, nil
}

func (c *Client) Set(path string, value []byte) error {
	return c.SetWithContext(context.TODO(), path, value)
}

func (c *Client) SetWithContext(ctx context.Context, path string, value []byte) error {
	_, err := c.client.Doc(path).Set(ctx, &data{value})
	return err
}

func (c *Client) Watch(path string, stop chan bool) <-chan *backend.Response {
	return c.WatchWithContext(context.TODO(), path, stop)
}

func (c *Client) WatchWithContext(ctx context.Context, path string, stop chan bool) <-chan *backend.Response {
	ch := make(chan *backend.Response, 0)
	t := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-t.C:
				v := &data{}
				snap, err := c.client.Doc(path).Get(ctx)
				if err == nil {
					err = snap.DataTo(&v)
				}
				ch <- &backend.Response{v.Data, err}
				if err != nil {
					time.Sleep(time.Second * 5)
				}
			case <-stop:
				close(ch)
				return
			}
		}
	}()
	return ch
}
