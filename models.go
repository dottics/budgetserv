package budget

import (
	"github.com/google/uuid"
	"time"
)

type Budget struct {
	UUID             uuid.UUID `json:"uuid"`
	UserUUID         uuid.UUID `json:"user_uuid"`
	OrganisationUUID uuid.UUID `json:"organisation_uuid"`
	Name             string    `json:"name"`
	Active           bool      `json:"active"`
}

type Budgets []Budget

type Group struct {
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	SubGroups []Group   `json:"sub_groups"`
}

type Groups []Group

type Item struct {
	UUID   uuid.UUID `json:"uuid"`
	Name   string    `json:"name"`
	Active bool      `json:"active"`
	Events Events `json:"-"`
}

// MonthlyTotal calculates to total monthly effect of each event of the item.
func (i *Item) MonthlyTotal(year int) [12]float64 {
	var a [12]float64

	for _, ev := range i.Events {
		evi := yearArray(year, ev)
		if ev.Debit {
			a = add(a, evi)
		}
		if ev.Credit {
			a = subtract(a, evi)
		}
	}
	return a
}

type Items []Item

type Event struct {
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	Debit     bool      `json:"debit"`
	Credit    bool      `json:"credit"`
	Amount    float64   `json:"amount"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Active    bool      `json:"active"`
}

type Events []Event
