package main

import (
	"context"
	"flag"
	"fmt"
	"myapp/pkg/config"
	"myapp/pkg/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var cs string
	flag.StringVar(&cs, "c", "", "config path")
	flag.Parse()
	cfg := config.New(cs)
	s := server.NewServer(cfg)
	s.AddServer()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)
	go func() {
		<-ch
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := s.Serv.Shutdown(ctx)
		if err != nil {
			fmt.Println("Server Shutdown Failed")
		}
	}()
	s.RunServer()
}
