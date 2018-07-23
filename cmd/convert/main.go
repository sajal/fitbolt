package main

//convert is tiny utility to change the database key from binary to string dates.

import (
	"bytes"
	"encoding/gob"
	"log"
	"os/user"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang/snappy"
	"github.com/sajal/fitbolt"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("user", err)
	}
	src, err := bolt.Open(usr.HomeDir+"/fitsyncgo.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	dst, err := fitbolt.NewBoltDB(usr.HomeDir + "/fitsyncgo_new.db")
	if err != nil {
		log.Fatal(err)
	}

	err = src.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DayDetail"))
		log.Println(b)
		return b.ForEach(func(k []byte, v []byte) error {
			var ts time.Time
			err = ts.GobDecode(k)
			if err != nil {
				log.Println(k)
				return err
			}
			log.Println(ts)
			ds := &fitbolt.DayDetail{}
			dec := gob.NewDecoder(snappy.NewReader(bytes.NewBuffer(v)))
			err = dec.Decode(ds)
			if err != nil {
				return err
			}
			return dst.StoreDayDetail(ts, ds)

		})
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dst)
}
