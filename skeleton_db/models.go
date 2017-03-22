package skeleton_db

import (
	"github.com/boltdb/bolt"
)

const (
	VERSION string = "0.0.1"
)

type DiffDb struct {
	Filename string
	Table    string
	db       *bolt.DB
}

/*

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

*/
