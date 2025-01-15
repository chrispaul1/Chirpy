package handlers

import (
	"net/http"
)

func (h *UserHandler) HandleGetAllChirps(w http.ResponseWriter, req *http.Request) {

	userChirps, err := h.cfg.DB.GetAllChirps(req.Context())
	allChirps := make([]Chirp, 0, len(userChirps))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Error, could not retrieve the chirps")
		return
	}

	for _, item := range userChirps {
		newChirp := Chirp{
			ID:        item.ID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			Body:      item.Body,
			UserID:    item.UserID,
		}
		allChirps = append(allChirps, newChirp)
	}
	RespondWithJson(w, http.StatusOK, allChirps)
}
