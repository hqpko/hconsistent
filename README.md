# hconsistent

### example:

```go
package main

import (
	"fmt"

	"github.com/hqpko/hconsistent"
)

func main() {
	c := hconsistent.NewConsistentHash()
	c.Add("host-1")
	c.Add("host-2")
	c.Add("host-3")
	host := c.Get("key-01")
	fmt.Println(host)
}
```