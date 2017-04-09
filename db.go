package fitbolt

import "time"

//DB interface to allow pluggable database backends in future... some day...
type DB interface {
	GetDayDetail(ts time.Time) (*DayDetail, error)
	StoreDayDetail(ts time.Time, ds *DayDetail) error
	Close()
}
