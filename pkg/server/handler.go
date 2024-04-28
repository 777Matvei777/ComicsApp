package server

import (
	"context"
	"encoding/json"
	"myapp/pkg/app"
	"myapp/pkg/config"
	"myapp/pkg/models"
	"net/http"
)

type Handler struct {
	Cfg    *config.Config
	Ctx    context.Context
	Comics models.Comic
	Client *app.Client
}

func (h *Handler) getPicsHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search")
	comics := h.Client.SearhDatabase(searchQuery)
	results := models.ImageSearchResult{
		URLs: comics,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *Handler) updateComicsHandler(w http.ResponseWriter, r *http.Request) {
	curr_total := h.Client.SizeDatabase()
	h.Client.Start()
	new_total := h.Client.SizeDatabase()
	h.Comics.New = new_total - curr_total
	h.Comics.Total = new_total
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.Comics)
}
