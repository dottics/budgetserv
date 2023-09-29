package budget

import (
	"errors"
	"github.com/dottics/dutil"
	"github.com/google/uuid"
	"github.com/johannesscr/micro/microtest"
	"testing"
	"time"
)

func TestService_CreateEvent(t *testing.T) {
	type E struct {
		event Event
		e     dutil.Error
	}
	tt := []struct {
		name     string
		UUID     uuid.UUID
		event    EventCreate
		exchange *microtest.Exchange
		E        E
	}{
		{
			name: "403 Permission Required",
			UUID: uuid.MustParse("ae9f5130-81fe-4526-9573-f7e892cc2e01"),
			event: EventCreate{
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
					Body: `{
						"message":"Forbidden: unable to process request",
						"data":{},
						"errors":{
							"permission":["Please ensure you have permission"]
						}
					}`,
				},
			},
			E: E{
				event: Event{},
				e: &dutil.Err{
					Status: 403,
					Errors: map[string][]string{
						"permission": {"Please ensure you have permission"},
					},
				},
			},
		},
		{
			name: "500 Internal Server Error",
			UUID: uuid.MustParse("ae9f5130-81fe-4526-9573-f7e892cc2e01"),
			event: EventCreate{
				Name:      "test event two",
				Amount:    12.19,
				Debit:     true,
				Credit:    false,
				StartDate: time.Date(2021, 11, 18, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2021, 11, 19, 0, 0, 0, 0, time.UTC),
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 500,
					Body: `{
						"message":"InternalServerError: unable to process request",
						"data":{},
						"errors":{
							"internal_server_error":["some random internal error"]
						}
					}`,
				},
			},
			E: E{
				event: Event{},
				e: &dutil.Err{
					Status: 500,
					Errors: map[string][]string{
						"internal_server_error": {"some random internal error"},
					},
				},
			},
		},
		{
			name: "200 Successful",
			UUID: uuid.MustParse("ae9f5130-81fe-4526-9573-f7e892cc2e01"),
			event: EventCreate{
				Name:      "test event three",
				Amount:    12.19,
				Debit:     true,
				Credit:    false,
				StartDate: time.Date(2021, 11, 18, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2021, 11, 19, 0, 0, 0, 0, time.UTC),
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
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
			E: E{
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
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add budget-micro-service exchange
			ms.Append(tc.exchange)

			ev, e := s.CreateEvent(tc.event)

			// test errors
			if tc.E.e != nil {
				if tc.E.e.Error() != e.Error() {
					t.Errorf("expected error '%s' got '%s'", tc.E.e.Error(), e.Error())
				}
			} else if e != nil {
				t.Errorf("unexpected error: %s", e.Error())
			}

			// test event
			if ev != tc.E.event {
				t.Errorf("expected an event '%v' got '%v'", tc.E.event, ev)
			}

			// TODO: test exchange request body
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	type E struct {
		event Event
		e     error
	}
	tt := []struct {
		name     string
		uuid     uuid.UUID
		event    EventUpdate
		exchange *microtest.Exchange
		E        E
	}{
		{
			name: "error",
			event: EventUpdate{
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
			E: E{
				event: Event{},
				e:     errors.New("no permission"),
			},
		},
		{
			name: "success",
			event: EventUpdate{
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
								"uuid":"b86768ee-69de-4fb2-81eb-ab96d14e37ae",
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
			E: E{
				event: Event{
					UUID:      uuid.MustParse("b86768ee-69de-4fb2-81eb-ab96d14e37ae"),
					Name:      "test event",
					Debit:     false,
					Credit:    true,
					Amount:    125.67,
					StartDate: time.Date(2021, 11, 18, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2021, 11, 19, 0, 0, 0, 0, time.UTC),
				},
				e: nil,
			},
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add budget-micro-service exchange
			ms.Append(tc.exchange)

			ev, e := s.UpdateEvent(tc.uuid, tc.event)
			// test error
			if tc.E.e != nil {
				if tc.E.e.Error() != e.Error() {
					t.Errorf("expected error '%s' got '%s'", tc.E.e.Error(), e.Error())
				}
			} else if e != nil {
				t.Errorf("unexpected error: %s", e.Error())
			}
			// test event
			if tc.E.event != ev {
				t.Errorf("expected event '%+v' got '%+v'", tc.E.event, ev)
			}
		})
	}
}

func TestService_DeleteEvent(t *testing.T) {
	type E struct {
		e        error
		exReqURI string
	}
	tt := []struct {
		name     string
		UUID     uuid.UUID
		exchange *microtest.Exchange
		E        E
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
			E: E{
				exReqURI: "/event/7dee09f0-2ba2-4b10-9c88-0c973ad4ebd0",
				e:        errors.New("no permission"),
			},
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
			E: E{
				exReqURI: "/event/2c7b7d76-c9e0-49f8-b585-166ac70dba6f",
				e:        nil,
			},
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
			if tc.E.e != nil {
				if tc.E.e.Error() != e.Error() {
					t.Errorf("expected error '%v' got '%v'", tc.E.e.Error(), e.Error())
				}
			} else if e != nil {
				t.Errorf("unpexpected error: %s", e.Error())
			}
			URI := tc.exchange.Request.RequestURI
			if URI != tc.E.exReqURI {
				t.Errorf("expected URI '%s' got '%s'", tc.E.exReqURI, URI)
			}
		})
	}
}
