package handler

import (
	"cat-manager/backend/internal/service"
	"net/http"
)

type DashboardHandler struct {
	svc *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) Get(w http.ResponseWriter, r *http.Request) {
	dashboard, err := h.svc.Get()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, dashboard)
}
