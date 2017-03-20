package main

import (
	"flag"
	"fmt"
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
	PrintVersion     bool
	Verbose          bool
}

var (
	diffDb skeleton_db.DiffDb
)

// func errorHandler(err) {

// 	os.Exit(1)
// }

func main() {
	cwd, _ := os.Getwd()
	databaseFile := path.Join(cwd, "data.db")

	flag.Usage = func() {
		fmt.Printf("%s %s\n\n", NAME, skeleton_db.VERSION)
		fmt.Printf("Usage:\n\t%s [options...] action key [action_args...]\n\n", BINARY)
		fmt.Println(" * action:\tThe action to preform. Supported action(s): GET, SET, DEL")
		fmt.Println(" * action_args:\tVariadic arguments provided to the requested action. Different actions require different arguments")
		fmt.Println("\n")
	}

	flag.StringVar(&RuntimeArgs.DatabaseLocation, "db", databaseFile, "location of database file")
	flag.BoolVar(&RuntimeArgs.PrintVersion, "v", false, "version")
	flag.BoolVar(&RuntimeArgs.Verbose, "verbose", false, "verbose")
	flag.Parse()

	// list version
	if RuntimeArgs.PrintVersion {
		fmt.Println("SkeletonDb", skeleton_db.VERSION)
		os.Exit(0)
	}

	// create database object
	diffDb = skeleton_db.NewDiffDb(RuntimeArgs.DatabaseLocation)

	// get args
	args := flag.Args()
	if 0 == len(args) {
		fmt.Printf("An action is required. Usage %s [options...] action key [action_args...]\n", BINARY)
		os.Exit(1)
	}
	if 2 > len(args) {
		fmt.Printf("A key is required. Usage %s [options...] action key [action_args...]\n", BINARY)
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
			fmt.Printf("An value is required. Usage %s [options...] action key [action_args...]\n", BINARY)
			os.Exit(1)
		}

		// load key
		ddata, err := diffDb.Load(key)
		if nil != err {
			if err.Error() == "Not found" {
				// create new diffstore if key not found in database
				ddata = skeleton_db.NewDiffStore(key)
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		// update diffstore
		ddata.Update(args[2])

		// save to database
		diffDb.Save(ddata)

		// print result
		fmt.Printf("%s\n", ddata.GetCurrent())

	// delete key
	case "DEL":

		ddata, err := diffDb.Load(key)
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}

		err = diffDb.Remove(ddata)
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unsupported action %s, cannot process.", args[0])
		os.Exit(1)
	}

}
