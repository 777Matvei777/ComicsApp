package main

import (
	"context"
	"flag"
	"myapp/pkg/config"
	"myapp/pkg/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var cs string
	flag.StringVar(&cs, "c", "", "config path")
	flag.Parse()
	cfg := config.New(cs)
	ctx, cancelFunc := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)
	go func() {
		<-ch
		cancelFunc()
	}()
	s := server.NewServer(cfg, ctx)
	s.RunServer()
}
