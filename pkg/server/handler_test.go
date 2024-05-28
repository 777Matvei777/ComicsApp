package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"myapp/pkg/app"
	"myapp/pkg/config"
	"myapp/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

type MockDatabase struct {
	Comics map[int]models.Item
	Users  map[string]models.User
}

func (m *MockDatabase) GetUrlByComicId(id int) string {
	if comic, ok := m.Comics[id]; ok {
		return comic.URL
	}
	return ""
}

func (m *MockDatabase) GetComicDatabase() map[int]bool {
	exist := make(map[int]bool)
	for id := range m.Comics {
		exist[id] = true
	}
	return exist
}

func (m *MockDatabase) CheckDataBase(ctx context.Context) (int, map[int]bool) {
	exist := m.GetComicDatabase()

	return 0, exist
}

func (m *MockDatabase) SizeDatabase() (int, error) {
	return len(m.Comics), nil
}

func (m *MockDatabase) GetUserByusername(user *models.User, creds *models.Credentials) error {
	if u, ok := m.Users[creds.Username]; ok {
		*user = u
		return nil
	}
	var err error
	return err
}

func (m *MockDatabase) CreateComic(items []models.Item) error {

	return nil
}

func (m *MockDatabase) BuildIndex() ([]models.KeywordIndex, error) {

	return nil, nil
}

func (m *MockDatabase) CreateIndex(index []models.KeywordIndex) error {

	return nil
}

func (m *MockDatabase) GetComicsByQuery(query []string) []string {

	return nil
}

func TestLoginWithDb(t *testing.T) {
	cfg := &config.Config{
		Token_max_time: 30, // Пример значения
	}
	us := models.User{
		Username: "testuser",
		Password: "password",
		Role:     "user",
	}
	mp := make(map[string]models.User)
	mp["testuser"] = us

	mockDb := &MockDatabase{
		Users: mp,
	}
	client := &app.Client{
		Cfg:   cfg,
		Num:   0,
		Exist: make(map[int]bool),
		Db:    mockDb,
	}

	handler := &Handler{
		Cfg:    cfg,
		Client: client,
	}

	creds := models.Credentials{
		Username: "testuser",
		Password: "password",
	}

	body, _ := json.Marshal(creds)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.loginHandler(rr, req, client)

	assert.Equal(t, http.StatusOK, rr.Code)
	body, err := io.ReadAll(rr.Body)
	assert.Nil(t, err)
	tokenString := string(body)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil
	})
	assert.Nil(t, err)
	assert.True(t, token.Valid)
}
