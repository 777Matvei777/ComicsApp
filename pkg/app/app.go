package app

import (
	"context"
	"fmt"
	"myapp/pkg/config"
	"myapp/pkg/database"
	"myapp/pkg/models"
	"myapp/pkg/words"
	"myapp/pkg/xkcd"
	"os"
	"time"
)

type Client struct {
	Cfg   *config.Config
	Num   int
	Exist map[int]bool
	Db    *database.DataBase
}

func NewClient(cfg *config.Config, ctx context.Context, num int) *Client {
	c := &Client{
		Cfg:   cfg,
		Num:   num,
		Exist: make(map[int]bool),
		Db: &database.DataBase{
			Db_path: cfg.DbFile,
			Items:   make(map[int]models.Item),
			Index:   make(map[string][]int),
		},
	}
	checkForNewComics(22, 15, 0, c)
	return c
}

func (c *Client) CreateJson(ctx context.Context) {
	Db := xkcd.Parse(c.Cfg.Url, c.Cfg.Parallel, ctx, c.Num, c.Exist)
	data := make(map[int]interface{})
	for i := 0; i < len(Db); i++ {
		keywords := fmt.Sprintf("%s %s", (Db)[i].Alt, (Db)[i].Transcript)
		normalized := words.Normalize(keywords)
		value := map[string]interface{}{
			"url":      (Db)[i].Url,
			"keywords": normalized,
		}
		data[(Db)[i].Id] = value
	}
	c.Db.CreateDataBase(data)
}

func (c *Client) Start(ctx context.Context) {
	exist_flag := false
	if _, err := os.Stat(c.Cfg.DbFile); err == nil {
		fmt.Println("File already exist")
		exist_flag = true
	}
	if exist_flag {
		c.Num, c.Exist = c.CheckDataBase(ctx)
		if c.Num != 0 {
			c.CreateJson(ctx)
		} else {
			fmt.Println("All comics in file")
		}
	} else {
		c.CreateJson(ctx)
	}
	c.Db.CreateIndexFile()
}

func (c *Client) SearhDatabase(searchFlag string, ctx context.Context) []string {
	comics_url := make([]string, 0)
	select {
	case <-ctx.Done():
		return comics_url
	default:
		normalized_query := words.Normalize(searchFlag)
		comics_url = c.Db.SearchByIndex(normalized_query)

	}
	return comics_url
}

func (c *Client) CheckDataBase(ctx context.Context) (int, map[int]bool) {
	return c.Db.CheckDataBase(ctx)

}

func (c *Client) SizeDatabase() int {
	return c.Db.SizeDatabase()
}

func checkForNewComics(hour, min, sec int, c *Client) {
	now := time.Now()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, now.Location())
	if now.After(nextRun) {
		nextRun = nextRun.Add(24 * time.Hour)
	}
	durationUntilNextRun := nextRun.Sub(now)

	timer := time.NewTimer(durationUntilNextRun)

	ctx := context.Background()

	go func() {
		for {
			<-timer.C
			c.Start(ctx)
			timer.Reset(24 * time.Hour)
		}
	}()
}
