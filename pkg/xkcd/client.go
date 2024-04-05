package xkcd

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func Parse(client *http.Client, Url string) *[][]byte {
	var Db [][]byte
	for i := 1; i < 2915; i++ { //2914
		adress := fmt.Sprintf("%s/%d/info.0.json", Url, i)
		resp, err := client.Get(adress)
		if err != nil {
			log.Fatalf("error getting %s", err)
		}
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		Db = append(Db, data)
	}
	return &Db
}
