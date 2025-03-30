package storage

import (
	"fmt"
	"sync/atomic"
	"time"

	tarantool "github.com/tarantool/go-tarantool"
)

// Storage - структура хранилища
type Storage struct {
	conn   *tarantool.Connection
	lastID atomic.Uint64
}

// NewStorage - создаёт соединение с Tarantool и возвращает Storage
func NewStorage(host string, port string) (*Storage, error) {
	opts := tarantool.Opts{
		Timeout: time.Second,
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	conn, err := tarantool.Connect(addr, opts)
	if err != nil {
		return nil, err
	}

	return &Storage{
		conn:   conn,
		lastID: atomic.Uint64{},
	}, nil
}
