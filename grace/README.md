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
