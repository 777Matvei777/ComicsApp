package main

import (
	"context"
	"flag"
	"fmt"
	"myapp/pkg/app"
	"myapp/pkg/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var cs string
	flag.StringVar(&cs, "c", "", "config path")
	flag.Parse()
	cfg := config.New(cs)
	exist_flag := false
	if _, err := os.Stat(cfg.Db_file); err == nil {
		fmt.Println("File already exist")
		exist_flag = true
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)
	go func() {
		for {
			<-ch
			cancelFunc()
		}
	}()
	if exist_flag {
		num, exist := app.CheckDataBase(cfg.Db_file)
		if num != 0 {
			app.Start(cfg.Url, cfg.Db_file, cfg.Parallel, ctx, num, exist)
		} else {
			fmt.Println("All comics in file")
		}
	} else {
		mp := make(map[int]int)
		app.Start(cfg.Url, cfg.Db_file, cfg.Parallel, ctx, 1, mp)
	}
}
