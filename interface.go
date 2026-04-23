package memkv

type KVDB interface {
	Put(key, value []byte) error
	Delete(key []byte) error
	// Contains(key []byte) bool
	Get(key []byte) (value []byte, err error)
	// Find(key []byte) (rkey, value []byte, err error)
	Num() int
	DumpKeys() [][]byte
	Dump() ([][]byte, [][]byte)
}
