package app

import (
	"context"
	"fmt"
	"myapp/pkg/database"
	"myapp/pkg/words"
	"myapp/pkg/xkcd"
	"os"
)

func Normalize(keywords string) []string {
	arr_words := words.SplitString(keywords)
	normalized, _ := words.Stemming(arr_words)
	return normalized
}

func CreateJson(Url string, Db_path string, Parallel int, ctx context.Context, num int, exist map[int]bool) {
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
}

func Start(Url string, Db_path string, parallel int, ctx context.Context, num int) {
	exist_flag := false
	if _, err := os.Stat(Db_path); err == nil {
		fmt.Println("File already exist")
		exist_flag = true
	}
	if exist_flag {
		num, exist := CheckDataBase(Db_path)
		if num != 0 {
			CreateJson(Url, Db_path, parallel, ctx, num, exist)
		} else {
			fmt.Println("All comics in file")
		}
	} else {
		mp := make(map[int]bool)
		CreateJson(Url, Db_path, parallel, ctx, num, mp)
	}
}
func SearhDatabase(searchFlag *string, indexFlag *bool) {
	if *searchFlag != "" {
		fmt.Println("Найденные комиксы: ")
		split_query := words.SplitString(*searchFlag)
		normalized_query, _ := words.Stemming(split_query)
		if *indexFlag {
			comics_url := database.SearchByIndex(normalized_query)
			for k, url := range comics_url {
				fmt.Println(k+1, " ", url)
			}
		} else {
			comics_url := database.SearchDatabase(normalized_query)
			for k, url := range comics_url {
				fmt.Println(k+1, " ", url)
			}
		}

	}
}

func CheckDataBase(Db_path string) (int, map[int]bool) {
	return database.CheckDataBase(Db_path)

}
