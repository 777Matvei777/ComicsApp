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

var Db []XkcdStruct

func Parse(Url string) {
	for i := 1; ; i++ { //2914
		adress := fmt.Sprintf("%s/%d/info.0.json", Url, i)
		resp, _ := http.Get(adress)
		if resp.StatusCode == 404 && i != 404 {
			resp.Body.Close()
			fmt.Printf("Загрузилось %d комиксов\n", i)
			break
		}
		defer resp.Body.Close()
		var one_data XkcdStruct
		data, _ := io.ReadAll(resp.Body)
		_ = json.Unmarshal([]byte(data), &one_data)
		Db = append(Db, one_data)
		if i%100 == 0 {
			fmt.Printf("Загрузилось %d комиксов\n", i)
		}
	}
}
