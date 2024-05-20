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

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewComic() *Comic {
	return &Comic{
		New:   0,
		Total: 0,
	}
}
