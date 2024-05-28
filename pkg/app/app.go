package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"myapp/pkg/config"
	"myapp/pkg/database"
	"myapp/pkg/models"
	"myapp/pkg/words"
	"myapp/pkg/xkcd"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Client struct {
	Cfg   *config.Config
	Num   int
	Exist map[int]bool
	Db    models.Database
}

func NewClient(cfg *config.Config, num int) *Client {
	c := &Client{
		Cfg:   cfg,
		Num:   num,
		Exist: make(map[int]bool),
	}
	Db, err := database.NewPostgreSQL(cfg.Postgresql)
	if err != nil {
		log.Fatal(err)
	}
	c.Db = Db
	checkForNewComics(18, 31, 0, c)
	return c
}

func (c *Client) CreateDataBase(ctx context.Context) {
	Db := xkcd.Parse(c.Cfg.Url, c.Cfg.Parallel, ctx, c.Num, c.Exist)
	values := make([]models.Item, 0)
	for i := 0; i < len(Db); i++ {
		keywords := fmt.Sprintf("%s %s", (Db)[i].Alt, (Db)[i].Transcript)
		normalized := words.Normalize(keywords)
		value := models.Item{
			Id:       Db[i].Id,
			URL:      (Db)[i].Url,
			Keywords: normalized,
		}
		values = append(values, value)

	}
	err := c.Db.CreateComic(values)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Client) Start(ctx context.Context) {
	c.Num, c.Exist = c.CheckDataBase(ctx)
	if c.Num != 0 {
		c.CreateDataBase(ctx)
	} else {
		fmt.Println("All comics in database")
	}
	index, err := c.Db.BuildIndex()
	if err != nil {
		log.Fatal(err)
	}
	c.Db.CreateIndex(index)
}

func (c *Client) SearhDatabase(searchFlag string, ctx context.Context) []string {
	comics_url := make([]string, 0)
	select {
	case <-ctx.Done():
		return comics_url
	default:
		normalized_query := words.Normalize(searchFlag)
		comics_url = c.Db.GetComicsByQuery(normalized_query)

	}
	return comics_url
}

func (c *Client) CheckDataBase(ctx context.Context) (int, map[int]bool) {
	return c.Db.CheckDataBase(ctx)

}

func (c *Client) SizeDatabase() (int, error) {
	count, err := c.Db.SizeDatabase()
	if err != nil {
		return 0, nil
	}
	return count, nil
}

func checkForNewComics(hour, min, sec int, c *Client) {
	now := time.Now()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, now.Location())
	if now.After(nextRun) {
		nextRun = nextRun.Add(24 * time.Hour)
	}
	durationUntilNextRun := nextRun.Sub(now)

	timer := time.NewTimer(durationUntilNextRun)

	ctx := context.Background()

	go func() {
		for {
			<-timer.C
			c.Start(ctx)
			timer.Reset(24 * time.Hour)
		}
	}()
}

func (c *Client) LoginWithDb(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	err = c.Db.GetUserByusername(&user, &creds)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if creds.Password != user.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": creds.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Minute * time.Duration(c.Cfg.Token_max_time)),
	})

	tokenString, err := token.SignedString([]byte("secret_key"))
	if err != nil {
		http.Error(w, "Error signing token", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(tokenString))
}
