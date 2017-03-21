package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
)

import (
	"./skeleton_db"
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

func errorHandler(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func incorrectUsageError() {
	err := fmt.Errorf("Incorrect usage! Nonsensical argument!")
	errorHandler(err)
}

func successHandler(msg string) {
	fmt.Println(msg)
	os.Exit(0)
}

func usage() {
	fmt.Printf("%s %s\n\n", NAME, skeleton_db.VERSION)
	fmt.Printf("Usage:\n\t%s [options...] action key [action_args...]\n\n", BINARY)
	fmt.Println(" * action:\tThe action to preform. Supported action(s): GET, SET, DEL")
	fmt.Println(" * action_args:\tVariadic arguments provided to the requested action. Different actions require different arguments")
	fmt.Println("\n")
}

func main() {
	cwd, _ := os.Getwd()
	databaseFile := path.Join(cwd, "data.db")

	// handle command line arguements
	flag.Usage = usage
	flag.StringVar(&RuntimeArgs.DatabaseLocation, "db", databaseFile, "location of database file")
	flag.BoolVar(&RuntimeArgs.PrintVersion, "v", false, "version")
	flag.BoolVar(&RuntimeArgs.Verbose, "verbose", false, "verbose")
	flag.Parse()

	// list version
	if RuntimeArgs.PrintVersion {
		fmt.Println("SkeletonDb ", skeleton_db.VERSION)
		os.Exit(0)
	}

	// create database object
	diffDb = skeleton_db.NewDiffDb(RuntimeArgs.DatabaseLocation)

	// get args
	args := flag.Args()
	if 2 > len(args) {
		incorrectUsageError()
	}

	// command line args
	action := args[0]
	key := args[1]

	switch action {
	// get value of key
	case "GET":
		ddata, err := diffDb.Load(key)
		if nil != err {
			errorHandler(err)
		}

		if 2 == len(args) {
			enc, _ := ddata.Encode()
			msg := fmt.Sprintf("%s", enc)
			successHandler(msg)
		}

		if "VALUE" == args[2] {
			successHandler(ddata.GetCurrent())
		}

		if "SNAPSHOTS" == args[2] {
			msg := fmt.Sprintf("%v", ddata.GetSnapshots())
			successHandler(msg)
		}

		if 4 == len(args) {

			num, err := strconv.ParseInt(args[3], 10, 64)
			if nil != err {
				errorHandler(err)
			}

			if "TIMESTAMP" == args[2] {
				val, err := ddata.GetPreviousByTimestamp(num)
				if nil != err {
					errorHandler(err)
				}
				successHandler(val)
			}

			if "INDEX" == args[2] {
				val, err := ddata.GetPreviousByIndex(int(num))
				if nil != err {
					errorHandler(err)
				}
				successHandler(val)
			}

		}

		incorrectUsageError()

	// set new value for key
	case "SET":

		// check for data to set as new value
		if 3 > len(args) {
			incorrectUsageError()
		}

		// load key
		ddata, err := diffDb.Load(key)
		if nil != err {
			if err.Error() == "Not found" {
				// create new diffstore if key not found in database
				ddata = skeleton_db.NewDiffStore(key)
			} else {
				errorHandler(err)
			}
		}

		// update diffstore
		ddata.Update(args[2])

		// save to database
		diffDb.Save(ddata)

		// print result
		successHandler(ddata.GetCurrent())

	// delete key
	case "DEL":

		ddata, err := diffDb.Load(key)
		if nil != err {
			errorHandler(err)
		}

		err = diffDb.Remove(ddata)
		if nil != err {
			errorHandler(err)
		}

	default:
		err := fmt.Errorf("Unsupported action %s, cannot process.", args[0])
		errorHandler(err)
	}

}
