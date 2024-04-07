package database

import (
	"encoding/json"
	"fmt"
	"myapp/pkg/app"
	"os"
)

type Item struct {
	URL      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

func CreateJson(Url string, Db_path string) *[]byte {
	app.Parse(Url)
	data := app.NewJson()
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	_ = os.WriteFile(Db_path, jsonData, 0644)
	return &jsonData
}

func WriteData(n int, Db_path string) {
	f, _ := os.Open(Db_path)
	defer f.Close()
	jsondata, _ := os.ReadFile(Db_path)
	var items map[int]Item
	_ = json.Unmarshal(jsondata, &items)
	for i := 1; i < n+1; i++ {
		fmt.Println(items[i])
	}
}
