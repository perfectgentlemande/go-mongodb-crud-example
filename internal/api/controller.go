package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Config struct {
	Addr   string `yaml:"addr"`
	Prefix string `yaml:"prefix"`
}
type ServerParams struct {
	Cfg *Config
}

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

func NewServer(params *ServerParams) *http.Server {
	ctrl := NewController()

	router := chi.NewRouter()
	router.Use(NewLoggingMiddleware())
	apiRouter := chi.NewRouter()
	HandlerFromMux(ctrl, apiRouter)

	if params.Cfg.Prefix == "" {
		params.Cfg.Prefix = "/"
	}
	router.Route(params.Cfg.Prefix, func(r chi.Router) {
		r.Mount("/", apiRouter)
	})

	return &http.Server{
		Addr:    params.Cfg.Addr,
		Handler: router,
	}
}

func NewLoggingMiddleware() func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			handler.ServeHTTP(w, r)
		})
	}
}
func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}
func WriteError() {

}
func WriteSuccessful() {

}
