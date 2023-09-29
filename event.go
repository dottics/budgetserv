package budget

import (
	"fmt"
	"github.com/dottics/dutil"
	"github.com/google/uuid"
	"github.com/johannesscr/micro/msp"
	"time"
)

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

// CreateEvent makes the request to the budget-microservice to create a new
// event and associate that event with an item
func (s *Service) CreateEvent(event EventCreate) (Event, dutil.Error) {
	s.URL.Path = "/event"

	payload, e := dutil.MarshalReader(event)
	if e != nil {
		return Event{}, e
	}
	res, e := s.DoRequest("POST", s.URL, nil, nil, payload)
	if e != nil {
		return Event{}, e
	}

	type data struct {
		Event Event `json:"event"`
	}
	resp := struct {
		Message string              `json:"message"`
		Data    data                `json:"data"`
		Errors  map[string][]string `json:"errors"`
	}{}

	_, e = msp.Decode(res, &resp)
	if e != nil {
		return Event{}, nil
	}

	if res.StatusCode != 200 {
		e := &dutil.Err{
			Status: res.StatusCode,
			Errors: resp.Errors,
		}
		return Event{}, e
	}

	return resp.Data.Event, nil
}

// UpdateEvent makes the request to the budget-microservice to update an event.
func (s *Service) UpdateEvent(UUID uuid.UUID, event EventUpdate) (Event, error) {
	s.URL.Path = fmt.Sprintf("/event/%s", UUID.String())

	payload, e := dutil.MarshalReader(event)
	if e != nil {
		return Event{}, e
	}
	res, e := s.DoRequest("PUT", s.URL, nil, nil, payload)
	if e != nil {
		return Event{}, nil
	}

	type data struct {
		Event Event `json:"event"`
	}
	resp := struct {
		Message string `json:"message"`
		Data    data   `json:"data"`
	}{}

	err := marshalResponse(200, res, &resp)
	if err != nil {
		return Event{}, err
	}

	return resp.Data.Event, nil
}

// DeleteEvent makes the request to the budget-microservice to delete an event.
func (s *Service) DeleteEvent(UUID uuid.UUID) error {
	s.URL.Path = fmt.Sprintf("/event/%s", UUID.String())

	res, e := s.DoRequest("DELETE", s.URL, nil, nil, nil)
	if e != nil {
		return e
	}

	resp := struct {
		Message string            `json:"message"`
		Data    map[string]string `json:"data"`
	}{}

	err := marshalResponse(200, res, &resp)
	if err != nil {
		return err
	}

	return nil
}
