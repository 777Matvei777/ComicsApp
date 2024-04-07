package main

import (
	"flag"
	"fmt"
	"myapp/pkg/config"
	"myapp/pkg/database"
)

func main() {
	var o bool
	var n int
	flag.BoolVar(&o, "o", false, "json structure")
	flag.IntVar(&n, "n", 0, "first n comics")
	flag.Parse()

	cfg := config.New()
	jsonData := database.CreateJson(cfg.Url, cfg.Db_file)
	if o {
		if n > 0 {
			database.WriteData(n, cfg.Db_file)
		} else {
			fmt.Println(string(*jsonData))
		}
	}
}
