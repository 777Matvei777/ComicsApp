package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"myapp/pkg/models"
	"sort"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

type PostgreSQL struct {
	DB *sql.DB
}

func NewPostgreSQL(connString string) (*PostgreSQL, error) {
	latestVersion := 1
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatal(err)
	}
	if dirty {
		log.Fatal("Database is in a dirty state")
	}
	if version < uint(latestVersion) {
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	}

	return &PostgreSQL{DB: db}, nil
}

func (p *PostgreSQL) CreateComic(value []models.Item) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			err = tx.Rollback()
			if err != nil {
				log.Println("error rollback", err)
			}
			panic(p)
		} else if err != nil {
			err = tx.Rollback()
			if err != nil {
				log.Println("error rollback", err)
			}
		}
	}()
	stmtComic, err := tx.Prepare("INSERT INTO comics (url) VALUES ($1) RETURNING id")
	if err != nil {
		return err
	}
	defer stmtComic.Close()

	stmtKeyword, err := tx.Prepare("INSERT INTO keywords (keyword, comic_id) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer stmtKeyword.Close()

	for id, comic := range value {
		var lastInsertId int
		err = stmtComic.QueryRow(comic.URL).Scan(&lastInsertId)
		if err != nil {
			return err
		}

		for _, keyword := range comic.Keywords {
			_, err = stmtKeyword.Exec(keyword, lastInsertId)
			if err != nil {
				return err
			}
		}
		if id%100 == 0 {
			fmt.Printf("Download %d comics\n", id)
		}
	}
	return tx.Commit()
}

func (p *PostgreSQL) BuildIndex() ([]models.KeywordIndex, error) {
	keywordMap := make(map[string][]int)

	rows, err := p.DB.Query("SELECT k.keyword, k.comic_id FROM keywords k JOIN comics c ON k.comic_id = c.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var keyword string
		var comicID int
		if err := rows.Scan(&keyword, &comicID); err != nil {
			return nil, err
		}
		keywordMap[keyword] = append(keywordMap[keyword], comicID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var keywordIndices []models.KeywordIndex
	for k, ids := range keywordMap {
		keywordIndices = append(keywordIndices, models.KeywordIndex{
			Keyword: k,
			Index:   ids,
		})
	}
	return keywordIndices, nil
}
func (p *PostgreSQL) CreateIndex(keywordIndices []models.KeywordIndex) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			err = tx.Rollback()
			if err != nil {
				log.Fatal("error rollback", err)
			}
			panic(p)
		} else if err != nil {
			err = tx.Rollback()
			if err != nil {
				log.Fatal("error rollback", err)
			}
		}
	}()
	stmt, err := tx.Prepare("INSERT INTO keyword_index (keyword, comic_ids) VALUES ($1, $2) ON CONFLICT (keyword) DO UPDATE SET comic_ids = excluded.comic_ids")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, ki := range keywordIndices {
		indices := "{" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ki.Index)), ","), "[]") + "}"
		_, err = stmt.Exec(ki.Keyword, indices)
		if err != nil {
			return err
		}

	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil

}

func (p *PostgreSQL) GetComicsByQuery(searchQuery []string) []string {
	comics := make([]string, 0)
	stat := make(map[int]int)
	for _, word := range searchQuery {
		var idstring string
		err := p.DB.QueryRow("SELECT comic_ids FROM keyword_index WHERE keyword=$1", word).Scan(&idstring)
		if err != nil {
			return nil
		}
		idstring = strings.Trim(idstring, "{}")
		idStrings := strings.Split(idstring, ",")
		var ids []int
		for _, idStr := range idStrings {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return nil
			}
			ids = append(ids, id)
		}
		for _, id := range ids {
			stat[id]++
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
		comic_url := p.GetUrlByComicId(k)
		comics = append(comics, comic_url)
		if len(comics) >= 10 {
			break
		}
	}
	return comics
}

func (p *PostgreSQL) GetUrlByComicId(id int) string {
	row := p.DB.QueryRow("SELECT url FROM comics WHERE id=$1", id)
	var comic_url string
	err := row.Scan(&comic_url)
	if err != nil {
		log.Fatal("error scan comic_url")
	}
	return comic_url
}

func (p *PostgreSQL) GetComicDatabase() map[int]bool {
	rows, _ := p.DB.Query("SELECT id FROM comics")
	defer rows.Close()

	data := make(map[int]bool)
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal("error scan id")
		}

		data[id] = true
	}
	return data
}

func (p *PostgreSQL) CheckDataBase(ctx context.Context) (int, map[int]bool) {
	res_id := 0
	flag := false
	exist := make(map[int]bool)
	data := p.GetComicDatabase()
	if len(data) == 0 {
		return 1, exist
	}
	select {
	case <-ctx.Done():
		return res_id, exist
	default:
		for comics_id := 1; comics_id < len(data); comics_id++ {
			if comics_id != 404 {
				if _, ok := data[comics_id]; !ok {
					if !flag {
						res_id = comics_id
						flag = true
					}
				} else {
					exist[comics_id] = true
				}
			}
		}
	}
	if res_id == 0 {
		return 0, exist
	}
	fmt.Printf("Missed comics with id %d\n", res_id)
	return res_id, exist
}

func (p *PostgreSQL) SizeDatabase() (int, error) {
	var count int
	err := p.DB.QueryRow("SELECT COUNT(*) FROM comics").Scan(&count)
	if err != nil {
		return 0, nil
	}
	return count, nil
}

func (p *PostgreSQL) GetUserByusername(user *models.User, creds *models.Credentials) error {
	err := p.DB.QueryRow("SELECT username, pass, roles FROM users WHERE username = $1", creds.Username).Scan(&user.Username, &user.Password, &user.Role)
	if err != nil {
		return err
	}
	return nil

}
