package handlers

import (
	"chrispaul1/chirpy/internal/auth"
	"net/http"

	"github.com/google/uuid"
)

func (h *UserHandler) HandleDeleteChirp(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(id)

	if err != nil {
		errMsg := "Error, ID not acceptable"
		RespondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	chirp, err := h.cfg.DB.GetChirpByID(req.Context(), chirpID)

	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Error, chirp could not be found")
		return
	}

	tokenString, err := auth.GetBearerToken(req.Header)

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	jwtUserID, err := auth.ValidateJWT(tokenString, h.cfg.JWT_SECRET)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if chirp.UserID != jwtUserID {
		RespondWithError(w, http.StatusForbidden, "Not allowed to delete the chirp")
		return
	}

	_, err = h.cfg.DB.DeleteChirp(req.Context(), chirpID)
	if err != nil {
		RespondWithError(w, http.StatusForbidden, "Could not delete the chirp")
		return
	}

	RespondWithJson(w, http.StatusNoContent, "")
}
