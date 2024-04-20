package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
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
	fmt.Printf("%d comics in file\n", len(items))
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

func createIndexFile() map[string][]int {
	var items map[int]Item
	if _, err := os.Stat("pkg/database/index.json"); err != nil {
		file, _ := os.ReadFile("pkg/database/database.json")
		_ = json.Unmarshal(file, &items)
		index := make(map[string][]int)
		for k, words := range items {
			for _, word := range words.Keywords {
				index[word] = append(index[word], k)
			}
		}
		jsonData, err := json.Marshal(index)
		if err != nil {
			fmt.Println("Marshaling error")
		}
		err = os.WriteFile("pkg/database/index.json", jsonData, 0644)
		if err != nil {
			fmt.Println("writing error")
		}
		return index
	} else {
		file, err := os.ReadFile("pkg/database/index.json")
		if err != nil {
			fmt.Println("Reading error")
		}
		index := make(map[string][]int)
		err = json.Unmarshal(file, &index)
		if err != nil {
			fmt.Println("Unmarshaling error")
		}
		return index
	}

}

func SearchDatabase(query []string) []string {
	var comics []string
	stat := make(map[int]int)
	for index, comic := range items {
		for _, query_word := range query {
			for _, keyword := range comic.Keywords {
				if query_word == keyword {
					stat[index]++
				}
			}
		}
	}
	keys := make([]int, 0, len(stat))
	for k := range stat {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return stat[keys[i]] > stat[keys[j]]
	})
	for _, k := range keys {
		comic_url := items[k].URL
		comics = append(comics, comic_url)
		if len(comics) >= 10 {
			break
		}
	}
	return comics
}
func SearchByIndex(query []string) []string {
	index := createIndexFile()
	stat := make(map[int]int)
	var comics []string
	for _, v := range query {
		if ids, found := index[v]; found {
			for _, i := range ids {
				stat[i]++
			}
		}
	}
	keys := make([]int, 0, len(stat))
	for k := range stat {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return stat[keys[i]] > stat[keys[j]]
	})
	for _, k := range keys {
		comic_url := items[k].URL
		if found := slices.Contains(comics, comic_url); !found {
			comics = append(comics, comic_url)
		}
		if len(comics) >= 10 {
			break
		}
	}
	return comics
}
