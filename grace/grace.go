package grace

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Shutdown func(context.Context)

type Shutdowns []Shutdown

func (shutdowns Shutdowns) Close(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(len(shutdowns))
	for _, shutdown := range shutdowns {
		go func(shutdown Shutdown) {
			defer wg.Done()
			shutdown(ctx)
		}(shutdown)
	}
	wg.Wait()
}

func Signal() <-chan os.Signal {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	return quit
}

// New returns a new shutdown function given a http.Handler function.
func New(handler http.Handler, port string) Shutdown {
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        handler,
		WriteTimeout:   10 * time.Second,
		ReadTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	idleConnsClosed := make(chan struct{})
	go func() {
		log.Printf("listening to port *:%s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
		close(idleConnsClosed)
	}()

	return func(ctx context.Context) {
		log.Println("shutting down")
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("server shutdown:", err)
		}
		select {
		case <-idleConnsClosed:
			log.Println("shutdown gracefully")
			return
		case <-ctx.Done():
			log.Println("shutdown abruptly after 5 seconds timeout")
			return
		}
	}
}
