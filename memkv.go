package memkv

type node struct {
	k     []byte
	v     []byte
	level []*node
}
