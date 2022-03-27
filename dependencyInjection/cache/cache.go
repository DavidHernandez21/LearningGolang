package cache

import "errors"

type contextKeyType int

const (
	ExecutorContextKey contextKeyType = iota
)

var (
	errNotFoundRecord = errors.New("could not find record")
)

type Executor interface {
	Set(key string, data []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}

type InMemory struct {
	data map[string][]byte
}

func (i *InMemory) Get(key string) ([]byte, error) {
	value, ok := i.data[key]

	if !ok {

		return nil, errNotFoundRecord

	}

	return value, nil
}

func (i *InMemory) Set(key string, data []byte) error {
	i.data[key] = data
	return nil
}

func (i *InMemory) Delete(key string) error {
	delete(i.data, key)
	return nil
}

func NewInMemoryCache( /*config*/ ) *InMemory {
	return &InMemory{
		data: make(map[string][]byte, 100),
	}
}
