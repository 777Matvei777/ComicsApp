package app

import (
	"fmt"
	"myapp/pkg/words"
	"myapp/pkg/xkcd"
)

func Parse(Url string) {
	xkcd.Parse(Url)
}

func Normalize(keywords string) []string {
	arr_words := words.SplitString(keywords)
	normalized, _ := words.Stemming(arr_words)
	return normalized
}

func NewJson() map[int]interface{} {
	data := make(map[int]interface{})
	for i := 0; i < len(xkcd.Db); i++ {
		keywords := fmt.Sprintf("%s %s", (xkcd.Db)[i].Alt, (xkcd.Db)[i].Transcript)
		normalized := Normalize(keywords)
		value := map[string]interface{}{
			"url":      (xkcd.Db)[i].Url,
			"keywords": normalized,
		}
		data[(xkcd.Db)[i].Id] = value
	}
	return data
}
