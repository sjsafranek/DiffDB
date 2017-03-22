package skeleton_db

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

import "github.com/boltdb/bolt"

// Creates and initializes DiffDb
func NewDiffDb(db_file string) DiffDb {
	var diffDb = DiffDb{Filename: db_file, Table: "DiffData"}
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
	if "" == self.Table {
		self.Table = "DiffData"
	}
	err := self.CreateTable(self.Table)
	if nil != err {
		log.Fatal(err)
	}
}

func (self DiffDb) getBucketName() []byte {
	return []byte(self.Table)
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

func (self *DiffDb) Load(name string) ([]byte, error) {
	self.db = self.Open()
	defer self.db.Close()

	name = strings.ToLower(name)

	var data []byte

	err := self.db.View(func(tx *bolt.Tx) error {
		//var err error
		bucket := tx.Bucket(self.getBucketName())
		if bucket == nil {
			panic(fmt.Errorf("Bucket does not exist"))
		}

		k := []byte(name)
		val := bucket.Get(k)
		if val == nil {
			return fmt.Errorf("Not found")
		}

		// decompress data
		data = self.decompressByte(val)

		return nil
	})
	return data, err
}

func (self *DiffDb) Save(name string, data []byte) error {
	self.db = self.Open()
	defer self.db.Close()

	err := self.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(self.getBucketName())
		if bucket == nil {
			panic(fmt.Errorf("Bucket does not exist"))
		}

		// compress data
		cmp := self.compressByte(data)

		// store in database
		err := bucket.Put([]byte(name), cmp)
		if err != nil {
			return fmt.Errorf("could not add to bucket: %s", err)
		}
		return err
	})
	return err
}

func (self *DiffDb) Remove(name string) error {
	self.db = self.Open()
	defer self.db.Close()

	err := self.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(self.getBucketName())
		if bucket == nil {
			panic(fmt.Errorf("Bucket does not exist"))
		}

		err := bucket.Delete([]byte(name))
		if err != nil {
			return fmt.Errorf("could not delete key: %s", err)
		}
		return err
	})
	return err
}

func (self *DiffDb) SelectAll() ([]string, error) {
	self.db = self.Open()
	defer self.db.Close()
	data := []string{}
	err := self.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(self.getBucketName()))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", self.getBucketName())
		}
		bucket.ForEach(func(key, _ []byte) error {
			data = append(data, string(key))
			return nil
		})
		return nil
	})
	return data, err
}

// Methods: Compression
// Source: https://github.com/schollz/gofind/blob/master/utils.go#L146-L169
//         https://github.com/schollz/gofind/blob/master/fingerprint.go#L43-L54
// Description:
//		Compress and Decompress bytes
func (self *DiffDb) compressByte(src []byte) []byte {
	compressedData := new(bytes.Buffer)
	self.compress(src, compressedData, 9)
	return compressedData.Bytes()
}

func (self *DiffDb) decompressByte(src []byte) []byte {
	compressedData := bytes.NewBuffer(src)
	deCompressedData := new(bytes.Buffer)
	self.decompress(compressedData, deCompressedData)
	return deCompressedData.Bytes()
}

func (self *DiffDb) compress(src []byte, dest io.Writer, level int) {
	compressor, _ := flate.NewWriter(dest, level)
	compressor.Write(src)
	compressor.Close()
}

func (self *DiffDb) decompress(src io.Reader, dest io.Writer) {
	decompressor := flate.NewReader(src)
	io.Copy(dest, decompressor)
	decompressor.Close()
}
