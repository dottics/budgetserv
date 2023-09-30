package budget

import (
	"fmt"
	"github.com/dottics/dutil"
	"github.com/google/uuid"
)

// CreateEvent makes the request to the budget-microservice to create a new
// event and associate that event with an item
func (s *Service) CreateEvent(event EventCreate) (Event, error) {
	s.URL.Path = "/event/"

	p, err := marshalReader(event)
	if err != nil {
		return Event{}, err
	}
	res, e := s.DoRequest("POST", s.URL, nil, nil, p)
	if e != nil {
		return Event{}, e
	}

	type data struct {
		Event Event `json:"event"`
	}
	resp := struct {
		Message string `json:"message"`
		Data    data   `json:"data"`
	}{}

	err = marshalResponse(201, res, &resp)
	if err != nil {
		return Event{}, err
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
		Message string `json:"message"`
	}{}

	err := marshalResponse(200, res, &resp)
	if err != nil {
		return err
	}

	return nil
}
