package fitbolt

//DaySteps returns total steps taken in a day
func (ds *DayDetail) DaySteps() (steps uint64) {
	for _, s := range ds.Steps.Steps.Dataset {
		steps += s.Value
	}
	return
}
