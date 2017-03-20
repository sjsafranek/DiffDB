package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
)

import (
	"skeleton_db"
)

const (
	NAME   = "SkeletonDB Client"
	BINARY = "skeleton_cli"
)

// RuntimeArgs contains all runtime
// arguments available
var RuntimeArgs struct {
	DatabaseLocation string
	Debug            bool
}

var (
	// timestamp    int64
	// index        int
	// insertText   string
	// key          string
	// snapshots    bool
	// current      bool
	// about        bool
	print_version bool
	diffDb        skeleton_db.DiffDb
)

func main() {
	cwd, _ := os.Getwd()
	databaseFile := path.Join(cwd, "data.db")

	flag.Usage = func() {
		fmt.Printf("%s %s\n\n", NAME, skeleton_db.VERSION)
		fmt.Printf("Usage:\n\t%s [options...] action key [action_args...]\n\n", BINARY)
	}

	flag.StringVar(&RuntimeArgs.DatabaseLocation, "db", databaseFile, "location of database file")
	flag.BoolVar(&print_version, "v", false, "version")
	flag.Parse()

	// list version
	if print_version {
		fmt.Println("SkeletonDb", skeleton_db.VERSION)
		os.Exit(0)
	}

	// create database object
	diffDb = skeleton_db.NewDiffDb(RuntimeArgs.DatabaseLocation)

	// get args
	args := flag.Args()
	if 0 == len(args) {
		fmt.Printf("An action is required. Usage %s [options...] action key [action_args...]", BINARY)
		os.Exit(1)
	}
	if 2 > len(args) {
		fmt.Printf("A key is required. Usage %s [options...] action key [action_args...]", BINARY)
		os.Exit(1)
	}

	// command line args
	action := args[0]
	key := args[1]

	switch action {
	// get value of key
	case "GET":
		ddata, err := diffDb.Load(key)
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}

		if 2 == len(args) {
			val, _ := ddata.Encode()
			fmt.Printf("%s\n", val)
			os.Exit(0)
		}

		if "VALUE" == args[2] {
			fmt.Println(ddata.GetCurrent())
			os.Exit(0)
		}

		if "SNAPSHOTS" == args[2] {
			fmt.Println(ddata.GetSnapshots())
			os.Exit(0)
		}

		if 4 == len(args) {

			num, err := strconv.ParseInt(args[3], 10, 64)
			if nil != err {
				fmt.Println(err)
				os.Exit(1)
			}

			if "TIMESTAMP" == args[2] {
				oldValue := ddata.GetPreviousByTimestamp(num)
				fmt.Println(oldValue)
				os.Exit(0)
			}

			if "INDEX" == args[2] {
				oldValue := ddata.GetPreviousByIndex(int(num))
				fmt.Println(oldValue)
				os.Exit(0)
			}

		}

		// alert error
		fmt.Println("Incorrect usage! Nonsensical argument!")
		os.Exit(1)

	// set new value for key
	case "SET":

		// check for data to set as new value
		if 3 > len(args) {
			fmt.Printf("An value is required. Usage %s [options...] action key [action_args...]", BINARY)
			os.Exit(1)
		}

		// load key
		ddata, err := diffDb.Load(key)
		if nil != err {
			if err.Error() == "Not found" {
				// create new diffstore if key not found in database
				ddata = skeleton_db.NewDiffStore(key)
			} else {
				log.Fatal(err)
			}
		}

		// update diffstore
		ddata.Update(args[2])

		// save to database
		diffDb.Save(ddata)

		// print result
		fmt.Printf("%s\n", ddata.GetCurrent())

	default:
		msg := fmt.Sprintf("Unsupported action %s, cannot process.", args[0])
		log.Fatal(msg)
	}

}
