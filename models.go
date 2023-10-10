package budget

import (
	"github.com/google/uuid"
	"time"
)

type BudgetCreatePayload struct {
	EntityUUID  uuid.UUID `json:"entity_uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type BudgetUpdatePayload struct {
	UUID        uuid.UUID `json:"uuid"`
	EntityUUID  uuid.UUID `json:"entity_uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type Budget struct {
	UUID        uuid.UUID `json:"uuid"`
	EntityUUID  uuid.UUID `json:"entity_uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Groups      Groups    `json:"groups"`
}

type Budgets []Budget

type GroupCreatePayload struct {
	BudgetUUID  uuid.UUID `json:"budget_uuid,omitempty"`
	GroupUUID   uuid.UUID `json:"group_uuid,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type GroupUpdatePayload struct {
	UUID        uuid.UUID `json:"uuid"`
	BudgetUUID  uuid.UUID `json:"budget_uuid,omitempty"`
	GroupUUID   uuid.UUID `json:"group_uuid,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type Group struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SubGroups   []Group   `json:"sub_groups"`
	Items       []Item    `json:"items"`
}

type Groups []Group

type ItemCreatePayload struct {
	GroupUUID   uuid.UUID `json:"group_uuid"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
}
type ItemUpdatePayload struct {
	UUID        uuid.UUID `json:"uuid"`
	GroupUUID   uuid.UUID `json:"group_uuid"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
}

type Item struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Events      Events    `json:"events"`
}

// MonthlyTotal calculates to total monthly effect of each event of the item.
func (i Item) MonthlyTotal(year int) [12]float64 {
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

type EventCreate struct {
	ItemUUID    uuid.UUID `json:"item_uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Debit       bool      `json:"debit"`
	Credit      bool      `json:"credit"`
	Amount      float64   `json:"amount"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type EventUpdate = EventCreate

type Event struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Debit       bool      `json:"debit"`
	Credit      bool      `json:"credit"`
	Amount      float64   `json:"amount"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type Events []Event

/* *** Equality Functions *************************************************** */

// EqualBudgets returns true if the two budgets are equal.
func EqualBudgets(a, b Budgets) bool {
	if len(a) != len(b) {
		return false
	}
	for i, budget := range a {
		if !EqualBudget(budget, b[i]) {
			return false
		}
	}
	return true
}

// EqualBudget returns true if the two budgets are equal.
func EqualBudget(a, b Budget) bool {
	switch {
	case a.UUID != b.UUID:
		return false
	case a.EntityUUID != b.EntityUUID:
		return false
	case a.Name != b.Name:
		return false
	case a.Description != b.Description:
		return false
	case !EqualGroups(a.Groups, b.Groups):
		return false
	default:
		return true
	}
}

// EqualGroups returns true if the two groups are equal.
func EqualGroups(a, b Groups) bool {
	if len(a) != len(b) {
		return false
	}
	for i, group := range a {
		if !EqualGroup(group, b[i]) {
			return false
		}
	}
	return true
}

// EqualGroup returns true if the two groups are equal.
func EqualGroup(a, b Group) bool {
	switch {
	case a.UUID != b.UUID:
		return false
	case a.Name != b.Name:
		return false
	case a.Description != b.Description:
		return false
	case !EqualGroups(a.SubGroups, b.SubGroups):
		return false
	default:
		return true
	}
}

// EqualItems returns true if the two items are equal.
func EqualItems(a, b Items) bool {
	if len(a) != len(b) {
		return false
	}
	for i, item := range a {
		if !EqualItem(item, b[i]) {
			return false
		}
	}
	return true
}

// EqualItem returns true if the two items are equal.
func EqualItem(a, b Item) bool {
	if !EqualEvents(a.Events, b.Events) {
		return false
	}
	switch {
	case !EqualEvents(a.Events, b.Events):
		return false
	case a.UUID != b.UUID:
		return false
	case a.Name != b.Name:
		return false
	case a.Description != b.Description:
		return false
	default:
		return true
	}
}

// EqualEvents returns true if the two events are equal.
func EqualEvents(a, b Events) bool {
	if len(a) != len(b) {
		return false
	}
	for i, event := range a {
		if !EqualEvent(event, b[i]) {
			return false
		}
	}
	return true
}

// EqualEvent returns true if the two events are equal.
func EqualEvent(a, b Event) bool {
	return a == b
}
