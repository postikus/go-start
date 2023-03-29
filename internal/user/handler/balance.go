package handler

import (
	"github.com/postikus/go-starter/internal/user"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	log     *zap.Logger
	service *user.Service
}

func NewHandler(log *zap.Logger, service *user.Service) *Handler {
	return &Handler{
		log:     log.Named("handler.user"),
		service: service,
	}
}

func (h *Handler) ServeHTTP(r *http.Request, w http.ResponseWriter) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
