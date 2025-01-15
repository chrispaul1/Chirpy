package handlers

import (
	"chrispaul1/chirpy/internal/auth"
	"net/http"
	"strings"
	"time"
)

func (h *UserHandler) HandlerRefresh(w http.ResponseWriter, req *http.Request) {
	type token struct {
		Token string `json:"token"`
	}

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

	refreshToken, err := h.cfg.DB.GetRefreshToken(req.Context(), refreshTokenStr)

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		RespondWithError(w, http.StatusUnauthorized, "Error, token is expired")
		return
	}

	if refreshToken.RevokedAt.Valid {
		RespondWithError(w, http.StatusUnauthorized, "Error, token is already revoked")
		return
	}

	newAccessToken, err := auth.MakeJWT(refreshToken.UserID, h.cfg.JWT_SECRET)

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	responseTokenStruct := token{
		Token: newAccessToken,
	}
	RespondWithJson(w, http.StatusOK, responseTokenStruct)
}
