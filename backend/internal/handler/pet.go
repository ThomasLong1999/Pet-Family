package handler

import (
	"cat-manager/backend/internal/model"
	"cat-manager/backend/internal/service"
	"encoding/json"
	"net/http"
)

type PetHandler struct {
	svc *service.PetService
}

func NewPetHandler(svc *service.PetService) *PetHandler {
	return &PetHandler{svc: svc}
}

func (h *PetHandler) List(w http.ResponseWriter, r *http.Request) {
	pets, err := h.svc.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if pets == nil {
		pets = []model.Pet{}
	}
	writeJSON(w, http.StatusOK, pets)
}

func (h *PetHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	pet, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "pet not found")
		return
	}
	writeJSON(w, http.StatusOK, pet)
}

func (h *PetHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreatePetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.Birthday == "" {
		writeError(w, http.StatusBadRequest, "birthday is required")
		return
	}
	if req.Species == "" {
		req.Species = "cat"
	}

	pet, err := h.svc.Create(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, pet)
}

func (h *PetHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req model.UpdatePetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	pet, err := h.svc.Update(id, req)
	if err != nil {
		writeError(w, http.StatusNotFound, "pet not found")
		return
	}
	writeJSON(w, http.StatusOK, pet)
}

func (h *PetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, "pet not found")
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}

func (h *PetHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "failed to parse form")
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		writeError(w, http.StatusBadRequest, "avatar file is required")
		return
	}
	defer file.Close()

	avatarURL, dominantColor, err := saveAndProcessImage(file, header, id, "avatars")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.svc.UpdateAvatar(id, avatarURL, dominantColor); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"avatar_url":      avatarURL,
		"dominant_color":  dominantColor,
	})
}
