package diff_db

import (
	"strings"
)

import "github.com/sjsafranek/SkeletonDB"

func New(db_file string) DiffDb {
	db := DiffDb{File: db_file}
	db.Init()
	return db
}

func (self *DiffDb) Init() {

	self.DB = skeleton.Database{File: self.getFile()}
	self.DB.Init()

	err := self.DB.CreateTable(self.getTable())
	if nil != err {
		panic(err)
	}
}

func (self *DiffDb) getFile() string {
	if "" == self.File {
		return "diff.db"
	}
	return self.File
}

func (self *DiffDb) getTable() string {
	if "" == self.Table {
		return "DiffData"
	}
	return self.Table
}

func (self *DiffDb) Load(name string) ([]byte, error) {
	name = strings.ToLower(name)
	data, err := self.DB.Select(self.getTable(), name)
	return data, err
}

func (self *DiffDb) Save(name string, data []byte) error {
	err := self.DB.Insert(self.getTable(), name, data)
	return err
}

func (self *DiffDb) Remove(name string) error {
	err := self.DB.Remove(name, self.getTable())
	return err
}

func (self *DiffDb) SelectAll() ([]string, error) {
	data, err := self.DB.SelectAll(self.getTable())
	return data, err
}
