package budget

import (
	"fmt"
	"github.com/dottics/dutil"
	"github.com/google/uuid"
	"github.com/johannesscr/micro/msp"
)

// TODO: remove this: endpoint no longer exists
//// GetEvents retrieves all the events from the budget-micro-service that are
//// related to an item.
//func (s *Service) GetEvents(UUID uuid.UUID) (Events, dutil.Error) {
//	s.URL.Path = "/budget/group/item/-/event"
//	q := url.Values{
//		"uuid": {UUID.String()},
//	}
//	s.URL.RawQuery = q.Encode()
//
//	type data struct {
//		Events Events `json:"events"`
//	}
//	resp := struct {
//		Message string              `json:"message"`
//		Data    data                `json:"data"`
//		Errors  map[string][]string `json:"errors"`
//	}{}
//
//	res, e := s.DoRequest("GET", s.URL, nil, nil, nil)
//	if e != nil {
//		return nil, e
//	}
//	_, e = msp.Decode(res, &resp)
//	if e != nil {
//		return nil, e
//	}
//
//	if res.StatusCode != 200 {
//		e := &dutil.Err{
//			Status: res.StatusCode,
//			Errors: resp.Errors,
//		}
//		return nil, e
//	}
//
//	return resp.Data.Events, nil
//}

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

func (s *Service) UpdateEvent(event Event) (Event, dutil.Error) {
	s.URL.Path = fmt.Sprintf("/event/%s", event.UUID.String())

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
		Message string              `json:"message"`
		Data    data                `json:"data"`
		Errors  map[string][]string `json:"errors"`
	}{}
	_, e = msp.Decode(res, &resp)
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
	s.URL.Path = fmt.Sprintf("/event/%s", UUID.String())

	res, e := s.DoRequest("DELETE", s.URL, nil, nil, nil)
	if e != nil {
		return e
	}

	resp := struct {
		Message string              `json:"message"`
		Data    map[string]string   `json:"data"`
		Errors  map[string][]string `json:"errors"`
	}{}
	_, e = msp.Decode(res, &resp)
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
