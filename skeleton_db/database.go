package skeleton_db

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"os"
	"time"
)

import (
	"github.com/boltdb/bolt"
)

// Database strust for application.
type Database struct {
	File string
}

// Create to bolt database. Returns open database connection.
// @returns *bolt.DB
func (self *Database) createDb() {
	conn, err := bolt.Open(self.File, 0644, nil)
	if err != nil {
		conn.Close()
		panic(err)
	}
	conn.Close()
}

// Connect to bolt database. Returns open database connection.
// @returns *bolt.DB
func (self *Database) Connect() *bolt.DB {
	// Check if file exists
	_, err := os.Stat(self.File)
	if err != nil {
		panic("Database not found!")
	}

	// Open database connection
	config := &bolt.Options{Timeout: 30 * time.Second}
	conn, err := bolt.Open(self.File, 0644, config)
	if err != nil {
		conn.Close()
		panic(err)
	}
	return conn
}

// Init creates bolt database if existing one not found.
// Creates layers and apikey tables. Starts database caching for layers
// @returns Error
func (self *Database) Init() {
	self.createDb()
}

// CreateTable creates bucket to store data
// @param table
// @returns Error
func (self *Database) CreateTable(conn *bolt.DB, table string) error {
	err := conn.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(table))
		return err
	})
	return err
}

func (self *Database) Insert(table string, key string, value []byte) error {
	// connect to database and write to table
	conn := self.Connect()
	defer conn.Close()
	err := conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(table))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", table)
		}
		err := bucket.Put([]byte(key), self.compressByte(value))
		return err
	})
	return err
}

func (self *Database) Select(table string, key string) ([]byte, error) {
	conn := self.Connect()
	defer conn.Close()
	val := []byte{}
	err := conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(table))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", table)
		}
		val = self.decompressByte(bucket.Get([]byte(key)))
		return nil
	})
	return val, err
}

func (self *Database) SelectAll(table string) ([]string, error) {
	conn := self.Connect()
	defer conn.Close()
	data := []string{}
	err := conn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(table))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", table)
		}
		bucket.ForEach(func(key, _ []byte) error {
			data = append(data, string(key))
			return nil
		})
		return nil
	})
	return data, err
}

func (self *Database) Remove(table string, name string) error {
	conn := self.Connect()
	defer conn.Close()

	err := conn.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(table))
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

// Methods: Compression
// Source: https://github.com/schollz/gofind/blob/master/utils.go#L146-L169
//         https://github.com/schollz/gofind/blob/master/fingerprint.go#L43-L54
// Description:
//		Compress and Decompress bytes
func (self *Database) compressByte(src []byte) []byte {
	compressedData := new(bytes.Buffer)
	self.compress(src, compressedData, 9)
	return compressedData.Bytes()
}

func (self *Database) decompressByte(src []byte) []byte {
	compressedData := bytes.NewBuffer(src)
	deCompressedData := new(bytes.Buffer)
	self.decompress(compressedData, deCompressedData)
	return deCompressedData.Bytes()
}

func (self *Database) compress(src []byte, dest io.Writer, level int) {
	compressor, _ := flate.NewWriter(dest, level)
	compressor.Write(src)
	compressor.Close()
}

func (self *Database) decompress(src io.Reader, dest io.Writer) {
	decompressor := flate.NewReader(src)
	io.Copy(dest, decompressor)
	decompressor.Close()
}
