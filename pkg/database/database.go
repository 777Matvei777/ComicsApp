package database

import (
	"encoding/json"
	"fmt"
	"myapp/pkg/words"
	"myapp/pkg/xkcd"
	"net/http"
	"os"
)

type DataBase struct {
	Id         int    `json:"num"`
	Alt        string `json:"alt"`
	Transcript string `json:"transcript"`
	Url        string `json:"img"`
}

type Item struct {
	URL      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

var Db []DataBase

func CreateDb(client *http.Client, Url string) {
	data := xkcd.Parse(client, Url)
	for i := 0; i < len(*data); i++ {
		var one_data DataBase
		_ = json.Unmarshal([]byte((*data)[i]), &one_data)
		Db = append(Db, one_data)
	}
}

func CreateJson(client *http.Client, Url string, Db_path string) *[]byte {
	CreateDb(client, Url)
	data := make(map[int]interface{})
	for i := 0; i < len(Db); i++ {
		keywords := fmt.Sprintf("%s %s", (Db)[i].Alt, (Db)[i].Transcript)
		arr_words := words.SplitString(keywords)
		normalized, _ := words.Stemming(arr_words)
		value := map[string]interface{}{
			"url":      (Db)[i].Url,
			"keywords": normalized,
		}
		data[(Db)[i].Id] = value
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	path := fmt.Sprintf("pkg/database/%s", Db_path)
	_ = os.WriteFile(path, jsonData, 0644)
	return &jsonData
}

func WriteData(n int, Db_path string) {
	path := fmt.Sprintf("pkg/database/%s", Db_path)
	f, _ := os.Open(path)
	defer f.Close()
	jsondata, _ := os.ReadFile(path)
	var items map[int]Item
	_ = json.Unmarshal(jsondata, &items)
	for i := 1; i < n+1; i++ {
		fmt.Println(items[i])
	}
}
