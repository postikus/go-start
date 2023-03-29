package handler

import (
	"encoding/json"
	"github.com/postikus/go-starter/internal/user"
	"github.com/postikus/go-starter/model"
	"go.uber.org/zap"
	"net/http"
)

type Post struct {
	log     *zap.Logger
	service *user.Service
}

func NewPost(log *zap.Logger, service *user.Service) *Post {
	return &Post{
		log:     log.Named("handler.user"),
		service: service,
	}
}

func (h *Post) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		user = new(model.User)
		err  error
	)

	if err = json.NewDecoder(r.Body).Decode(user); err != nil {
		h.log.Warn("Error decoding request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user, err = h.service.New(r.Context(), user); err != nil {
		h.log.Error("Error creating new user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		h.log.Warn("Error encoding response", zap.Error(err))
		return
	}
}
