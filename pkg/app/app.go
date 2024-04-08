package app

import (
	"fmt"
	"myapp/pkg/database"
	"myapp/pkg/words"
	"myapp/pkg/xkcd"
)

func Parse(Url string) {
	xkcd.Parse(Url)
}

func Normalize(keywords string) []string {
	arr_words := words.SplitString(keywords)
	normalized, _ := words.Stemming(arr_words)
	return normalized
}

func CreateJson(Url string, Db_path string) {
	Db := xkcd.Parse(Url)
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
}

func WriteData(n int, Db_path string) {
	database.WriteData(n, Db_path)
}

func Start(Url string, Db_path string) {
	CreateJson(Url, Db_path)
}
