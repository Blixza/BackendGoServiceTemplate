package user_handler

import (
	user_service "backend-service-template/internal/service/user"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct {
	svc    *user_service.Service
	logger *zap.Logger
}

func NewHandler(svc *user_service.Service, logger *zap.Logger) *Handler {
	return &Handler{svc: svc, logger: logger}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Nickname string `json:"nickname"`
		Password string `json:"password"`
		Discord  string `json:"discord"`
		Email    string `json:"email"`
		Balance  int    `json:"balance"`
		Towns    sql.NullString `json:"towns"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.svc.Register(
		r.Context(), req.Nickname, req.Password, req.Discord, req.Email,
		req.Balance, req.Towns,
	)
	if err != nil {
		h.logger.Error("failed to create user", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("user created", zap.Any("user_id", user.ID))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("failed to get user", zap.Error(err))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	user, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("user not found", zap.Any("id", id), zap.Error(err))
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
