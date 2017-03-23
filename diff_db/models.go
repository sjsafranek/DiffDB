package diff_db

import "github.com/sjsafranek/SkeletonDB"

type DiffDb struct {
	File string
	Table    string
	DB       skeleton.Database
}
