package server

import (
	"log"
	"myapp/pkg/app"
	"myapp/pkg/config"
	"myapp/pkg/models"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
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
	var wg sync.WaitGroup
	comics := models.NewComic()
	client := app.NewClient(s.Cfg, 1)
	lim := rate.NewLimiter(rate.Limit(s.Cfg.RateLimit), 1)
	h := Handler{
		Cfg:     s.Cfg,
		Comics:  *comics,
		Client:  client,
		sem:     make(chan struct{}, s.Cfg.ConcurrencyLimit),
		limiter: lim,
		wg:      &wg,
	}
	s.Router.HandleFunc("GET /pics", h.getPicsHandler)
	s.Router.HandleFunc("POST /update", h.authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		h.updateComicsHandler(w, r)
	}, "admin"))
	s.Router.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		h.loginHandler(w, r, client)
	})
	h.wg.Wait()
}

func (s *Server) RunServer() {
	err := s.Serv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
