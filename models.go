package fitbolt

import (
	"time"

	fitbit "github.com/lenkaiser/go.fitbit"
)

//LogWeight stores Weight details. Upstream is broken...
type LogWeight struct {
	WeightLog []*fitbit.Weight `json:"weight"`
}

//StepDataset Stores intra-day steps activity dataset
type StepDataset struct {
	Time  string `json:"time"`
	Value uint64 `json:"value"`
}

//ActivitiesSteps Stores intra-day steps activity
type ActivitiesSteps struct {
	Dataset []*StepDataset `json:"dataset"`
}

//StepIntra Stores intra-day steps activity
type StepIntra struct {
	Steps *ActivitiesSteps `json:"activities-steps-intraday"`
}

//HeartRateZone stores heart rate
type HeartRateZone struct {
	CaloriesOut float64 `json:"caloriesOut"`
	Max         int     `json:"max"`
	Min         int     `json:"min"`
	Name        string  `json:"name"`
}

//HeartRateZoneValue Stores heartrate values
type HeartRateZoneValue struct {
	HeartRateZone []*HeartRateZone `json:"heartRateZones"`
	Resting       int              `json:"restingHeartRate"`
}

//HeartActivityZones Stores heartrate zones
type HeartActivityZones struct {
	HeartRateZones *HeartRateZoneValue `json:"value"`
}

//Heart Stores details about the heart
type Heart struct {
	ActivityZones []*HeartActivityZones `json:"activities-heart"`
	IntraDay      *ActivitiesSteps      `json:"activities-heart-intraday"`
}

//DayDetail is all fitbit data we track for a single day
type DayDetail struct {
	Date      time.Time          `json:"date"` //What date is this report about
	Fetched   time.Time          `json:"fetched"`
	Weight    *LogWeight         `json:"weight"`
	Activity  *fitbit.Activities `json:"activity"`
	Sleep     *fitbit.Sleep      `json:"sleep"`
	IsFinal   bool               `json:"isfinal"` //Is this report finalized? to not spam fitbit api again for it...
	Steps     *StepIntra         `json:"steps"`
	HeartRate *Heart             `json:"heart"`
}
