package main

import (
	"github.com/boltdb/bolt"
)

type DiffDb struct {
	Filename string
	db       *bolt.DB
}

// type versionsInfo struct {
// 	VersionDate string
// 	VersionNum  int
// }

type DiffStore struct {
	Title       string
	CurrentText string
	//Diffs       []string
	//Timestamps  []int64
	Diffs map[int64]string
	//Encrypted bool
}
