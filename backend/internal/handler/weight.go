package handler

import (
	"cat-manager/backend/internal/model"
	"cat-manager/backend/internal/service"
	"encoding/json"
	"net/http"
)

type WeightHandler struct {
	svc *service.WeightService
}

func NewWeightHandler(svc *service.WeightService) *WeightHandler {
	return &WeightHandler{svc: svc}
}

func (h *WeightHandler) List(w http.ResponseWriter, r *http.Request) {
	petID := r.PathValue("id")
	records, err := h.svc.List(petID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if records == nil {
		records = []model.WeightRecord{}
	}
	writeJSON(w, http.StatusOK, records)
}

func (h *WeightHandler) Create(w http.ResponseWriter, r *http.Request) {
	petID := r.PathValue("id")

	var req model.CreateWeightRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Weight <= 0 {
		writeError(w, http.StatusBadRequest, "weight must be positive")
		return
	}
	if req.RecordedAt == "" {
		writeError(w, http.StatusBadRequest, "recorded_at is required")
		return
	}

	record, err := h.svc.Create(petID, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, record)
}

func (h *WeightHandler) Delete(w http.ResponseWriter, r *http.Request) {
	petID := r.PathValue("id")
	recordID := r.PathValue("rid")

	if err := h.svc.Delete(petID, recordID); err != nil {
		writeError(w, http.StatusNotFound, "record not found")
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}
