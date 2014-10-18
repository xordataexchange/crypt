package backend

type Store interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}
