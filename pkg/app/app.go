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
	Ctx   context.Context
	Num   int
	Exist map[int]bool
	Db    *database.DataBase
}

func NewClient(cfg *config.Config, ctx context.Context, num int) *Client {
	c := &Client{
		Cfg:   cfg,
		Ctx:   ctx,
		Num:   num,
		Exist: make(map[int]bool),
		Db: &database.DataBase{
			Db_path: cfg.DbFile,
			Items:   make(map[int]models.Item),
			Index:   make(map[string][]int),
		},
	}
	go checkForNewComics(c)
	return c
}

func (c *Client) CreateJson() {
	Db := xkcd.Parse(c.Cfg.Url, c.Cfg.Parallel, c.Ctx, c.Num, c.Exist)
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

func (c *Client) Start() {
	exist_flag := false
	if _, err := os.Stat(c.Cfg.DbFile); err == nil {
		fmt.Println("File already exist")
		exist_flag = true
	}
	if exist_flag {
		c.Num, c.Exist = c.CheckDataBase()
		if c.Num != 0 {
			c.CreateJson()
		} else {
			fmt.Println("All comics in file")
		}
	} else {
		c.CreateJson()
	}
	c.Db.CreateIndexFile()
}

func (c *Client) SearhDatabase(searchFlag string) []string {
	normalized_query := words.Normalize(searchFlag)
	comics_url := c.Db.SearchByIndex(normalized_query)
	return comics_url
}

func (c *Client) CheckDataBase() (int, map[int]bool) {
	return c.Db.CheckDataBase()

}

func (c *Client) SizeDatabase() int {
	return c.Db.SizeDatabase()
}

func checkForNewComics(c *Client) {
	for {
		c.Start()
		time.Sleep(24 * time.Hour)
	}
}
