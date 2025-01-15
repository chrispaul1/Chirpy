package handlers

import (
	"chrispaul1/chirpy/internal/database"
	"database/sql"
	"net/http"
	"strings"
	"time"
)

func (h *UserHandler) HandleRevoke(w http.ResponseWriter, req *http.Request) {

	authHeader := req.Header.Get("Authorization")

	if authHeader == "" {
		RespondWithError(w, http.StatusUnauthorized, "Error,Header not found")
		return
	}

	refreshTokenStr, ok := strings.CutPrefix(authHeader, "Bearer ")
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "Token could not be found")
		return
	}

	revokeStruct := database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Token: refreshTokenStr,
	}
	err := h.cfg.DB.RevokeRefreshToken(req.Context(), revokeStruct)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	RespondWithJson(w, http.StatusNoContent, "")
}
