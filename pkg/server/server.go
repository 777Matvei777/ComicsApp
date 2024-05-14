package server

import (
	"myapp/pkg/app"
	"myapp/pkg/config"
	"myapp/pkg/models"
	"net/http"
)

type Server struct {
	Router *http.ServeMux
	Cfg    *config.Config
	Serv   *http.Server
}

func NewServer(cfg *config.Config) *Server {
	s := &Server{
		Router: http.NewServeMux(),
		Cfg:    cfg,
	}
	s.initHandlers()
	s.Serv = &http.Server{
		Addr:    s.Cfg.Port,
		Handler: s.Router,
	}
	return s
}

func (s *Server) initHandlers() {
	comics := models.NewComic()
	client := app.NewClient(s.Cfg, 1)
	h := Handler{
		Cfg:    s.Cfg,
		Comics: *comics,
		Client: client,
	}
	s.Router.HandleFunc("GET /pics", h.getPicsHandler)
	s.Router.HandleFunc("POST /update", h.updateComicsHandler)
}

func (s *Server) RunServer() {
	s.Serv.ListenAndServe()
}
