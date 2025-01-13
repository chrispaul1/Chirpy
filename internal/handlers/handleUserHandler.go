package handlers

import "chrispaul1/chirpy/internal/config"

type UserHandler struct {
	cfg *config.ApiConfig // or whatever package ApiConfig is in
}

func NewUserHandler(cfg *config.ApiConfig) *UserHandler {
	return &UserHandler{
		cfg: cfg,
	}
}
