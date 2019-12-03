// Package grace implements graceful shutdown for your server and provides an
// interface for graceful termination of other infrastructure code.
package grace

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// New returns a new shutdown function given a http.Handler function.
func New(handler http.Handler, port int) Shutdown {
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        handler,
		WriteTimeout:   10 * time.Second,
		ReadTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	idleConnsClosed := make(chan struct{})
	go func() {
		log.Printf("[grace] listening to port *:%d\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[grace] listen: %s", err)
		}
		close(idleConnsClosed)
	}()

	return func(ctx context.Context) {
		log.Println("[grace] shutting down")
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("[grace] server shutdown:", err)
		}
		select {
		case <-idleConnsClosed:
			log.Println("[grace] shutdown gracefully")
			return
		case <-ctx.Done():
			log.Println("[grace] shutdown abruptly after 5 seconds timeout")
			return
		}
	}
}
