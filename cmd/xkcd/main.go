package main

import (
	"flag"
	"myapp/pkg/app"
	"myapp/pkg/config"
)

func main() {
	var o bool
	var n int
	flag.BoolVar(&o, "o", false, "json structure")
	flag.IntVar(&n, "n", 0, "first n comics")
	flag.Parse()

	cfg := config.New()
	app.Start(cfg.Url, cfg.Db_file)
	if o {
		if n > 0 {
			app.WriteData(n, cfg.Db_file)
		} else {
			app.WriteData(0, cfg.Db_file)
		}
	}
}
