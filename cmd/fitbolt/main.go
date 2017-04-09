package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	cli "gopkg.in/urfave/cli.v1"

	"github.com/sajal/fitbolt"
)

func sync(dbpath, fname, fclient, fsecret, genesis string) error {
	if dbpath == "" {
		return fmt.Errorf("dbpath is blank")
	}
	if fname == "" {
		return fmt.Errorf("creds is blank")
	}
	if genesis == "" {
		return fmt.Errorf("genesis is blank")
	}
	if fclient == "" {
		return fmt.Errorf("FITBIT_CLIENT must be set")
	}
	if fsecret == "" {
		return fmt.Errorf("FITBIT_SECRET must be set")
	}
	db, err := fitbolt.NewBoltDB(dbpath)
	if err != nil {
		return err
	}
	defer db.Close()
	syncer, err := fitbolt.NewSyncer(fname, fclient, fsecret, db)
	if err != nil {
		return err
	}
	return syncer.Sync(genesis)
}

func main() {
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

	app := cli.NewApp()
	app.Name = "fitbolt"
	app.Usage = "sync and query fitbit data locally"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "dbpath",
			Value: dbpath,
			Usage: "path to bolt database(gets created if not exists)",
		},
		cli.StringFlag{
			Name:  "creds",
			Value: fname,
			Usage: "path to bolt database(gets created if not exists)",
		},
		cli.StringFlag{
			Name:  "genesis",
			Value: "2016-01-01",
			Usage: "The date to start fetching data from",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "sync",
			Usage: "Sync with fitbit, needs FITBIT_CLIENT and FITBIT_SECRET to be set",
			Action: func(c *cli.Context) error {
				return sync(c.GlobalString("dbpath"), c.GlobalString("creds"), fClient, fSecret, c.GlobalString("genesis"))
			},
		},
		{
			Name:  "steps",
			Usage: "list steps by day",
			Action: func(c *cli.Context) error {
				return fmt.Errorf("not implimented yet")
			},
		},
	}

	app.Run(os.Args)

	/*
		defer db.Close()
		syncer, err := fitbolt.NewSyncer(*fitbitCred, fClient, fSecret, db)
		if err != nil {
			log.Fatal(err)
		}
		err = syncer.Sync(*genisis)
		if err != nil {
			log.Fatal(err)
		}*/
}
