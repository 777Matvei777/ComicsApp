package app

import (
	"context"
	"fmt"
	"myapp/pkg/database"
	"myapp/pkg/words"
	"myapp/pkg/xkcd"
)

func Normalize(keywords string) []string {
	arr_words := words.SplitString(keywords)
	normalized, _ := words.Stemming(arr_words)
	return normalized
}

func CreateJson(Url string, Db_path string, Parallel int, ctx context.Context, num int, exist map[int]int) {
	Db := xkcd.Parse(Url, Parallel, ctx, num, exist)
	data := make(map[int]interface{})
	for i := 0; i < len(Db); i++ {
		keywords := fmt.Sprintf("%s %s", (Db)[i].Alt, (Db)[i].Transcript)
		normalized := Normalize(keywords)
		value := map[string]interface{}{
			"url":      (Db)[i].Url,
			"keywords": normalized,
		}
		data[(Db)[i].Id] = value
	}
	database.CreateDataBase(data, Db_path)
	select {}
}

func Start(Url string, Db_path string, parallel int, ctx context.Context, num int, exist map[int]int) {
	CreateJson(Url, Db_path, parallel, ctx, num, exist)
}

func CheckDataBase(Db_path string) (int, map[int]int) {
	return database.CheckDataBase(Db_path)

}
