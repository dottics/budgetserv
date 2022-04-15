package budget

import "testing"

func TestItem_MonthlyTotal(t *testing.T) {
	tt := []struct {
		name string
		item Item
		total [12]float64
	}{
		{
			name: "no event",
			item: Item{Events: Events{}},
			total: [12]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "single event single month",
			item: Item{Events: Events{
				Event{
					StartDate: MustParse("2021-03-12"),
					EndDate: MustParse("2021-03-12"),
					Debit: true,
					Amount: 5,
				},
			}},
			total: [12]float64{0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "single event multiple months",
			item: Item{Events: Events{
				Event{
					StartDate: MustParse("2021-03-12"),
					EndDate: MustParse("2021-06-12"),
					Debit: true,
					Amount: 5,
				},
			}},
			total: [12]float64{0, 0, 5, 5, 5, 5, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "single event multiple months start previous year",
			item: Item{Events: Events{
				Event{
					StartDate: MustParse("2020-03-12"),
					EndDate: MustParse("2021-03-12"),
					Debit: true,
					Amount: 5,
				},
			}},
			total: [12]float64{5, 5, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "single event multiple months end next year",
			item: Item{Events: Events{
				Event{
					StartDate: MustParse("2021-11-12"),
					EndDate: MustParse("2022-03-12"),
					Debit: true,
					Amount: 5,
				},
			}},
			total: [12]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 5},
		},
		{
			name: "multiple events",
			item: Item{Events: Events{
				Event{
					StartDate: MustParse("2021-03-12"),
					EndDate: MustParse("2021-03-20"),
					Debit: true,
					Amount: 5,
				},
				Event{
					StartDate: MustParse("2021-03-12"),
					EndDate: MustParse("2021-04-20"),
					Debit: true,
					Amount: 3,
				},
			}},
			total: [12]float64{0, 0, 8, 3, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "multiple events debit credit",
			item: Item{Events: Events{
				Event{
					StartDate: MustParse("2021-03-12"),
					EndDate: MustParse("2021-03-20"),
					Debit: true,
					Amount: 5,
				},
				Event{
					StartDate: MustParse("2021-03-12"),
					EndDate: MustParse("2021-04-20"),
					Credit: true,
					Amount: 3,
				},
			}},
			total: [12]float64{0, 0, 2, -3, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			x := tc.item.MonthlyTotal(2021)
			if x != tc.total {
				t.Errorf("expected '%v' got '%v'", tc.total, x)
			}
		})
	}
}
