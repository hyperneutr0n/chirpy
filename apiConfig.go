package main

import (
	"sync/atomic"

	"github.com/hyperneutr0n/chirpy/internal/database"
)

type apiConfig struct {
	fsHit	atomic.Int32
	db		*database.Queries
}