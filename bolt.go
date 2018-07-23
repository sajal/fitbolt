package fitbolt

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang/snappy"
)

//BoltDB stores/retrieves our data from boltdb
type BoltDB struct {
	db *bolt.DB
}

//NewBoltDB returns a new instance of BoltDB
func NewBoltDB(dbpath string) (*BoltDB, error) {
	db, err := bolt.Open(dbpath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	//Init bucket(s) if not exists...
	err = db.Update(func(tx *bolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists([]byte("DayDetail"))
		if e != nil {
			return fmt.Errorf("create bucket: %s", e)
		}
		return e
	})
	if err != nil {
		db.Close() //Cleanup
		return nil, err
	}
	return &BoltDB{db}, nil
}

//GetDayDetail retrieves data for a given day represented by ts.
func (db *BoltDB) GetDayDetail(ts time.Time) (ds *DayDetail, err error) {
	err = db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DayDetail"))
		t := []byte(ts.Format("2006-01-02"))
		v := b.Get(t)
		//log.Println(v)
		if v == nil {
			return nil
		}
		ds = &DayDetail{}
		dec := gob.NewDecoder(snappy.NewReader(bytes.NewBuffer(v)))
		err = dec.Decode(ds)
		return err
	})
	return
}

//StoreDayDetail stores a day's data into the db
func (db *BoltDB) StoreDayDetail(ts time.Time, ds *DayDetail) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("DayDetail"))
		var network bytes.Buffer
		snp := snappy.NewWriter(&network)
		enc := gob.NewEncoder(snp)
		err := enc.Encode(ds)
		if err != nil {
			return fmt.Errorf("encode gob DayDetail: %s", err)
		}
		t := []byte(ts.Format("2006-01-02"))
		bucket.Put(t, network.Bytes())
		return nil
	})
}

//Close the bolt db
func (db *BoltDB) Close() {
	db.db.Close()
}
