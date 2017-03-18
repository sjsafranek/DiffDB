package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

import "github.com/boltdb/bolt"

// Creates and initializes DiffDb
func NewDiffDb(db_file string) DiffDb {
	var diffDb = DiffDb{Filename: db_file}
	diffDb.Init()
	return diffDb
}

// Open to create the database and open
func (self DiffDb) Open() *bolt.DB {
	var err error
	config := &bolt.Options{Timeout: 30 * time.Second}
	conn, err := bolt.Open(self.Filename, 0600, config)
	if err != nil {
		log.Println("Opening BoltDB timed out")
		log.Fatal(err)
	}
	return conn
}

func (self DiffDb) Init() {
	err := self.CreateTable("datas")
	if nil != err {
		log.Fatal(err)
	}
}

func (self DiffDb) CreateTable(table string) error {
	self.db = self.Open()
	defer self.db.Close()
	err := self.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(table))
		return err
	})
	return err
}

// Close database
func (self DiffDb) Close() {
	self.db.Close()
}

func (self *DiffDb) Load(title string) (DiffData, error) {
	self.db = self.Open()
	defer self.db.Close()

	title = strings.ToLower(title)

	log.Println("[DiffDb] [DEBUG] Searching for key:", title)
	var ddata DiffData
	err := self.db.View(func(tx *bolt.Tx) error {
		//var err error
		bucket := tx.Bucket([]byte("datas"))
		if bucket == nil {
			//return fmt.Errorf("Bucket does not exist")
			panic(fmt.Errorf("Bucket does not exist"))
		}

		k := []byte(title)
		val := bucket.Get(k)

		if val == nil {
			// make new one
			ddata.Title = title
			ddata.CurrentText = ""
			ddata.Diffs = []string{}
			ddata.Timestamps = []int64{}
			return nil
		}

		err := ddata.decode(val)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Could not get DiffData: %s", err)
		return ddata, err
	}
	return ddata, nil
}

func (self *DiffDb) Save(ddata DiffData) error {
	self.db = self.Open()
	defer self.db.Close()

	err := self.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("datas"))
		if bucket == nil {
			//return fmt.Errorf("Bucket does not exist")
			panic(fmt.Errorf("Bucket does not exist"))
		}

		enc, err := ddata.encode()
		if err != nil {
			return fmt.Errorf("could not encode DiffData: %s", err)
		}

		err = bucket.Put([]byte(ddata.Title), enc)
		if err != nil {
			return fmt.Errorf("could add to bucket: %s", err)
		}
		return err
	})
	return err
}
