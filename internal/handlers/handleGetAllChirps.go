package handlers

import (
	"chrispaul1/chirpy/internal/database"
	"net/http"
	"sort"
	"strings"

	"github.com/google/uuid"
)

func (h *UserHandler) HandleGetAllChirps(w http.ResponseWriter, req *http.Request) {

	authorIDString := req.URL.Query().Get("author_id")
	sortString := req.URL.Query().Get("sort")
	var userChirps []database.Chirp
	var err error

	authorID := uuid.Nil
	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid Author ID ")
		}
	}

	allChirps := make([]Chirp, 0, len(userChirps))
	userChirps, err = h.cfg.DB.GetAllChirps(req.Context())
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Error, could not retrieve the chirps")
		return
	}

	for _, itemChirp := range userChirps {
		if authorID != uuid.Nil && itemChirp.UserID != authorID {
			continue
		}
		newChirp := Chirp{
			ID:        itemChirp.ID,
			CreatedAt: itemChirp.CreatedAt,
			UpdatedAt: itemChirp.UpdatedAt,
			Body:      itemChirp.Body,
			UserID:    itemChirp.UserID,
		}
		allChirps = append(allChirps, newChirp)
	}

	if strings.ToLower(sortString) == "desc" {
		sort.Slice(allChirps, func(i, j int) bool {
			return allChirps[i].CreatedAt.After(allChirps[j].CreatedAt)
		})
	} else {
		sort.Slice(allChirps, func(i, j int) bool {
			return allChirps[i].CreatedAt.Before(allChirps[j].CreatedAt)
		})
	}
	RespondWithJson(w, http.StatusOK, allChirps)
}
