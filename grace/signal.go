package grace

import (
	"os"
	"os/signal"
	"syscall"
)

func Signal() <-chan os.Signal {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	return quit
}
