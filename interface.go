package memkv

import "bytes"

type KVDB interface {
	Put(key, value any) error
	Delete(key any) error
	Contains(key any) bool
	Get(key any) (value any, err error)
	Find(key any) (rkey, value any, err error)
	Num() int
}

type CompareFunc func(a, b any) int

type CopyKeyFunc func(v any) any

type CopyValueFunc func(v any) any

func BytesCompare(a, b any) int {
	return bytes.Compare(a.([]byte), b.([]byte))
}

func NoCopy(v any) any {
	return v
}

func BytesClone(v any) any {
	vbuf, ok := v.([]byte)
	if !ok {
		return v
	}
	buf := make([]byte, len(vbuf))
	copy(buf, vbuf)
	return buf
}
