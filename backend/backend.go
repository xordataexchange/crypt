package backend

type Response struct {
	Value []byte
	Error error
}

type Store interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Watch(key string, stop chan bool) <-chan *Response
}
