package budget

import (
	"time"
)

// add is a basic implementation of vector addition to add two vectors
// (or arrays) of length 12 together.
func add(x [12]float64, y [12]float64) [12]float64 {
	var t [12]float64

	for i := 0; i < 12; i++ {
		t[i] = x[i] + y[i]
	}
	return t
}


// subtract is a basic implementation of vector subtraction to subtract two
// vectors (or arrays) of length 12 together.
func subtract(x [12]float64, y [12]float64) [12]float64 {
	var t [12]float64

	for i := 0; i < 12; i++ {
		t[i] = x[i] - y[i]
	}
	return t
}

// yearArray converts the event from a time period it and array of monthly amounts.
// TODO: improve to incorporate specific dates, such as specific date (day) transactions or closest (last friday of month)
func yearArray(year int, e Event) [12]float64 {
	// declare variable
	var m [12]float64

	// do computation
	// loop through months
	for i := 1; i <= 12; i++ {
		tmi := time.Month(i)
		firstOfMonth := time.Date(year, tmi, 1, 0, 0, 0, 0, time.UTC)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

		afterStart := e.StartDate.Before(lastOfMonth) || e.StartDate.Equal(lastOfMonth)
		beforeEnd := e.EndDate.After(firstOfMonth) || e.EndDate.Equal(firstOfMonth)

		if afterStart && beforeEnd {
			m[i-1] = e.Amount
		}
	}

	return m
}
