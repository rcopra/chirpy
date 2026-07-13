package main

import (
	"sync/atomic"

	"github.com/rcopra/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}
