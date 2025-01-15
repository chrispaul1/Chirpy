package handlers

import (
	"net/http"

	"github.com/google/uuid"
)

func (h *UserHandler) HandleGetChirpByID(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(id)

	if err != nil {
		errMsg := "Error, ID not acceptable"
		RespondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	userChirp, err := h.cfg.DB.GetChirpByID(req.Context(), chirpID)
	if err != nil {
		errMsg := "Error, could not retrieve the chirp"
		RespondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	newChirp := Chirp{
		ID:        userChirp.ID,
		CreatedAt: userChirp.CreatedAt,
		UpdatedAt: userChirp.UpdatedAt,
		Body:      userChirp.Body,
		UserID:    userChirp.UserID,
	}

	RespondWithJson(w, http.StatusOK, newChirp)

}
