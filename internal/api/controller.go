package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/perfectgentlemande/go-mongodb-crud-example/internal/logger"
	"github.com/perfectgentlemande/go-mongodb-crud-example/internal/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Addr   string `yaml:"addr"`
	Prefix string `yaml:"prefix"`
}
type ServerParams struct {
	Cfg  *Config
	Log  *zap.Logger
	Srvc *service.Service
}

type Controller struct {
	ServerInterface

	srvc *service.Service
}

func NewController(srvc *service.Service) *Controller {
	return &Controller{
		srvc: srvc,
	}
}

func NewServer(params *ServerParams) *http.Server {
	ctrl := NewController(params.Srvc)

	router := chi.NewRouter()
	router.Use(logger.NewLoggingMiddleware(params.Log))
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

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}
func WriteError(ctx context.Context, w http.ResponseWriter, status int, message string) {
	log := logger.GetLogger(ctx)

	err := RespondWithJSON(w, status, APIError{Message: message})
	if err != nil {
		log.Error("write response", zap.Field{
			Key:       logger.ErrorField,
			Type:      zapcore.ErrorType,
			Interface: err,
		})
	}
}
func WriteSuccessful(ctx context.Context, w http.ResponseWriter, payload interface{}) {
	log := logger.GetLogger(ctx)

	err := RespondWithJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Error("write response", zap.Field{
			Key:       logger.ErrorField,
			Type:      zapcore.ErrorType,
			Interface: err,
		})
	}
}
