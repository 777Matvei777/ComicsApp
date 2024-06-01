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
	"sync"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/time/rate"
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

type MockClient struct {
	mock.Mock
}

func (m *MockClient) SizeDatabase() (int, error) {
	return 0, nil
}

func (m *MockClient) Start(ctx context.Context) {
	m.Called(ctx)
}

func (m *MockClient) CreateDataBase(ctx context.Context) {

}

func (m *MockClient) SearhDatabase(searchFlag string, ctx context.Context) []string {
	if searchFlag == "test" {
		return []string{"http://example.com/comic1.jpg", "http://example.com/comic2.jpg"}
	}

	return []string{}
}

func (m *MockClient) CheckDataBase(ctx context.Context) (int, map[int]bool) {
	return 0, nil
}

func (m *MockClient) LoginWithDb(w http.ResponseWriter, r *http.Request) {

}

func TestGetPicsHandler(t *testing.T) {
	cfg := &config.Config{}
	client := &MockClient{}
	limiter := rate.NewLimiter(1, 1)
	h := &Handler{
		Cfg:     cfg,
		Client:  client,
		limiter: limiter,
		sem:     make(chan struct{}, 1),
	}

	req, err := http.NewRequest("GET", "/pics?search=test", nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	h.getPicsHandler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := `{"urls": ["http://example.com/comic1.jpg", "http://example.com/comic2.jpg"]}`
	assert.JSONEq(t, expected, rr.Body.String())

	limiter.SetBurst(0)
	reqOverLimit, _ := http.NewRequest("GET", "/pics?search=test", nil)
	rrOverLimit := httptest.NewRecorder()
	h.getPicsHandler(rrOverLimit, reqOverLimit)
	assert.Equal(t, http.StatusTooManyRequests, rrOverLimit.Code)

	reqInvalidSearch, _ := http.NewRequest("GET", "/pics?search=", nil)
	rrInvalidSearch := httptest.NewRecorder()
	h.getPicsHandler(rrInvalidSearch, reqInvalidSearch)
	assert.Equal(t, http.StatusTooManyRequests, rrInvalidSearch.Code)
}
func TestUpdateComicsHandlerWithMock(t *testing.T) {
	mockClient := new(MockClient)
	mockClient.On("Start", mock.Anything).Once()
	h := &Handler{
		Cfg:     &config.Config{},
		Comics:  models.Comic{},
		Client:  mockClient,
		mu:      sync.Mutex{},
		sem:     make(chan struct{}, 1),
		limiter: rate.NewLimiter(rate.Every(1*time.Second), 1),
		wg:      &sync.WaitGroup{},
	}
	req, err := http.NewRequest("POST", "/update", nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	h.updateComicsHandler(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := "{\"new\":0,\"total\":0}"
	assert.JSONEq(t, expected, rr.Body.String(), "handler returned unexpected body")
	mockClient.AssertExpectations(t)
}

func TestAuthMiddleware(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	middleware := (&Handler{}).authMiddleware(testHandler, "admin")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "admin",
	})
	tokenString, _ := token.SignedString([]byte("secret_key"))
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", tokenString)
	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddlewareUnauthorized(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	middleware := (&Handler{}).authMiddleware(testHandler, "admin")
	tokenWithWrongKey := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "admin",
	})
	wrongKeyTokenString, _ := tokenWithWrongKey.SignedString([]byte("wrong_key"))

	reqWithWrongKey, _ := http.NewRequest("GET", "/", nil)
	reqWithWrongKey.Header.Set("Authorization", wrongKeyTokenString)

	rrWithWrongKey := httptest.NewRecorder()
	middleware.ServeHTTP(rrWithWrongKey, reqWithWrongKey)

	assert.Equal(t, http.StatusUnauthorized, rrWithWrongKey.Code)
	tokenWithWrongRole := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "user",
	})
	wrongRoleTokenString, _ := tokenWithWrongRole.SignedString([]byte("secret_key"))

	reqWithWrongRole, _ := http.NewRequest("GET", "/", nil)
	reqWithWrongRole.Header.Set("Authorization", wrongRoleTokenString)

	rrWithWrongRole := httptest.NewRecorder()
	middleware.ServeHTTP(rrWithWrongRole, reqWithWrongRole)

	assert.Equal(t, http.StatusUnauthorized, rrWithWrongRole.Code)

	reqWithoutToken, _ := http.NewRequest("GET", "/", nil)

	rrWithoutToken := httptest.NewRecorder()
	middleware.ServeHTTP(rrWithoutToken, reqWithoutToken)

	assert.Equal(t, http.StatusUnauthorized, rrWithoutToken.Code)
}
