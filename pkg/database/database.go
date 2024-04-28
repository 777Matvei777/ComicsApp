package database

import (
	"encoding/json"
	"fmt"
	"log"
	"myapp/pkg/models"
	"os"
	"slices"
	"sort"
)

type DataBase struct {
	Db_path string
	Items   map[int]models.Item
	Index   map[string][]int
}

func (d *DataBase) CreateDataBase(data map[int]interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Marshaling error: ", err)
	}
	if _, err := os.Stat(d.Db_path); err == nil {
		file, err := os.ReadFile(d.Db_path)
		if err != nil {
			log.Println("Reading error: ", err)
		}
		var curr_items map[int]models.Item
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
		os.WriteFile(d.Db_path, curr_data, 0644)
		d.Items = curr_items
	} else {
		err = os.WriteFile(d.Db_path, jsonData, 0644)
		if err != nil {
			log.Println("Writing error: ", err)
		}
		err = json.Unmarshal(jsonData, &d.Items)
		if err != nil {
			log.Println("Unmarshaling error: ", err)
		}
	}
	fmt.Printf("Data saved to %s\n", d.Db_path)
	fmt.Printf("%d comics in file\n", len(d.Items))
}

func (d *DataBase) CheckDataBase() (int, map[int]bool) {
	file, err := os.Open(d.Db_path)
	if err != nil {
		fmt.Println("error opened file")
	}
	defer file.Close()
	json.NewDecoder(file).Decode(&d.Items)
	res_id := 0
	flag := false
	exist := make(map[int]bool)
	for comics_id := 1; comics_id < len(d.Items); comics_id++ {
		if comics_id != 404 {
			if _, ok := d.Items[comics_id]; !ok {
				if !flag {
					res_id = comics_id
					flag = true
				}
			} else {
				exist[comics_id] = true
			}
		}
	}
	if res_id == 0 {
		return 0, exist
	}
	fmt.Printf("Missed comics with id %d\n", res_id)
	return res_id, exist
}

func (d *DataBase) CreateIndexFile() {
	if _, err := os.Stat("pkg/database/index.json"); err != nil {
		data := make(map[string][]int)
		for k, words := range d.Items {
			for _, word := range words.Keywords {
				data[word] = append(data[word], k)
			}
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Marshaling error")
		}
		err = os.WriteFile("pkg/database/index.json", jsonData, 0644)
		if err != nil {
			fmt.Println("writing error")
		}
		d.Index = data
	} else {
		file, err := os.ReadFile("pkg/database/index.json")
		if err != nil {
			fmt.Println("Reading error")
		}
		data := make(map[string][]int)
		err = json.Unmarshal(file, &data)
		if err != nil {
			fmt.Println("Unmarshaling error")
		}
		d.Index = data
	}

}
func (d *DataBase) SizeDatabase() int {
	return len(d.Items)

}
func (d *DataBase) SearchDatabase(query []string) []string {
	var comics []string
	stat := make(map[int]int)
	for index, comic := range d.Items {
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
		comic_url := d.Items[k].URL
		comics = append(comics, comic_url)
		if len(comics) >= 10 {
			break
		}
	}
	return comics
}
func (d *DataBase) SearchByIndex(query []string) []string {
	stat := make(map[int]int)
	var comics []string
	for _, v := range query {
		if ids, found := d.Index[v]; found {
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
		comic_url := d.Items[k].URL
		if found := slices.Contains(comics, comic_url); !found {
			comics = append(comics, comic_url)
		}
		if len(comics) >= 10 {
			break
		}
	}
	return comics
}
