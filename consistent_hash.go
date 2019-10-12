package hconsistent

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"sort"
	"sync"
)

const replicationFactor = 1 << 4

type ConsistentHash struct {
	sync.RWMutex
	load  map[string]struct{}
	hosts map[uint64]string
	index []uint64
}

func NewConsistentHash() *ConsistentHash {
	return &ConsistentHash{load: map[string]struct{}{}, hosts: map[uint64]string{}, index: []uint64{}}
}

func (c *ConsistentHash) Add(host string) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.load[host]; ok {
		return
	}
	c.load[host] = struct{}{}

	for i := 0; i < replicationFactor; i++ {
		index := c.hash(fmt.Sprintf("%s-%d", host, i))

		for {
			if _, ok := c.hosts[index]; ok {
				index++
			} else {
				break
			}
		}

		c.hosts[index] = host
		c.index = append(c.index, index)
	}

	sort.Slice(c.index, func(i, j int) bool {
		return c.index[i] < c.index[j]
	})
}

func (c *ConsistentHash) Get(key string) string {
	c.RLock()
	defer c.RUnlock()

	if len(c.hosts) == 0 {
		panic("consistent no hosts")
	}

	index := c.search(c.hash(key))

	return c.hosts[c.index[index]]
}

func (c *ConsistentHash) search(sum uint64) int {
	index := sort.Search(len(c.index), func(i int) bool {
		return c.index[i] >= sum
	})
	if index >= len(c.index) {
		index = 0
	}
	return index
}

func (c *ConsistentHash) hash(s string) uint64 {
	sum := md5.Sum([]byte(s))
	return binary.LittleEndian.Uint64(sum[:])
}
