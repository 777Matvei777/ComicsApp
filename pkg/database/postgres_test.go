package database

import (
	"context"
	"fmt"
	"myapp/pkg/models"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCheckDataBaseWithMissedComic(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	p := &PostgreSQL{DB: db}
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1).
		AddRow(2).
		AddRow(3).
		AddRow(5).
		AddRow(6).
		AddRow(7)
	mock.ExpectQuery("SELECT id FROM comics").WillReturnRows(rows)
	res_id, _ := p.CheckDataBase(ctx)
	expectedResID := 4
	if res_id != expectedResID {
		t.Errorf("CheckDataBase() res_id = %d, want %d", res_id, expectedResID)
	}
}
func TestCheckDataBase(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	p := &PostgreSQL{DB: db}

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1).
		AddRow(2).
		AddRow(3)
	mock.ExpectQuery("SELECT id FROM comics").WillReturnRows(rows)

	res_id, exist := p.CheckDataBase(ctx)

	// Проверяем, что функция возвращает правильный ID первого пропущенного комикса.
	expectedResID := 0
	if res_id != expectedResID {
		t.Errorf("CheckDataBase() res_id = %d, want %d", res_id, expectedResID)
	}
	expectedExist := map[int]bool{1: true, 2: true}
	for id, e := range expectedExist {
		if exist[id] != e {
			t.Errorf("CheckDataBase() exist[%d] = %v, want %v", id, exist[id], e)
		}
	}

	// Проверяем, что функция не возвращает несуществующие комиксы.
	if len(exist) != len(expectedExist) {
		t.Errorf("CheckDataBase() returned unexpected number of existing comics: got %v, want %v", len(exist), len(expectedExist))
	}
}
func TestCreateComic(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := &PostgreSQL{DB: db}

	mock.ExpectBegin()

	mock.ExpectPrepare("INSERT INTO comics \\(url\\) VALUES \\(\\$1\\) RETURNING id")

	mock.ExpectPrepare("INSERT INTO keywords \\(keyword, comic_id\\) VALUES \\(\\$1, \\$2\\)")

	mock.ExpectQuery("INSERT INTO comics \\(url\\) VALUES \\(\\$1\\) RETURNING id").
		WithArgs("http://example.com/comic1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec("INSERT INTO keywords \\(keyword, comic_id\\) VALUES \\(\\$1, \\$2\\)").
		WithArgs("hero", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	comics := []models.Item{
		{
			URL:      "http://example.com/comic1",
			Keywords: []string{"hero"},
		},
	}

	err = repo.CreateComic(comics)
	if err != nil {
		t.Errorf("error was not expected while creating comics: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateIndex(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := &PostgreSQL{DB: db}

	mock.ExpectBegin()

	stmt := "INSERT INTO keyword_index \\(keyword, comic_ids\\) VALUES \\(\\$1, \\$2\\) ON CONFLICT \\(keyword\\) DO UPDATE SET comic_ids = excluded.comic_ids"
	mock.ExpectPrepare(stmt)

	keywordIndices := []models.KeywordIndex{
		{Keyword: "hero", Index: []int{1, 2, 3}},
	}

	for _, ki := range keywordIndices {
		indices := "{" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ki.Index)), ","), "[]") + "}"
		mock.ExpectExec(stmt).
			WithArgs(ki.Keyword, indices).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	mock.ExpectCommit()

	err = repo.CreateIndex(keywordIndices)
	if err != nil {
		t.Errorf("error was not expected while creating index: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetComicsByQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := &PostgreSQL{DB: db}

	mock.ExpectQuery("SELECT comic_ids FROM keyword_index WHERE keyword=\\$1").
		WithArgs("hero").
		WillReturnRows(sqlmock.NewRows([]string{"comic_ids"}).AddRow("1"))

	mock.ExpectQuery("SELECT url FROM comics WHERE id=\\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"url"}).AddRow("http://example.com/comic1"))

	searchQuery := []string{"hero"}
	comics := repo.GetComicsByQuery(searchQuery)
	expectedComics := []string{
		"http://example.com/comic1",
	}
	if !reflect.DeepEqual(comics, expectedComics) {
		t.Errorf("expected %v, got %v", expectedComics, comics)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUrlByComicId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repos := &PostgreSQL{DB: db}

	comicID := 1
	comicURL := "http://example.com/comic1"
	rows := sqlmock.NewRows([]string{"url"}).AddRow(comicURL)
	mock.ExpectQuery("SELECT url FROM comics WHERE id=\\$1").
		WithArgs(comicID).
		WillReturnRows(rows)

	url := repos.GetUrlByComicId(comicID)

	assert.Equal(t, comicURL, url)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSizeDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	p := &PostgreSQL{DB: db}

	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM comics").WillReturnRows(rows)

	count, err := p.SizeDatabase()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if count != 5 {
		t.Errorf("expected count to be 5, got %d", count)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUserByusername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	p := &PostgreSQL{DB: db}
	creds := &models.Credentials{Username: "testuser"}
	user := &models.User{}

	rows := sqlmock.NewRows([]string{"username", "pass", "roles"}).
		AddRow(creds.Username, "hashedpassword", "user")
	mock.ExpectQuery("SELECT username, pass, roles FROM users WHERE username = \\$1").
		WithArgs(creds.Username).
		WillReturnRows(rows)

	err = p.GetUserByusername(user, creds)

	if err != nil {
		t.Errorf("error was not expected while getting user: %s", err)
	}
	if user.Username != creds.Username || user.Password != "hashedpassword" || user.Role != "user" {
		t.Errorf("expected user to have username %s and role 'user', got %v", creds.Username, user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetComicDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	p := &PostgreSQL{DB: db}

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1).
		AddRow(2).
		AddRow(3)
	mock.ExpectQuery("SELECT id FROM comics").WillReturnRows(rows)

	data := p.GetComicDatabase()

	expectedData := map[int]bool{1: true, 2: true, 3: true}
	if len(data) != len(expectedData) {
		t.Errorf("expected data length to be %d, got %d", len(expectedData), len(data))
	}
	for id, exists := range expectedData {
		if data[id] != exists {
			t.Errorf("expected data for id %d to be %v, got %v", id, exists, data[id])
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestBuildIndex(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	p := &PostgreSQL{DB: db}

	rows := sqlmock.NewRows([]string{"keyword", "comic_id"}).
		AddRow("hero", 1).
		AddRow("villain", 2).
		AddRow("hero", 3)
	mock.ExpectQuery("SELECT k\\.keyword, k\\.comic_id FROM keywords k JOIN comics c ON k\\.comic_id = c\\.id").WillReturnRows(rows)

	keywordIndices, err := p.BuildIndex()

	expectedKeywordIndices := []models.KeywordIndex{
		{Keyword: "hero", Index: []int{1, 3}},
		{Keyword: "villain", Index: []int{2}},
	}
	if err != nil {
		t.Errorf("error was not expected while building index: %s", err)
	}
	if !reflect.DeepEqual(keywordIndices, expectedKeywordIndices) {
		t.Errorf("expected keyword indices to be %v, got %v", expectedKeywordIndices, keywordIndices)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestNewPostgreSQL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	pg, err := NewPostgreSQL("host=localhost dbname=postgres user=postgres port=5432 password=local sslmode=disable", "file://../../migrations")

	if err != nil {
		t.Errorf("error was not expected while creating PostgreSQL: %s", err)
	}

	if pg == nil {
		t.Errorf("expected non-nil PostgreSQL object")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
