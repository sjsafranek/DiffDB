package main

import (
	"flag"
	//"fmt"
	"log"
	"os"
	"path"
)

// RuntimeArgs contains all runtime
// arguments available
var RuntimeArgs struct {
	DatabaseLocation string
	Debug            bool
}

var (
	timestamp  int64
	index      int
	insertText string
	key        string
	snapshots  bool
	current    bool
	about      bool
	diffDb     DiffDb
)

func main() {
	cwd, _ := os.Getwd()
	databaseFile := path.Join(cwd, "data.db")
	flag.StringVar(&RuntimeArgs.DatabaseLocation, "db", databaseFile, "location of database file")
	flag.StringVar(&insertText, "insert", "", "insert text")
	flag.Int64Var(&timestamp, "ts", -1, "timestamp")
	flag.IntVar(&index, "i", -1, "index")
	flag.StringVar(&key, "key", "", "key name")
	flag.BoolVar(&snapshots, "s", false, "list snapshots")
	flag.BoolVar(&current, "c", false, "list current")
	flag.BoolVar(&about, "a", false, "about")
	flag.Parse()

	diffDb = NewDiffDb(RuntimeArgs.DatabaseLocation)

	if "" == key {
		log.Fatal("No key to lookup")
	}

	// Load key
	ddata, err := diffDb.Load(key)
	if nil != err {
		log.Fatal(err)
	}

	if about {
		log.Printf("%s\n", ddata)
	}

	// Print older value
	if -1 < timestamp {
		oldValue := ddata.GetPreviousByTimestamp(timestamp)
		log.Println(oldValue)
	}

	if -1 < index {
		oldValue := ddata.GetPreviousByIndex(index)
		log.Println(oldValue)
	}

	// Insert new value
	if "" != insertText {
		ddata.Update(insertText)
		diffDb.Save(ddata)
		log.Printf("%s\n", ddata)
	}

	// Print current value
	if current {
		log.Println(ddata.GetCurrent())
	}

	// List snapshots
	if snapshots {
		log.Println(ddata.GetSnapshots())
	}

}
