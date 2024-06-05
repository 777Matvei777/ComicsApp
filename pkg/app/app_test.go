package app

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"myapp/pkg/config"
	"myapp/pkg/database"
	"myapp/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (mdb *MockDB) CreateComic(values []models.Item) error {
	args := mdb.Called(values)
	return args.Error(0)
}

func (mdb *MockDB) GetUrlByComicId(id int) string {
	return ""
}

func (mdb *MockDB) GetComicDatabase() map[int]bool {
	return nil
}

func (mdb *MockDB) CheckDataBase(ctx context.Context) (int, map[int]bool) {
	return 0, nil
}

func (mdb *MockDB) SizeDatabase() (int, error) {
	return 0, nil
}

func (mdb *MockDB) GetUserByusername(user *models.User, creds *models.Credentials) error {
	return nil
}

func (mdb *MockDB) BuildIndex() ([]models.KeywordIndex, error) {
	return nil, nil
}

func (mdb *MockDB) CreateIndex(keywordIndices []models.KeywordIndex) error {
	return nil
}

func (mdb *MockDB) GetComicsByQuery(searchQuery []string) []string {
	return nil
}
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

func TestNewClient(t *testing.T) {
	cfg, err := config.New("../../config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	client := NewClient(cfg, 1)
	assert.NotNil(t, client)
	assert.Equal(t, cfg, client.Cfg)
	assert.Equal(t, 1, client.Num)

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
func TestCheckForNewComics(t *testing.T) {
	testHour := 12
	testMin := 48
	testSec := 50
	client := &Client{}
	checkForNewComics(testHour, testMin, testSec, client)

	now := time.Now()
	expectedNextRun := time.Date(now.Year(), now.Month(), now.Day(), testHour, testMin, testSec, 0, now.Location())
	if now.After(expectedNextRun) {
		expectedNextRun = expectedNextRun.Add(24 * time.Hour)
	}
	expectedDuration := expectedNextRun.Sub(now)

	if expectedDuration < 0 {
		t.Errorf("The duration until the next run is negative, got: %v", expectedDuration)
	}
}

func TestCreateDataBase(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.New("../../config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	client := &Client{
		Cfg:   cfg,
		Num:   10,
		Exist: make(map[int]bool),
	}
	mockDB := new(MockDB)
	client.Db = mockDB
	mockDB.On("CreateComic", mock.Anything).Return(nil)
	client.CreateDataBase(ctx)
	mockDB.AssertCalled(t, "CreateComic", mock.Anything)
}
