package server

import (
	"encoding/json"
	"myapp/pkg/app"
	"myapp/pkg/config"
	"myapp/pkg/models"
	"net/http"
	"sync"
)

type Handler struct {
	Cfg    *config.Config
	Comics models.Comic
	Client *app.Client
	mu     sync.Mutex
}

func (h *Handler) getPicsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	searchQuery := r.URL.Query().Get("search")
	comics := h.Client.SearhDatabase(searchQuery, ctx)
	results := models.ImageSearchResult{
		URLs: comics,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *Handler) updateComicsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if h.mu.TryLock() {
		curr_total := h.Client.SizeDatabase()
		h.Client.Start(ctx)
		new_total := h.Client.SizeDatabase()
		h.Comics.New = new_total - curr_total
		h.Comics.Total = new_total
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(h.Comics)
		h.mu.Unlock()
	}
}
