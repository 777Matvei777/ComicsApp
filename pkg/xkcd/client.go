package xkcd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type XkcdStruct struct {
	Id         int    `json:"num"`
	Alt        string `json:"alt"`
	Transcript string `json:"transcript"`
	Url        string `json:"img"`
}

func Parse(Url string) []XkcdStruct {
	var Db []XkcdStruct
	for i := 1; ; i++ { //2914
		adress := fmt.Sprintf("%s/%d/info.0.json", Url, i)
		resp, err := http.Get(adress)
		if err != nil {
			fmt.Println("getting error")
		}
		defer resp.Body.Close()
		if resp.StatusCode == 404 && i != 404 {
			resp.Body.Close()
			fmt.Printf("Загрузилось %d комиксов\n", i)
			break
		}
		var one_data XkcdStruct
		data, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal([]byte(data), &one_data)
		if err != nil {
			fmt.Println("error:", err)
		}
		Db = append(Db, one_data)
		if i%100 == 0 {
			fmt.Printf("Загрузилось %d комиксов\n", i)
		}
	}
	return Db
}
