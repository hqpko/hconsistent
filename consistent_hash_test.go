package hconsistent

import (
	"fmt"
	"testing"
)

func TestConsistentHash(t *testing.T) {
	c := NewConsistentHash()
	c.Add("host-1")
	for i := 0; i < 30; i++ {
		if c.Get(fmt.Sprintf("host-%d", i)) != "host-1" {
			t.Error("consistent hash get fail.")
		}
	}
}
