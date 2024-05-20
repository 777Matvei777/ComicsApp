package server

import (
	"encoding/json"
	"fmt"
	"log"
	"myapp/pkg/app"
	"myapp/pkg/config"
	"myapp/pkg/models"
	"net/http"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/time/rate"
)

type Handler struct {
	Cfg     *config.Config
	Comics  models.Comic
	Client  *app.Client
	mu      sync.Mutex
	sem     chan struct{}
	limiter *rate.Limiter
	wg      *sync.WaitGroup
}

func (h *Handler) getPicsHandler(w http.ResponseWriter, r *http.Request) {
	if !h.limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
	}
	h.sem <- struct{}{}
	defer func() { <-h.sem }()
	ctx := r.Context()
	searchQuery := r.URL.Query().Get("search")
	comics := h.Client.SearhDatabase(searchQuery, ctx)
	results := models.ImageSearchResult{
		URLs: comics,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *Handler) loginHandler(w http.ResponseWriter, r *http.Request, client *app.Client) {
	client.LoginWithDb(w, r)
}
func (h *Handler) updateComicsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if h.mu.TryLock() {
		if !h.limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		}
		curr_total, err := h.Client.SizeDatabase()
		if err != nil {
			log.Fatal(err)
		}
		h.Client.Start(ctx)
		new_total, err := h.Client.SizeDatabase()
		if err != nil {
			log.Fatal(err)
		}
		h.Comics.New = new_total - curr_total
		h.Comics.Total = new_total
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(h.Comics)
		h.mu.Unlock()
	}
}

func (h *Handler) authMiddleware(handler http.HandlerFunc, role string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret_key"), nil
		})
		if err != nil {
			http.Error(w, "error parsing token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		if token == nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == role {
				handler(w, r)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
		}
	}
}
