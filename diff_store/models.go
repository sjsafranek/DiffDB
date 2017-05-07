package diffstore

import (
	"sync"
)

// DiffStore struct for storing diff data
type DiffStore struct {
	// Name         string
	CurrentValue string
	Diffs        map[int64]string
	Shards       map[int64]DiffShard
	lock         sync.RWMutex
}

// DiffShard struct for storing pieces of diff data
type DiffShard struct {
	// Name         string
	CurrentValue string
	Diffs        map[int64]string
	lock         sync.RWMutex
}
