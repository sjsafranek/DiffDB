package main

import (
	"sync"
)

import (
	"github.com/boltdb/bolt"
)

type DiffDb struct {
	Filename string
	db       *bolt.DB
}

type DiffStore struct {
	Name         string
	CurrentValue string
	Diffs        map[int64]string
	lock         sync.RWMutex
	Shards       map[int64]DiffShard
}

type DiffShard struct {
	Name         string
	CurrentValue string
	Diffs        map[int64]string
	lock         sync.RWMutex
}
