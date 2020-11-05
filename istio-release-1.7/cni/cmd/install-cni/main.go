package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"istio.io/istio/cni/pkg/install-cni/cmd"
	"istio.io/pkg/log"
)

func main() {
	// 创建在终止信号上取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func(sigChan chan os.Signal, cancel context.CancelFunc) {
		sig := <-sigChan
		log.Infof("Exit signal received: %s", sig)
		cancel()
	}(sigChan, cancel)

	rootCmd := cmd.GetCommand()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
