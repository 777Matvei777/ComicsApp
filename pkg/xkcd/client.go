package xkcd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type XkcdStruct struct {
	Id         int    `json:"num"`
	Alt        string `json:"alt"`
	Transcript string `json:"transcript"`
	Url        string `json:"img"`
}

func Parse(Url string, Parallel int, ctx context.Context, num int, exist map[int]bool) []XkcdStruct {
	var Db []XkcdStruct
	var wg sync.WaitGroup
	var mutex sync.Mutex
	ch := make(chan int, Parallel)
	//found404 := make(chan bool)
	flag := false
	for i := num; !flag; i++ {
		select {
		case <-ctx.Done():
			return Db
		default:
			if _, ok := exist[i]; !ok {
				ch <- i
				wg.Add(1)
				i := i
				go func(i int) {
					defer wg.Done()
					defer func() { <-ch }()

					address := fmt.Sprintf("%s/%d/info.0.json", Url, i)
					resp, err := http.Get(address)
					if err != nil {
						fmt.Println("getting error:", err)
						return
					}
					defer resp.Body.Close()
					if resp.StatusCode == 404 && i != 404 {
						flag = true
						return
					}
					var oneData XkcdStruct
					data, err := io.ReadAll(resp.Body)
					if err != nil {
						fmt.Println("reading error:", err)
						return
					}
					_ = json.Unmarshal([]byte(data), &oneData)
					mutex.Lock()
					Db = append(Db, oneData)
					mutex.Unlock()
					if i%100 == 0 {
						fmt.Printf("Загружен %d-ый комикс\n", i)
					}
					// select {
					// case <-found404:
					// 	return
					// default:
					// 	defer func() { <-ch }()

					// 	address := fmt.Sprintf("%s/%d/info.0.json", Url, i)
					// 	resp, err := http.Get(address)
					// 	if err != nil {
					// 		fmt.Println("getting error:", err)
					// 		return
					// 	}
					// 	defer resp.Body.Close()
					// 	if resp.StatusCode == 404 && i != 404 {
					// 		flag = true
					// 		found404 <- true
					// 		return
					// 	}
					// 	var oneData XkcdStruct
					// 	data, err := io.ReadAll(resp.Body)
					// 	if err != nil {
					// 		fmt.Println("reading error:", err)
					// 		return
					// 	}
					// 	_ = json.Unmarshal([]byte(data), &oneData)
					// 	mutex.Lock()
					// 	Db = append(Db, oneData)
					// 	mutex.Unlock()
					// 	if i%100 == 0 {
					// 		fmt.Printf("Загружен %d-ый комикс\n", i)
					// 	}
					// }

				}(i)
			}

		}
	}
	wg.Wait()
	close(ch)
	fmt.Printf("Загружено %d комиксов\n", len(Db))
	return Db
}
