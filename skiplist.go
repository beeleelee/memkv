package memkv

import (
	"fmt"
	"math/rand"
	"sync"
)

const _max_height = 12

type node struct {
	k     any
	v     any
	level []*node
}

type skipList struct {
	rnd       *rand.Rand
	mu        sync.RWMutex
	prevNode  [_max_height]*node
	heads     [_max_height]*node
	maxHeight int
	n         int
	comfunc   CompareFunc
	cpv       CopyValueFunc
}

func (s *skipList) randHeight() (h int) {
	const branching = 4
	h = 1
	for h < _max_height && s.rnd.Int()%branching == 0 {
		h++
	}
	return
}

func (s *skipList) findGE(key any, prev bool) (*node, bool) {
	if prev {
		// reset prevNode
		for i := 0; i < s.maxHeight; i++ {
			s.prevNode[i] = nil
		}
	}
	h := s.maxHeight - 1
	var next, nd *node

	for {
		if nd != nil {
			next = nd.level[h]
		} else {
			next = s.heads[h]
		}
		cmp := 1
		if next != nil {
			cmp = s.comfunc(next.k, key)
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

func (s *skipList) Put(key, value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if nd, exact := s.findGE(key, true); exact {
		nd.v = s.cpv(value)
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
		k:     key,
		v:     s.cpv(value),
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

func (s *skipList) Delete(key any) error {
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

func (s *skipList) Get(key any) (value any, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if nd, exact := s.findGE(key, false); exact {
		value = s.cpv(nd.v)
	} else {
		err = ErrNotFound
	}

	return
}

func (s *skipList) Num() int {
	return s.n
}

func (s *skipList) Contains(key any) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, exact := s.findGE(key, false); exact {
		return true
	}
	return false
}

func (s *skipList) Find(key any) (rkey, value any, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if nd, _ := s.findGE(key, false); nd != nil {
		rkey = nd.k
		value = s.cpv(nd.v)
		return
	}
	err = ErrNotFound
	return
}

func (s *skipList) print() {
	h := s.maxHeight - 1
	for ; h >= 0; h-- {
		fmt.Printf("level %d head -> ", h)
		node := s.heads[h]
		for {
			if node == nil {
				fmt.Print(" nil")
				break
			}
			fmt.Printf("%s, %s -> ", node.k, node.v)
			node = node.level[h]
		}
		fmt.Println()
	}
}

func newSkipList(cmpfunc CompareFunc, cpv CopyValueFunc) *skipList {
	return &skipList{
		rnd:       rand.New(rand.NewSource(0xdeadbeaf)),
		maxHeight: 1,
		comfunc:   cmpfunc,
		cpv:       cpv,
	}
}

func NewBytesKV() KVDB {
	return newSkipList(BytesCompare, BytesClone)
}
