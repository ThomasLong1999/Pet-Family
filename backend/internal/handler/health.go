package handler

import (
	"cat-manager/backend/internal/model"
	"cat-manager/backend/internal/service"
	"encoding/json"
	"net/http"
)

type HealthHandler struct {
	svc *service.HealthService
}

func NewHealthHandler(svc *service.HealthService) *HealthHandler {
	return &HealthHandler{svc: svc}
}

func (h *HealthHandler) List(w http.ResponseWriter, r *http.Request) {
	petID := r.PathValue("id")
	recordType := r.URL.Query().Get("type")

	records, err := h.svc.List(petID, recordType)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if records == nil {
		records = []model.HealthRecord{}
	}
	writeJSON(w, http.StatusOK, records)
}

func (h *HealthHandler) Create(w http.ResponseWriter, r *http.Request) {
	petID := r.PathValue("id")

	// Support both JSON and multipart (for report upload)
	contentType := r.Header.Get("Content-Type")
	var req model.CreateHealthRequest
	var reportURL string

	if len(contentType) > 0 && contentType[:10] == "multipart/" {
		// Multipart form — may include a report file
		if err := r.ParseMultipartForm(20 << 20); err != nil {
			writeError(w, http.StatusBadRequest, "failed to parse form")
			return
		}

		req.Type = r.FormValue("type")
		req.Name = r.FormValue("name")
		req.Date = r.FormValue("date")

		if nd := r.FormValue("next_date"); nd != "" {
			req.NextDate = &nd
		}
		if note := r.FormValue("note"); note != "" {
			req.Note = &note
		}

		// Handle report file upload
		file, header, err := r.FormFile("report")
		if err == nil {
			defer file.Close()
			url, _, saveErr := saveAndProcessImage(file, header, petID, "reports")
			if saveErr == nil {
				reportURL = url
			}
		}
	} else {
		// JSON body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
	}

	if req.Type == "" || req.Name == "" || req.Date == "" {
		writeError(w, http.StatusBadRequest, "type, name, and date are required")
		return
	}

	record, err := h.svc.Create(petID, req, reportURL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, record)
}

func (h *HealthHandler) Update(w http.ResponseWriter, r *http.Request) {
	recordID := r.PathValue("rid")

	var req model.UpdateHealthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	record, err := h.svc.Update(recordID, req)
	if err != nil {
		writeError(w, http.StatusNotFound, "record not found")
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (h *HealthHandler) Delete(w http.ResponseWriter, r *http.Request) {
	petID := r.PathValue("id")
	recordID := r.PathValue("rid")

	if err := h.svc.Delete(petID, recordID); err != nil {
		writeError(w, http.StatusNotFound, "record not found")
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}
