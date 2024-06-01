package app

import (
	"bytes"
	"context"
	"encoding/json"
	"myapp/pkg/config"
	"myapp/pkg/database"
	"myapp/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestStart(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT .* FROM comics").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("SELECT .* FROM keywords").WillReturnRows(sqlmock.NewRows([]string{"keyword", "comic_id"}).AddRow("example", 1))

	client := &Client{
		Db: &database.PostgreSQL{DB: db},
	}

	ctx := context.Background()
	client.Start(ctx)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSearhDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	comicUrls := []string{"http://example.com/comic1"}

	client := &Client{
		Db: &database.PostgreSQL{DB: db},
	}
	mock.ExpectQuery("SELECT comic_ids FROM keyword_index WHERE keyword=\\$1").
		WithArgs("flag").
		WillReturnRows(sqlmock.NewRows([]string{"comic_ids"}).AddRow("1"))

	mock.ExpectQuery("SELECT url FROM comics WHERE id=\\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"url"}).AddRow("http://example.com/comic1"))

	ctx := context.Background()
	comics := client.SearhDatabase("flag", ctx)
	if len(comics) != len(comicUrls) {
		t.Errorf("expected %d comics, got %d", len(comicUrls), len(comics))
	}
	for i, comic := range comics {
		if comic != comicUrls[i] {
			t.Errorf("expected comic at index %d to be %s, got %s", i, comicUrls[i], comic)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSizeDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := &Client{
		Db: &database.PostgreSQL{DB: db},
	}

	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM comics").WillReturnRows(rows)

	count, err := client.SizeDatabase()
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

func TestLoginWithDb(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := &Client{
		Db: &database.PostgreSQL{DB: db},
		Cfg: &config.Config{
			Token_max_time: 15,
		},
	}

	creds := models.Credentials{Username: "testuser", Password: "password"}
	user := models.User{Username: "testuser", Password: "password", Role: "user"}

	mock.ExpectQuery("SELECT (.+) FROM users WHERE username=?").
		WithArgs(creds.Username).
		WillReturnRows(sqlmock.NewRows([]string{"username", "password", "role"}).
			AddRow(user.Username, user.Password, user.Role))

	body, _ := json.Marshal(creds)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(client.LoginWithDb)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// func TestCreateDataBase(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("Не удалось создать mock базы данных: %s", err)
// 	}
// 	defer db.Close()
// 	cfg, err := config.New("../../config.yaml")
// 	if err != nil {
// 		log.Println("Ошибка создания конфига")
// 	}
// 	// Настройка моковых вызовов для функции CreateComic
// 	mock.ExpectBegin()

// 	mock.ExpectPrepare("INSERT INTO comics \\(url\\) VALUES \\(\\$1\\) RETURNING id")

// 	mock.ExpectPrepare("INSERT INTO keywords \\(keyword, comic_id\\) VALUES \\(\\$1, \\$2\\)")

// 	mock.ExpectQuery("INSERT INTO comics \\(url\\) VALUES \\(\\$1\\) RETURNING id").
// 		WithArgs("https://imgs.xkcd.com/comics/barrel_cropped_(1).jpg").
// 		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
// 	keyword := "[[A boy sits in a barrel which is floating in an ocean.]]\nBoy: I wonder where I'll float next?\n[[The barrel drifts into the distance. Nothing else can be seen.]]\n{{Alt: Don't we all.}} Don't we all."
// 	normalized := words.Normalize(keyword)
// 	mock.ExpectExec("INSERT INTO keywords \\(keyword, comic_id\\) VALUES \\(\\$1, \\$2\\)").
// 		for i := 0; i < len(normalized); i++{
// 			WithArgs(normalized[i], 1).
// 		}
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	mock.ExpectCommit()
// 	exist := make(map[int]bool)
// 	for i := 2; i < 3000; i++ {
// 		exist[i] = true
// 	}
// 	// Создание клиента с моковой базой данных
// 	c := &Client{
// 		Db:    &database.PostgreSQL{DB: db},
// 		Cfg:   cfg,
// 		Exist: exist,
// 	}

// 	// Вызов функции CreateDataBase
// 	c.CreateDataBase(context.Background())

// 	// Проверка, что все ожидания были выполнены
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("Не все ожидания были выполнены: %s", err)
// 	}
// }
