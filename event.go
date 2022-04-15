package budget

import (
	"github.com/dottics/dutil"
	"github.com/google/uuid"
	"net/url"
)

// GetEvents retrieves all the events from the budget-micro-service that are
// related to an item.
func (s *Service) GetEvents(UUID uuid.UUID) (Events, dutil.Error) {
	s.URL.Path = "/budget/group/item/-/event"
	q := url.Values{
		"uuid": {UUID.String()},
	}
	s.URL.RawQuery = q.Encode()

	type data struct {
		Events Events `json:"events"`
	}
	resp := struct {
		Message string              `json:"message"`
		Data    data                `json:"data"`
		Errors  map[string][]string `json:"errors"`
	}{}

	res, e := s.newRequest("GET", s.URL.String(), nil, nil)
	if e != nil {
		return nil, e
	}
	_, e = s.decode(res, &resp)
	if e != nil {
		return nil, e
	}

	if res.StatusCode != 200 {
		e := &dutil.Err{
			Status: res.StatusCode,
			Errors: resp.Errors,
		}
		return nil, e
	}

	return resp.Data.Events, nil
}

// CreateEvent makes the request to the budget-micro-service to create a new
// event and associate that event with an item
func (s *Service) CreateEvent(UUID uuid.UUID, event Event) (Event, dutil.Error) {
	s.URL.Path = "/event"

	p := struct {
		ItemUUID uuid.UUID `json:"item_uuid"`
		Event    Event     `json:"event"`
	}{
		ItemUUID: UUID,
		Event:    event,
	}

	payload, e := dutil.MarshalReader(p)
	if e != nil {
		return Event{}, e
	}
	res, e := s.newRequest("POST", s.URL.String(), nil, payload)
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

	_, e = s.decode(res, &resp)
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

func (s *Service) UpdateEvent(event Event) (Event, dutil.Error) {
	s.URL.Path = "/event/-"

	payload, e := dutil.MarshalReader(event)
	if e != nil {
		return Event{}, e
	}
	res, e := s.newRequest("PUT", s.URL.String(), nil, payload)
	if e != nil {
		return Event{}, nil
	}

	type data struct {
		Event Event `json:"event"`
	}
	resp := struct {
		Message string              `json:"message"`
		Data    data                `json:"data"`
		Errors  map[string][]string `json:"errors"`
	}{}
	_, e = s.decode(res, &resp)
	if e != nil {
		return Event{}, e
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

func (s *Service) DeleteEvent(UUID uuid.UUID) dutil.Error {
	s.URL.Path = "/event/-"
	q := url.Values{
		"uuid": {UUID.String()},
	}
	s.URL.RawQuery = q.Encode()

	res, e := s.newRequest("DELETE", s.URL.String(), nil, nil)
	if e != nil {
		return e
	}

	resp := struct {
		Message string              `json:"message"`
		Data    map[string]string   `json:"data"`
		Errors  map[string][]string `json:"errors"`
	}{}
	_, e = s.decode(res, &resp)
	if e != nil {
		return e
	}

	if res.StatusCode != 200 {
		e := &dutil.Err{
			Status: res.StatusCode,
			Errors: resp.Errors,
		}
		return e
	}
	return nil
}
