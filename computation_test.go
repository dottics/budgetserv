package budget

import (
	"log"
	"testing"
	"time"
)

func MustParse(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		log.Panicln(err.Error())
	}
	return t
}

func Test_yearArray(t *testing.T) {
	tt := []struct {
		name string
		year int
		event Event
		z [12]float64
	}{
		{
			name: "blank event",
			year: 0,
			event: Event{},
			z: [12]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "year event mismatch",
			year: 2020,
			event: Event{
				Amount: 5.05,
				StartDate: MustParse("2021-04-01"),
				EndDate: MustParse("2021-04-01"),
			},
			z: [12]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "year event match single month",
			year: 2020,
			event: Event{
				Amount: 5.05,
				StartDate: MustParse("2020-04-01"),
				EndDate: MustParse("2020-04-01"),
			},
			z: [12]float64{0, 0, 0, 5.05, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "year event match multiple months",
			year: 2020,
			event: Event{
				Amount: 5.05,
				StartDate: MustParse("2020-04-01"),
				EndDate: MustParse("2020-06-01"),
			},
			z: [12]float64{0, 0, 0, 5.05, 5.05, 5.05, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "year event start",
			year: 2020,
			event: Event{
				Amount: 5.05,
				StartDate: MustParse("2020-04-01"),
				EndDate: MustParse("2021-04-01"),
			},
			z: [12]float64{0, 0, 0, 5.05, 5.05, 5.05, 5.05, 5.05, 5.05, 5.05, 5.05, 5.05},
		},
		{
			name: "year event end",
			year: 2021,
			event: Event{
				Amount: 5.05,
				StartDate: MustParse("2020-04-01"),
				EndDate: MustParse("2021-04-01"),
			},
			z: [12]float64{5.05, 5.05, 5.05, 5.05, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			z := yearArray(tc.year, tc.event)
			if z != tc.z {
				t.Errorf("expected '%v' got '%v'", tc.z, z)
			}
		})
	}
}

func Test_add(t *testing.T) {
	tt := []struct {
		name string
		x [12]float64
		y [12]float64
		z [12]float64
	}{
		{
			name: "zero value",
			z: [12]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "x and y",
			x: [12]float64{1, 0, 0, 0, 0, 0, 0, 0, 0, -5, -2, 2},
			y: [12]float64{0, 0, 12.34, 0, 0, 0, 0, 0, 0, 2, 2, 3},
			z: [12]float64{1, 0, 12.34, 0, 0, 0, 0, 0, 0, -3, 0, 5},
		},
	}

	for  _, tc := range  tt  {
		t.Run(tc.name, func(t *testing.T) {
			z := add(tc.x, tc.y)
			if z != tc.z {
				t.Errorf("expected total '%v' got '%v'", tc.z, z)
			}
		})
	}
}

func Test_subtract(t *testing.T) {
	tt := []struct {
		name string
		x [12]float64
		y [12]float64
		z [12]float64
	}{
		{
			name: "zero value",
			z: [12]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "x and y",
			x: [12]float64{1, 0, 0, 0, 0, 0, 0, 0, 0, -5, -2, 2},
			y: [12]float64{0, 0, 12.34, 0, 0, 0, 0, 0, 0, 2, -2, 3},
			z: [12]float64{1, 0, -12.34, 0, 0, 0, 0, 0, 0, -7, 0, -1},
		},
	}

	for  _, tc := range  tt  {
		t.Run(tc.name, func(t *testing.T) {
			z := subtract(tc.x, tc.y)
			if z != tc.z {
				t.Errorf("expected total '%v' got '%v'", tc.z, z)
			}
		})
	}
}
