package diff_db

import "github.com/sjsafranek/SkeletonDB"

type DiffDb struct {
	Filename string
	Table    string
	DB       skeleton.Database
}
