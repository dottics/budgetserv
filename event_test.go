package budget

import (
	"errors"
	"github.com/google/uuid"
	"github.com/johannesscr/micro/microtest"
	"testing"
	"time"
)

func TestService_CreateEvent(t *testing.T) {
	tt := []struct {
		name     string
		UUID     uuid.UUID
		payload  EventCreate
		exchange *microtest.Exchange
		event    Event
		e        error
	}{
		{
			name: "403 Permission Required",
			UUID: uuid.MustParse("ae9f5130-81fe-4526-9573-f7e892cc2e01"),
			payload: EventCreate{
				Name:      "test event one",
				Amount:    12.19,
				Debit:     true,
				Credit:    false,
				StartDate: time.Date(2021, 11, 18, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2021, 11, 19, 0, 0, 0, 0, time.UTC),
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   noPermission,
				},
			},
			event: Event{},
			e:     errors.New("no permission"),
		},
		{
			name: "200 Successful",
			UUID: uuid.MustParse("ae9f5130-81fe-4526-9573-f7e892cc2e01"),
			payload: EventCreate{
				Name:      "test event three",
				Amount:    12.19,
				Debit:     true,
				Credit:    false,
				StartDate: time.Date(2021, 11, 18, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2021, 11, 19, 0, 0, 0, 0, time.UTC),
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 201,
					Body: `{
						"message":"event created successfully",
						"data":{
							"event":{
								"uuid":"b86768ee-69de-4fb2-81eb-ab96d14e37ae",
								"name":"test event three",
								"amount":12.19,
								"debit":true,
								"credit":false,
								"start_date":"2021-11-18T00:00:00Z",
								"end_date":"2021-11-19T00:00:00Z",
								"active":true
							}
						},
						"errors":{}
					}`,
				},
			},
			event: Event{
				UUID:      uuid.MustParse("b86768ee-69de-4fb2-81eb-ab96d14e37ae"),
				Name:      "test event three",
				Amount:    12.19,
				Debit:     true,
				Credit:    false,
				StartDate: time.Date(2021, 11, 18, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2021, 11, 19, 0, 0, 0, 0, time.UTC),
			},
			e: nil,
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add budget-micro-service exchange
			ms.Append(tc.exchange)

			ev, e := s.CreateEvent(tc.payload)

			// test errors
			if NotEqualError(e, tc.e) {
				t.Errorf("expected error '%v' got '%v'", tc.e, e)
			}

			// test event
			if ev != tc.event {
				t.Errorf("expected an event\n'%+v'\ngot\n'%+v'", tc.event, ev)
			}

			// TODO: test exchange request body
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	tt := []struct {
		name     string
		uuid     uuid.UUID
		payload  EventUpdate
		exchange *microtest.Exchange
		uri      string
		event    Event
		e        error
	}{
		{
			name: "error",
			uuid: uuid.MustParse("8feec066-2dfb-44f5-b353-1cb6e75c3084"),
			payload: EventUpdate{
				ItemUUID:  uuid.MustParse("b86768ee-69de-4fb2-81eb-ab96d14e37ae"),
				Name:      "test event",
				Debit:     true,
				Credit:    false,
				StartDate: time.Date(2021, 11, 18, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2021, 11, 19, 0, 0, 0, 0, time.UTC),
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   noPermission,
				},
			},
			uri:   "/event/8feec066-2dfb-44f5-b353-1cb6e75c3084",
			event: Event{},
			e:     errors.New("no permission"),
		},
		{
			name: "success",
			uuid: uuid.MustParse("0aab721b-4224-474c-9df8-e77117a31a02"),
			payload: EventUpdate{
				ItemUUID:  uuid.MustParse("b86768ee-69de-4fb2-81eb-ab96d14e37ae"),
				Name:      "test event",
				Debit:     false,
				Credit:    true,
				Amount:    125.67,
				StartDate: time.Date(2021, 11, 18, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2021, 11, 19, 0, 0, 0, 0, time.UTC),
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message":"event updated successfully",
						"data":{
							"event":{
								"uuid":"0aab721b-4224-474c-9df8-e77117a31a02",
								"name":"test event",
								"debit":false,
								"credit":true,
								"amount":125.67,
								"start_date":"2021-11-18T00:00:00Z",
								"end_date":"2021-11-19T00:00:00Z",
								"active":true
							}
						},
						"errors":{}
					}`,
				},
			},
			uri: "/event/0aab721b-4224-474c-9df8-e77117a31a02",
			event: Event{
				UUID:      uuid.MustParse("0aab721b-4224-474c-9df8-e77117a31a02"),
				Name:      "test event",
				Debit:     false,
				Credit:    true,
				Amount:    125.67,
				StartDate: time.Date(2021, 11, 18, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2021, 11, 19, 0, 0, 0, 0, time.UTC),
			},
			e: nil,
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add budget-micro-service exchange
			ms.Append(tc.exchange)

			ev, e := s.UpdateEvent(tc.uuid, tc.payload)
			// test error
			if NotEqualError(e, tc.e) {
				t.Errorf("expected error '%v' got '%v'", tc.e, e)
			}
			// test event
			if tc.event != ev {
				t.Errorf("expected event '%+v' got '%+v'", tc.event, ev)
			}
			// test exchange request uri
			if tc.uri != tc.exchange.Request.RequestURI {
				t.Errorf("expected uri '%s' got '%s'", tc.uri, tc.exchange.Request.RequestURI)
			}
		})
	}
}

func TestService_DeleteEvent(t *testing.T) {
	tt := []struct {
		name     string
		UUID     uuid.UUID
		exchange *microtest.Exchange
		uri      string
		e        error
	}{
		{
			name: "403 Bad Request",
			UUID: uuid.MustParse("7dee09f0-2ba2-4b10-9c88-0c973ad4ebd0"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   noPermission,
				},
			},
			uri: "/event/7dee09f0-2ba2-4b10-9c88-0c973ad4ebd0",
			e:   errors.New("no permission"),
		},
		{
			name: "200 Success",

			UUID: uuid.MustParse("2c7b7d76-c9e0-49f8-b585-166ac70dba6f"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message":"event deleted successfully"
					}`,
				},
			},
			uri: "/event/2c7b7d76-c9e0-49f8-b585-166ac70dba6f",
			e:   nil,
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add budget-micro-service exchange
			ms.Append(tc.exchange)

			e := s.DeleteEvent(tc.UUID)
			if NotEqualError(e, tc.e) {
				t.Errorf("expected error '%v' got '%v'", tc.e, e)
			}

			if tc.exchange.Request.RequestURI != tc.uri {
				t.Errorf("expected URI '%s' got '%s'", tc.uri, tc.exchange.Request.RequestURI)
			}
		})
	}
}
