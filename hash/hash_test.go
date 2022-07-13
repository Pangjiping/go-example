package hash

import "testing"

func Test_HashVnode(t *testing.T) {
	t.Log("start test consistent hash with virtual nodes...")
	ChVnode()
}

func Test_CommonHash(t *testing.T) {
	t.Log("start test common hash...")
	ComHash()
}

func Test_ConsistentHash(t *testing.T) {
	t.Log("start test consistent hash...")
	ConsistentHash()
}
