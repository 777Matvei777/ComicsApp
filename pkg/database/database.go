package database

import (
	"encoding/json"
	"fmt"
	"os"
)

type Item struct {
	URL      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

var items map[int]Item

func CreateDataBase(data map[int]interface{}, Db_path string) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	_ = os.WriteFile(Db_path, jsonData, 0644)
	_ = json.Unmarshal(jsonData, &items)
}

func WriteData(n int, Db_path string) {
	// f, _ := os.Open(Db_path)
	// defer f.Close()
	// jsondata, _ := os.ReadFile(Db_path)
	// var items map[int]Item
	// _ = json.Unmarshal(jsondata, &items)
	if n > 0 {
		for i := 1; i < n+1; i++ {
			fmt.Println(items[i])
		}
	} else {
		for k := range items {
			fmt.Println(items[k])
		}
	}

}
