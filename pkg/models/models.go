package models

import (
	"context"
	"net/http"
)

type Item struct {
	Id       int      `json:"id"`
	URL      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

type Database interface {
	GetUrlByComicId(id int) string
	GetComicDatabase() map[int]bool
	CheckDataBase(ctx context.Context) (int, map[int]bool)
	SizeDatabase() (int, error)
	GetUserByusername(user *User, creds *Credentials) error
	CreateComic(value []Item) error
	BuildIndex() ([]KeywordIndex, error)
	CreateIndex(keywordIndices []KeywordIndex) error
	GetComicsByQuery(searchQuery []string) []string
}

type Client interface {
	CreateDataBase(ctx context.Context)
	Start(ctx context.Context)
	SearhDatabase(searchFlag string, ctx context.Context) []string
	CheckDataBase(ctx context.Context) (int, map[int]bool)
	SizeDatabase() (int, error)
	LoginWithDb(w http.ResponseWriter, r *http.Request)
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
