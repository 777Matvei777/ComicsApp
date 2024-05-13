package models

type Item struct {
	Id       int      `json:"id"`
	URL      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

type Comic struct {
	New   int `json:"new"`
	Total int `json:"total"`
}

type ImageSearchResult struct {
	URLs []string `json:"urls"`
}

type KeywordIndex struct {
	Keyword string
	Index   []int
}

func NewComic() *Comic {
	return &Comic{
		New:   0,
		Total: 0,
	}
}
