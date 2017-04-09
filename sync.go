package fitbolt

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	fitbit "github.com/lenkaiser/go.fitbit"
	"github.com/sajal/fitbitclient"
)

const (
	baseURL = "https://api.fitbit.com/1/user/-/"
)

//Syncer syncs fitbit data to a DB
type Syncer struct {
	db        DB
	client    *http.Client
	lastSync  time.Time
	freshSync time.Time
	loc       *time.Location
}

//NewSyncer creates a new syncer
func NewSyncer(fname, fClient, fSecret string, db DB) (*Syncer, error) {
	//Initialise FitbitClient
	conf := &fitbitclient.Config{
		ClientID:     fClient,
		ClientSecret: fSecret,
		Scopes:       []string{"activity", "heartrate", "location", "nutrition", "profile", "settings", "sleep", "social", "weight"},
		CredFile:     fname,
	}
	client, err := fitbitclient.NewFitBitClient(conf)
	if err != nil {
		return nil, err
	}
	syncer := &Syncer{client: client, db: db}
	//Load timezone from fitbit API
	profile := &fitbit.Profile{}
	err = syncer.geturl("profile.json", profile)
	if err != nil {
		return nil, err
	}
	//log.Fatal(profile.User.Timezone)
	syncer.loc, err = time.LoadLocation(profile.User.Timezone)
	if err != nil {
		return nil, err
	}
	return syncer, nil
}

//Sync syncs from fitbit api starting from genisis
func (sync *Syncer) Sync(genesis string) (err error) {
	sync.lastSync, sync.freshSync, err = sync.getlastsync()
	if err != nil {
		return
	}
	current, err := time.ParseInLocation("2006-01-02", genesis, sync.loc)
	if err != nil {
		log.Fatal(err)
	}
	for !current.After(time.Now()) {
		isfinal, err := sync.isDayComplete(current)
		if err != nil {
			return err
		}
		if !isfinal {
			//Only sync if its not final...
			ds, err := sync.sync(current)
			if err != nil {
				return err
			}
			err = sync.db.StoreDayDetail(current, ds)
			if err != nil {
				return err
			}
		}
		current = current.Add(time.Hour * 24)
	}
	return
}

func (sync *Syncer) sync(dt time.Time) (*DayDetail, error) {
	log.Println("Syncing: ", dt)
	ds := &DayDetail{
		Date:    dt,
		Fetched: time.Now(),
	}
	ds.IsFinal = dt.Before(sync.lastSync.Add(time.Hour * -24))
	log.Println("Fetching activities...")
	act := &fitbit.Activities{}
	err := sync.geturl(fmt.Sprintf("activities/date/%s.json", dt.Format("2006-01-02")), act)
	if err != nil {
		return nil, err
	}
	ds.Activity = act

	log.Println("Fetching weight...")

	bodyData := &LogWeight{}
	err = sync.geturl(fmt.Sprintf("body/log/weight/date/%s.json", dt.Format("2006-01-02")), bodyData)
	if err != nil {
		return nil, err
	}
	ds.Weight = bodyData

	log.Println("Fetching sleep...")
	sl := &fitbit.Sleep{}
	err = sync.geturl(fmt.Sprintf("sleep/date/%s.json", dt.Format("2006-01-02")), sl)
	if err != nil {
		return nil, err
	}
	ds.Sleep = sl

	log.Println("Fetching steps...")
	stepintra := &StepIntra{}
	err = sync.geturl(fmt.Sprintf("activities/steps/date/%s/1d/1min.json", dt.Format("2006-01-02")), stepintra)
	if err != nil {
		return nil, err
	}
	ds.Steps = stepintra

	log.Println("Fetching heart...")
	h := &Heart{}
	err = sync.geturl(fmt.Sprintf("activities/heart/date/%s/1d/1sec.json", dt.Format("2006-01-02")), h)
	if err != nil {
		return nil, err
	}
	ds.HeartRate = h
	return ds, nil
}

func (sync *Syncer) isDayComplete(ts time.Time) (bool, error) {
	ds, err := sync.db.GetDayDetail(ts)
	//log.Println(ds)
	if err != nil {
		return false, err
	}
	if ds == nil {
		return false, nil
	}
	//If there has been no tracker activity since we last fetched this day, no point in continuing...
	//Check if freshSync is older than ds.Fetched
	if ds.Fetched.After(sync.freshSync) {
		//log.Println(ds.Fetched, freshSync)
		return true, nil
	}
	return ds.IsFinal, nil
}

//Helper function to load result of get into arbitary structure
func (sync *Syncer) geturl(url string, val interface{}) error {
	resp, err := sync.client.Get(baseURL + url)
	if err != nil {
		log.Println(resp.Header)
		return err
	}
	if resp.StatusCode != 200 {
		log.Println(resp.Header)
		return errors.New(resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(val)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (sync *Syncer) getlastsync() (ltime, ftime time.Time, err error) {
	devs := fitbit.GetDevices{}
	err = sync.geturl("devices.json", &devs)
	if err != nil {
		return
	}
	log.Println(devs)
	var t time.Time
	for _, dev := range devs {
		log.Println(dev)
		log.Println(dev.LastSyncTime)
		t, err = time.ParseInLocation("2006-01-02T15:04:05", dev.LastSyncTime, sync.loc)
		if err != nil {
			return
		}
		if ltime.IsZero() || ltime.After(t) {
			ltime = t
		}
		if ftime.IsZero() || ftime.Before(t) {
			ftime = t
		}
	}
	return
}
