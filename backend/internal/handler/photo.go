package handler

import (
	"cat-manager/backend/internal/model"
	"cat-manager/backend/internal/service"
	"net/http"
)

type PhotoHandler struct {
	svc        *service.PhotoService
	petService *service.PetService
}

func NewPhotoHandler(svc *service.PhotoService, petService *service.PetService) *PhotoHandler {
	return &PhotoHandler{svc: svc, petService: petService}
}

func (h *PhotoHandler) List(w http.ResponseWriter, r *http.Request) {
	petID := r.PathValue("id")
	photos, err := h.svc.List(petID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if photos == nil {
		photos = []model.Photo{}
	}
	writeJSON(w, http.StatusOK, photos)
}

func (h *PhotoHandler) Upload(w http.ResponseWriter, r *http.Request) {
	petID := r.PathValue("id")

	pet, err := h.petService.GetByID(petID)
	if err != nil {
		writeError(w, http.StatusNotFound, "pet not found")
		return
	}

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "failed to parse form")
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		writeError(w, http.StatusBadRequest, "photo file is required")
		return
	}
	defer file.Close()

	photoURL, _, err := saveAndProcessImage(file, header, petID, "photos")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	caption := r.FormValue("caption")
	var captionPtr *string
	if caption != "" {
		captionPtr = &caption
	}

	// Accept age_group from frontend; fall back to auto-calculation
	ageGroup := r.FormValue("age_group")
	if ageGroup == "" {
		ageGroup = service.CalculateAgeGroup(pet.Birthday)
	}

	photo, err := h.svc.CreateWithAgeGroup(petID, photoURL, ageGroup, captionPtr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, photo)
}

func (h *PhotoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	petID := r.PathValue("id")
	photoID := r.PathValue("pid")

	if err := h.svc.Delete(petID, photoID); err != nil {
		writeError(w, http.StatusNotFound, "photo not found")
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}
