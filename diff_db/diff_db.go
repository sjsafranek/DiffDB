package diff_db

import (
	"fmt"
	"strings"
)

import "github.com/sjsafranek/SkeletonDB"

func NewDiffDb(db_file string) DiffDb {
	var diffDb = DiffDb{
		Filename: db_file,
		Table:    "DiffData",
		DB:       skeleton.Database{File: db_file}}
	diffDb.Init()
	return diffDb
}

func (self DiffDb) Init() {
	if "" == self.Table {
		self.Table = "DiffData"
	}

	self.DB.Init()

	conn := self.DB.Connect()
	defer conn.Close()

	err := self.DB.CreateTable(conn, self.Table)
	if nil != err {
		panic(err)
	}
}

func (self *DiffDb) Load(name string) ([]byte, error) {
	name = strings.ToLower(name)
	data, err := self.DB.Select(self.Table, name)
	return data, err
}

func (self *DiffDb) Save(name string, data []byte) error {
	err := self.DB.Insert(self.Table, name, data)

	fmt.Println(err)

	return err
}

func (self *DiffDb) Remove(name string) error {
	err := self.DB.Remove(name, self.Table)
	return err
}

func (self *DiffDb) SelectAll() ([]string, error) {
	data, err := self.DB.SelectAll(self.Table)
	return data, err
}
