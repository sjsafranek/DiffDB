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

//var VersionNum string

var inputText string

func main() {
	cwd, _ := os.Getwd()
	databaseFile := path.Join(cwd, "data.db")
	flag.StringVar(&RuntimeArgs.DatabaseLocation, "db", databaseFile, "location of database file")
	flag.StringVar(&inputText, "text", "", "inputText")

	flag.Parse()

	// create programdata bucket
	var diffDb = DiffDb{Filename: RuntimeArgs.DatabaseLocation}
	diffDb.Init()

	//var ddata DiffData
	ddata, err := diffDb.load("test_file")
	if nil != err {
		log.Fatal(err)
	}

	log.Printf("%s\n", ddata)

	ddata.Update("TESTING")
	diffDb.save(ddata)
	log.Printf("%s\n", ddata)

	ddata.Update(`{"method":"test"}`)
	diffDb.save(ddata)
	log.Printf("%s\n", ddata)

	str0, err := ddata.rebuildTextsToDiffN(0)
	if nil != err {
		log.Fatal(err)
	}
	log.Printf("%s\n", str0)

	log.Println(ddata.GetImportantVersions())

	//log.Println(ddata.GetCurrent())

}
