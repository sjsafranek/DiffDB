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
	insertText string
	key        string
	snapshots  bool
	current    bool
	diffDb     DiffDb
)

func main() {
	cwd, _ := os.Getwd()
	databaseFile := path.Join(cwd, "data.db")
	flag.StringVar(&RuntimeArgs.DatabaseLocation, "db", databaseFile, "location of database file")
	flag.StringVar(&insertText, "insert", "", "insert text")
	flag.Int64Var(&timestamp, "ts", -1, "ts")
	flag.StringVar(&key, "key", "", "key name")
	flag.BoolVar(&snapshots, "s", false, "list snapshots")
	flag.BoolVar(&current, "c", false, "list current")
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
	log.Printf("%s\n", ddata)

	// Print older value
	if -1 < timestamp {
		oldValue := ddata.GetPrevious(timestamp)
		log.Println(timestamp, oldValue)
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
