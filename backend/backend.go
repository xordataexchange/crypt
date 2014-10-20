// Package backend provides the K/V store interface for crypt backends.
package backend

// Response represents a response from a backend store.
type Response struct {
	Value []byte
	Error error
}

// A Store is a K/V store backend that retrieves and sets, and monitors
// data in a K/V store.
type Store interface {
	// Get retrieves a value from a K/V store for the provided key.
	Get(key string) ([]byte, error)

	// Set sets the provided key to value.
	Set(key string, value []byte) error

	// Watch monitors a K/V store for changes to key.
	Watch(key string, stop chan bool) <-chan *Response
}
