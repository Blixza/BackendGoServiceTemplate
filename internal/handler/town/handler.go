package town_handler

import (
	town_service "backend-service-template/internal/service/town"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct {
	svc    *town_service.Service
	logger *zap.Logger
}

func NewHandler(svc *town_service.Service, logger *zap.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID              uuid.UUID `json:"id"`
		Name            string    `json:"name"`
		Balance         int       `json:"balance"`
		OwnerNickname   string    `json:"owner_nickname"`
		XCoordOverworld int       `json:"x_coord_overworld"`
		YCoordOverworld int       `json:"y_coord_overworld"`
		ZCoordOverworld int       `json:"z_coord_overworld"`
		XCoordNether    int       `json:"x_coord_nether"`
		YCoordNether    int       `json:"y_coord_nether"`
		ZCoordNether    int       `json:"z_coord_nether"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	town, err := h.svc.Register(
		r.Context(), req.Name, req.Balance, req.OwnerNickname,
		req.XCoordOverworld, req.YCoordOverworld, req.ZCoordOverworld,
		req.XCoordNether, req.YCoordNether, req.ZCoordNether,
	)
	if err != nil {
		h.logger.Error("failed to create town", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("town created", zap.Any("town_id", town.ID))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(town)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.FromBytes([]byte(idStr))
	if err != nil {
		h.logger.Error("failed to get town", zap.Error(err))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	town, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("town not found", zap.Any("id", id), zap.Error(err))
		http.Error(w, "town not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(town)
}
