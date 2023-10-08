package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

func Shutdown(ctx context.Context, cancel context.CancelFunc, closers ...io.Closer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		for _, cls := range closers {
			if err := cls.Close(); err != nil {
				fmt.Println(err)
			}
		}
	}()

	select {
	case <-ctx.Done():
		return
	case <-c:
		cancel()
		return
	}
}
