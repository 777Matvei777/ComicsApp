package benching

import (
	"encoding/json"
	"myapp/pkg/models"
	"myapp/pkg/words"
	"os"
	"slices"
	"sort"
	"testing"
)

var Items map[int]models.Item

func SearchDatabase(query []string) {
	file, _ := os.Open("myapp/pkg/database/database.json")
	json.NewDecoder(file).Decode(&Items)
	var comics []string
	stat := make(map[int]int)
	for index, comic := range Items {
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
		comic_url := Items[k].URL
		comics = append(comics, comic_url)
		if len(comics) >= 10 {
			break
		}
	}
}

func SearchByIndex(query []string) {
	file, _ := os.Open("myapp/pkg/database/index.json")
	Index := make(map[string][]int)
	json.NewDecoder(file).Decode(&Index)
	stat := make(map[int]int)
	var comics []string
	for _, v := range query {
		if ids, found := Index[v]; found {
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
		comic_url := Items[k].URL
		if found := slices.Contains(comics, comic_url); !found {
			comics = append(comics, comic_url)
		}
		if len(comics) >= 10 {
			break
		}
	}
}
func normal_query(query string) []string {
	arr_query := words.SplitString(query)
	normal_query, _ := words.Stemming(arr_query)
	return normal_query
}

func BenchmarkSearch(b *testing.B) {
	arr := normal_query("I'm following your questions")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		SearchDatabase(arr)
	}

}

func BenchmarkIndexedSearch(b *testing.B) {
	arr := normal_query("I'm following your questions")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SearchByIndex(arr)
	}
}
