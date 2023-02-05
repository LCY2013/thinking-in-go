package consistenthash

import (
	"hash/crc32"
	"sort"
)

// HashFunc defines function to generate hash code
type HashFunc func(data []byte) uint32

// NodeMap stores nodes and you can pick node from NodeMap
type NodeMap struct {
	hashFunc    HashFunc
	hashNodes   []int // sorted
	hashMapNode map[int]string
}

// NewNodeMap creates a new NodeMap
func NewNodeMap(fn HashFunc) *NodeMap {
	nm := &NodeMap{
		hashFunc:    fn,
		hashMapNode: make(map[int]string),
	}

	if fn == nil {
		nm.hashFunc = crc32.ChecksumIEEE
	}

	return nm
}

// IsEmpty returns if there is no node in NodeMap
func (nm *NodeMap) IsEmpty() bool {
	return len(nm.hashMapNode) == 0
}

// AddNode add the given nodes into consistent hash circle
func (nm *NodeMap) AddNode(nodes ...string) {
	for _, node := range nodes {
		if node == "" {
			continue
		}

		hash := int(nm.hashFunc([]byte(node)))
		nm.hashNodes = append(nm.hashNodes, hash)
		nm.hashMapNode[hash] = node
	}

	sort.Ints(nm.hashNodes)
}

// PickNode gets the closest item in the hash to the provided key.
func (nm *NodeMap) PickNode(key string) string {
	if nm.IsEmpty() {
		return ""
	}

	hash := int(nm.hashFunc([]byte(key)))

	// Binary search for appropriate replica.
	idx := sort.Search(len(nm.hashNodes), func(i int) bool {
		return nm.hashNodes[i] >= hash
	})

	// Means we have cycled back to the first replica.
	if idx == len(nm.hashNodes) {
		idx = 0
	}

	return nm.hashMapNode[nm.hashNodes[idx]]
}
