package memkv

import (
	"bytes"
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
	heads     [_max_height]*node
	maxHeight int
	n         int
}

func (s *skipList) randHeight() (h int) {
	const branching = 4
	h = 1
	for h < _max_height && s.rnd.Int()%branching == 0 {
		h++
	}
	return
}

func (s *skipList) findGE(key []byte, prev bool) (nd *node, exact bool) {
	if prev {
		// reset prevNode
		for i := 0; i < s.maxHeight; i++ {
			s.prevNode[i] = nil
		}
	}
	h := s.maxHeight - 1
	var next *node

	for {
		if nd != nil {
			next = nd.level[h]
		} else {
			next = s.heads[h]
		}
		cmp := 1
		if next != nil {
			cmp = bytes.Compare(next.k, key)
		}
		if cmp < 0 {
			// keep searching in this list
			nd = next
		} else {
			if prev {
				s.prevNode[h] = nd
			} else if cmp == 0 {
				return next, true
			}
			if h == 0 {
				return next, cmp == 0
			}
			h--
		}
	}

}

func (s *skipList) Put(key, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if nd, exact := s.findGE(key, true); exact {
		nd.v = append(nd.v[:0], value...)
		return nil
	}

	h := s.randHeight()
	if h > s.maxHeight {
		// for i := s.maxHeight; i < h; i++ {
		// 	s.prevNode[i] = nil
		// }
		s.maxHeight = h
	}
	nd := &node{
		k:     append([]byte{}, key...),
		v:     append([]byte{}, value...),
		level: make([]*node, h),
	}
	for i, n := range s.prevNode[:h] {
		if n == nil {
			head := s.heads[i]
			s.heads[i] = nd
			if head != nil {
				nd.level[i] = head
			}
		} else {
			nd.level[i] = n.level[i]
			n.level[i] = nd
		}
	}
	s.n++
	return nil

}

func (s *skipList) Delete(key []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	nd, exact := s.findGE(key, true)
	if !exact {
		return ErrNotFound
	}

	h := len(nd.level) - 1
	for i, n := range s.prevNode[:h] {
		if n != nil {
			n.level[i] = nd.level[i]
		} else {
			s.heads[i] = nd.level[i]
		}
	}
	s.n--
	return nil
}

func (s *skipList) Get(key []byte) (value []byte, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if nd, exact := s.findGE(key, false); exact {
		value = append([]byte{}, nd.v...)
	} else {
		err = ErrNotFound
	}

	return
}

func New() KVDB {
	return &skipList{
		rnd:       rand.New(rand.NewSource(0xabacabed)),
		maxHeight: 1,
	}

}
