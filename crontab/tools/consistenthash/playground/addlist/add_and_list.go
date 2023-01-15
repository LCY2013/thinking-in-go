package main

import (
	"fmt"
	consistent "github.com/LCY2013/thinking-in-go/crontab/tools/consistenthash"

	"github.com/cespare/xxhash"
)

type Member string

func (m Member) String() string {
	return string(m)
}

type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	return xxhash.Sum64(data)
}

func main() {
	members := []consistent.Member{}
	for i := 0; i < 8; i++ {
		member := Member(fmt.Sprintf("node%d.olricmq", i))
		members = append(members, member)
	}
	cfg := consistent.Config{
		PartitionCount:    71,
		ReplicationFactor: 20,
		Load:              1.25,
		Hasher:            hasher{},
	}

	c := consistent.New(members, cfg)
	owners := make(map[string]int)
	for partID := 0; partID < cfg.PartitionCount; partID++ {
		owner := c.GetPartitionOwner(partID)
		owners[owner.String()]++
	}
	fmt.Println("average load:", c.AverageLoad())
	fmt.Println("owners:", owners)
}
