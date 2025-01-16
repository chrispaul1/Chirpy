package handlers

import (
	"chrispaul1/chirpy/internal/auth"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (h *UserHandler) HandleWebhooks(w http.ResponseWriter, req *http.Request) {

	type RequestStruct struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIkey(req.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if apiKey != h.cfg.POLKA_KEY {
		RespondWithError(w, http.StatusUnauthorized, "Error, unauthorized access")
		return
	}

	decoder := json.NewDecoder(req.Body)
	userFields := RequestStruct{}
	err = decoder.Decode(&userFields)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Error, bad request")
		return
	}

	if userFields.Event != "user.upgraded" {
		RespondWithError(w, http.StatusNoContent, "")
		return
	}

	parsedID, err := uuid.Parse(userFields.Data.UserID)

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Bad request, id not formatted correctly")
		return
	}
	_, err = h.cfg.DB.UpgradeUserChirpyMembership(req.Context(), parsedID)

	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Error, user not found")
		return
	}

	RespondWithJson(w, http.StatusNoContent, "")

}
