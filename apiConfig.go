package main

import "sync/atomic"

type apiConfig struct {
	fsHit atomic.Int32
}