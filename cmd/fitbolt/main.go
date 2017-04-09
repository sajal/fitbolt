package main

import (
	"flag"
	"log"
	"os"
	"os/user"

	"github.com/sajal/fitbolt"
)

var (
	db      *fitbolt.BoltDB
	syncer  *fitbolt.Syncer
	genisis *string
)

func init() {
	var err error
	usr, err := user.Current()
	if err != nil {
		log.Fatal("user", err)
	}
	//log.Println( usr.HomeDir )
	fname := usr.HomeDir + "/.fitsyncgo"
	dbpath := usr.HomeDir + "/fitsyncgo.db"
	fClient := os.Getenv("FITBIT_CLIENT")
	fSecret := os.Getenv("FITBIT_SECRET")

	//Provide flags with defaults
	fitbitCred := flag.String("creds", fname, "Path to store/load fitbit tokens from")
	dbPath := flag.String("db", dbpath, "path to bolt database(gets created if not exists)")
	genisis = flag.String("genisis", "2016-01-01", "The date to start fetching data from")
	flag.Parse()

	//Check env variables
	if fClient == "" {
		log.Fatal("FITBIT_CLIENT must be set")
	}
	if fSecret == "" {
		log.Fatal("FITBIT_SECRET must be set")
	}

	db, err = fitbolt.NewBoltDB(*dbPath)
	if err != nil {
		log.Fatal(err)
	}

	syncer, err = fitbolt.NewSyncer(*fitbitCred, fClient, fSecret, db)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	defer db.Close()

	err := syncer.Sync(*genisis)
	if err != nil {
		log.Fatal(err)
	}
}
