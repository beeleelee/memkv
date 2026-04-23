package memkv

import (
	"math/rand"
	"sync"
)

const _max_height = 12

type node struct {
	k     []byte
	v     []byte
	level []*node
}

type skipList struct {
	rnd       *rand.Rand
	mu        sync.RWMutex
	prevNode  [_max_height]*node
	maxHeight int
	n         int
}

func (s *skipList) Put(key, value []byte) error {
	return nil
}

func (s *skipList) Delete(key []byte) error {
	return nil
}

func (s *skipList) Get(key []byte) (value []byte, err error) {
	return nil, nil
}

func New() KVDB {
	return &skipList{
		rnd:       rand.New(rand.NewSource(0xabacabed)),
		maxHeight: 1,
	}

}
