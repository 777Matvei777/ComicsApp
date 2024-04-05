package main

import (
	"flag"
	"fmt"
	"myapp/pkg/config"
	"myapp/pkg/database"
	"net/http"
)

func main() {
	var o bool
	var n int
	flag.BoolVar(&o, "o", false, "json structure")
	flag.IntVar(&n, "n", 0, "first n comics")
	flag.Parse()

	cfg := config.New()
	client := &http.Client{}
	jsonData := database.CreateJson(client, cfg.Url, cfg.Db_file)
	if o {
		fmt.Println(string(*jsonData))
	}
	if n > 0 {
		database.WriteData(n, cfg.Db_file)
	}
}
