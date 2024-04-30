package server

import (
	"context"
	"myapp/pkg/app"
	"myapp/pkg/config"
	"myapp/pkg/models"
	"net/http"
)

type Server struct {
	Router *http.ServeMux
	Cfg    *config.Config
}

func NewServer(cfg *config.Config, ctx context.Context) *Server {
	s := &Server{
		Router: http.NewServeMux(),
		Cfg:    cfg,
	}
	s.initHandlers(ctx)
	return s
}

func (s *Server) initHandlers(ctx context.Context) {
	comics := models.NewComic()
	client := app.NewClient(s.Cfg, ctx, 1)
	h := Handler{
		Cfg:    s.Cfg,
		Comics: *comics,
		Client: client,
	}
	s.Router.HandleFunc("GET /pics", h.getPicsHandler)
	s.Router.HandleFunc("POST /update", h.updateComicsHandler)
}

func (s *Server) RunServer() {
	http.ListenAndServe(s.Cfg.Port, s.Router)
}
