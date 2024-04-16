package database

import (
	"encoding/json"
	"fmt"
	"log"
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
		log.Println("Marshaling error: ", err)
	}
	if _, err := os.Stat(Db_path); err == nil {
		file, err := os.ReadFile(Db_path)
		if err != nil {
			log.Println("Reading error: ", err)
		}
		var curr_items map[int]Item
		err = json.Unmarshal(file, &curr_items)
		if err != nil {
			log.Println("Unmarshaling error: ", err)
		}
		err = json.Unmarshal(jsonData, &curr_items)
		if err != nil {
			log.Println("unmarshaling error: ", err)
		}
		curr_data, err := json.Marshal(curr_items)
		if err != nil {
			log.Println("Marshaling error: ", err)
		}
		os.WriteFile(Db_path, curr_data, 0644)
		items = curr_items
	} else {
		err = os.WriteFile(Db_path, jsonData, 0644)
		if err != nil {
			log.Println("Writing error: ", err)
		}
		err = json.Unmarshal(jsonData, &items)
		if err != nil {
			log.Println("Unmarshaling error: ", err)
		}
	}
	fmt.Printf("Data saved to %s\n", Db_path)
	fmt.Printf("%d comics in file", len(items))
	os.Exit(0)
}

func CheckDataBase(Db_path string) (int, map[int]bool) {
	file, err := os.Open(Db_path)
	if err != nil {
		fmt.Println("error opened file")
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&items)
	if err != nil {
		fmt.Println("error json decoding")
	}
	comics_id := 1
	flag := false
	exist := make(map[int]bool)
	for ; comics_id < len(items); comics_id++ {
		if comics_id != 404 {
			if _, ok := items[comics_id]; !ok {
				if !flag {
					break
				}
			} else {
				exist[comics_id] = true
			}
		}
	}
	if comics_id == len(items) {
		return 0, exist
	}
	fmt.Printf("Missed comics with id %d\n", comics_id)
	return comics_id, exist
}
