package config

import (
	"chrispaul1/chirpy/internal/database"
	"sync/atomic"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
	Platform       string
}
