package backend

type Backend interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}
