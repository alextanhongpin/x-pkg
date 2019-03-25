```go
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/alextanhongpin/pkg/ttlmap"
)

func main() {
	ttlMap := ttlmap.New(1 * time.Second)
	shutdown := ttlMap.Cleanup(time.Second)

	add := func(n int) {
		for i := 0; i < n; i++ {
			val := rand.Intn(100000)
			s := strconv.FormatInt(int64(val), 10)
			ttlMap.Put(s, s)
		}
		log.Println(ttlMap.Len())
	}

	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Duration(i*950) * time.Millisecond)
			n := 200 + rand.Intn(800)
			add(n)
		}(i)
	}

	done := make(chan interface{})
	go func() {
		time.Sleep(10 * time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdown(ctx)
		close(done)
	}()
	<-done

	log.Println(ttlMap.Len())
	fmt.Println("done")
}
```
