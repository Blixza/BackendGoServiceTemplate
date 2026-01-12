package user_handler

import (
	user_service "backend-service-template/internal/service/user"
	"encoding/json"
	"net/http"
	"strconv"

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
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.svc.Register(r.Context(), req.Name, req.Email)
	if err != nil {
		h.logger.Error("failed to create user", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("user created", zap.Int("user_id", user.ID))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.svc.Get(r.Context(), id)
	if err != nil {
		h.logger.Error("user not found", zap.Int("id", id), zap.Error(err))
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
