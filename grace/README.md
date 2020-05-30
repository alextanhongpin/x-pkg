# grace

Implements graceful shutdown for golang router that implements the `http.Handler`.

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/alextanhongpin/pkg/grace"
)

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello go")
	})

	// Handles graceful shutdown.
	shutdown := grace.New(r, 8080)

	// Listens to CTRL + C.
	<-grace.Signal()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shutdown(ctx)
}
```

## Shutdown Group

```go
func main() {
	sg := NewShutdownGroup()

	sg.Add(func(ctx context.Context)  {
		// Add a cancel mechanism here to force termination of a resource.
		fmt.Println("shutting down 1")
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		select {
		case <-time.After(5 * time.Second):
			fmt.Println("shut down 1")
		case <-ctx.Done():
			fmt.Printf("force shut down 1: %v\n", ctx.Err())
			return
		}
	})

	sg.Add(func(ctx context.Context) {
		fmt.Println("shutting down 2")
		time.Sleep(2 * time.Second)
		fmt.Println("shut down 2")
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	t0 := time.Now()

	fmt.Println("initializing graceful shutdown")
	sg.Close(ctx)

	fmt.Println("time taken:", time.Since(t0))
	fmt.Println("process terminated")
}
```
