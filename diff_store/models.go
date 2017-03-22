package diff_store

import (
	"sync"
)

type DiffStore struct {
	Name         string
	CurrentValue string
	Diffs        map[int64]string
	Shards       map[int64]DiffShard
	lock         sync.RWMutex
}

type DiffShard struct {
	Name         string
	CurrentValue string
	Diffs        map[int64]string
	lock         sync.RWMutex
}
