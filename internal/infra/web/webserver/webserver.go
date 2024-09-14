package webserver

import (
	"net/http"

	"github.com/carlosmeds/rate-limiter/internal/infra/database"
	md "github.com/carlosmeds/rate-limiter/internal/infra/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() {
	repo := database.NewRateLimiterRepository()
	rl := md.NewRateLimiterMiddleware(repo)

	s.Router.Use(middleware.Logger)
	s.Router.Use(rl.RateLimiter)
	for path, handler := range s.Handlers {
		s.Router.Handle(path, handler)
	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
