package budget

import (
	"github.com/dottics/dutil"
	"github.com/google/uuid"
	"github.com/johannesscr/micro/microtest"
	"testing"
	"time"
)

func TestService_GetEvents(t *testing.T) {
	type E struct {
		exReqURI string
		events   Events
		e        dutil.Error
	}
	tt := []struct {
		name     string
		uuid     uuid.UUID
		exchange *microtest.Exchange
		E        E
	}{
		{
			name: "403 Permission Required",
			uuid: uuid.MustParse("9db1590a-9f77-47af-baa6-0a786095e510"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body: `{
						"message":"Forbidden: unable to process request",
						"data":{},
						"errors":{
							"auth":["Please ensure you have permission"]
						}
					}`,
				},
			},
			E: E{
				exReqURI: "/budget/group/item/-/event?uuid=9db1590a-9f77-47af-baa6-0a786095e510",
				events:   Events{},
				e: &dutil.Err{
					Status: 403,
					Errors: map[string][]string{
						"auth": {"Please ensure you have permission"},
					},
				},
			},
		},
		{
			name: "500 Internal Server Error",
			uuid: uuid.MustParse("08aa5301-2285-4f2f-b930-677683468e0f"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 500,
					Body: `{
						"message":"InternalServerError: unable to process request",
						"data":{},
						"errors":{
							"internal_server_error":["some internal server error"]
						}
					}`,
				},
			},
			E: E{
				exReqURI: "/budget/group/item/-/event?uuid=08aa5301-2285-4f2f-b930-677683468e0f",
				events:   Events{},
				e: &dutil.Err{
					Status: 500,
					Errors: map[string][]string{
						"internal_server_error": {"some internal server error"},
					},
				},
			},
		},
		{
			name: "200 Successful Events",
			uuid: uuid.MustParse("bf905480-9c0f-42d3-85c0-cc98dd63be5c"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message":"events found successfully",
						"data":{
							"events":[
								{
									"uuid":"d89b8fc2-9e64-4706-a4db-17ffbe6d16a9",
									"name":"event one",
									"debit":true,
									"credit":false,
									"amount":1205.24,
									"start_date":"2021-11-12T00:00:00Z",
									"end_date":"2021-11-12T00:00:00Z",
									"active":true
								},
								{
									"uuid":"6aa6d86e-8bcc-4706-a2d1-82bbfc7e7c97",
									"name":"event two",
									"debit":true,
									"credit":false,
									"amount":178,
									"start_date":"2021-11-12T00:00:00Z",
									"end_date":"2021-11-12T00:00:00Z",
									"active":true
								},
								{
									"uuid":"603f9589-ecb3-437a-9406-1652eb5a38eb",
									"name":"event three",
									"debit":false,
									"credit":true,
									"amount":0.35,
									"start_date":"2021-11-12T00:00:00Z",
									"end_date":"2021-11-12T00:00:00Z",
									"active":true
								}
							]
						},
						"errors":{}
					}`,
				},
			},
			E: E{
				exReqURI: "/budget/group/item/-/event?uuid=bf905480-9c0f-42d3-85c0-cc98dd63be5c",
				events: Events{
					Event{
						UUID:   uuid.MustParse("d89b8fc2-9e64-4706-a4db-17ffbe6d16a9"),
						Name:   "event one",
						Debit:  true,
						Credit: false,
						Amount: 1205.24,
						Active: true,
					},
					Event{
						UUID:   uuid.MustParse("6aa6d86e-8bcc-4706-a2d1-82bbfc7e7c97"),
						Name:   "event two",
						Debit:  true,
						Credit: false,
						Amount: 178,
						Active: true,
					},
					Event{
						UUID:   uuid.MustParse("603f9589-ecb3-437a-9406-1652eb5a38eb"),
						Name:   "event three",
						Debit:  false,
						Credit: true,
						Amount: 0.35,
						Active: true,
					},
				},
				e: nil,
			},
		},
	}

	s := NewService("")
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add budget-micro-service exchange
			ms.Append(tc.exchange)

			xe, e := s.GetEvents(tc.uuid)
			// test errors
			if tc.E.e != nil {
				if tc.E.e.Error() != e.Error() {
					t.Errorf("expected error '%v' got '%v'", tc.E.e.Error(), e.Error())
				}
			} else if e != nil {
				t.Errorf("unexpected error: %s", e.Error())
			}
			// test events
			if len(tc.E.events) != len(xe) {
				t.Errorf("expected events length %d got %d", len(tc.E.events), len(xe))
			}
			for i, event := range tc.E.events {
				ev := xe[i]
				if ev.UUID != event.UUID {
					t.Errorf("expected event uuid '%v' got '%v'", event.UUID, ev.UUID)
				}
				if ev.Name != event.Name {
					t.Errorf("expected event name '%v' got '%v'", event.Name, ev.Name)
				}
				if ev.Debit != event.Debit {
					t.Errorf("expected event debit '%v' got '%v'", event.Debit, ev.Debit)
				}
				if ev.Credit != event.Credit {
					t.Errorf("expected event credit '%v' got '%v'", event.Credit, ev.Credit)
				}
				if ev.Amount != event.Amount {
					t.Errorf("expected event credit %f got %f", event.Amount, ev.Amount)
				}
				if ev.Active != event.Active {
					t.Errorf("expected event active %v got %v", event.Active, ev.Active)
				}
			}
			// test query string
			if tc.exchange.Request.RequestURI != tc.E.exReqURI {
				t.Errorf("expected exchange request URI '%s' got '%s'", tc.E.exReqURI, tc.exchange.Request.RequestURI)
			}
		})
	}
}

func TestService_CreateEvent(t *testing.T) {
	type E struct {
		event Event
		e dutil.Error
	}
	tt := []struct {
		name string
		UUID uuid.UUID
		event Event
		exchange *microtest.Exchange
		E E
	}{
		{
			name: "403 Permission Required",
			UUID: uuid.MustParse("ae9f5130-81fe-4526-9573-f7e892cc2e01"),
			event: Event{
				Name: "test event one",
				Amount: 12.19,
				Debit: true,
				Credit: false,
				StartDate: time.Date(2021, 11, 18, 0, 0, 0, 0, time.UTC),
				EndDate: time.Date(2021, 11, 19, 0, 0, 0, 0, time.UTC),
				Active: true,
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
			event: Event{
				Name: "test event two",
				Amount: 12.19,
				Debit: true,
				Credit: false,
				StartDate: time.Date(2021, 11, 18, 0, 0, 0, 0, time.UTC),
				EndDate: time.Date(2021, 11, 19, 0, 0, 0, 0, time.UTC),
				Active: true,
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
			event: Event{
				Name: "test event three",
				Amount: 12.19,
				Debit: true,
				Credit: false,
				StartDate: time.Date(2021, 11, 18, 0, 0, 0, 0, time.UTC),
				EndDate: time.Date(2021, 11, 19, 0, 0, 0, 0, time.UTC),
				Active: true,
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
					UUID: uuid.MustParse("b86768ee-69de-4fb2-81eb-ab96d14e37ae"),
					Name: "test event three",
					Amount: 12.19,
					Debit: true,
					Credit: false,
					StartDate: time.Date(2021, 11, 18, 0, 0, 0, 0, time.UTC),
					EndDate: time.Date(2021, 11, 19, 0, 0, 0, 0, time.UTC),
					Active: true,
				},
				e: nil,
			},
		},
	}

	s := NewService("")
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add budget-micro-service exchange
			ms.Append(tc.exchange)

			ev, e := s.CreateEvent(tc.UUID, tc.event)

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
		e dutil.Error
	}
	tt := []struct {
		name string
		event Event
		exchange *microtest.Exchange
		E E
	}{
		{
			name: "error",
			event: Event{
				UUID: uuid.MustParse("b86768ee-69de-4fb2-81eb-ab96d14e37ae"),
				Name: "test event",
				Debit: true,
				Credit: false,
				StartDate: time.Date(2021, 11, 18, 0, 0, 0, 0, time.UTC),
				EndDate: time.Date(2021, 11, 19, 0, 0, 0, 0, time.UTC),
				Active: true,
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
			name: "success",
			event: Event{
				UUID: uuid.MustParse("b86768ee-69de-4fb2-81eb-ab96d14e37ae"),
				Name: "test event",
				Debit: false,
				Credit: true,
				Amount: 125.67,
				StartDate: time.Date(2021, 11, 18, 0, 0, 0, 0, time.UTC),
				EndDate: time.Date(2021, 11, 19, 0, 0, 0, 0, time.UTC),
				Active: true,
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
					UUID: uuid.MustParse("b86768ee-69de-4fb2-81eb-ab96d14e37ae"),
					Name: "test event",
					Debit: false,
					Credit: true,
					Amount: 125.67,
					StartDate: time.Date(2021, 11, 18, 0, 0, 0, 0, time.UTC),
					EndDate: time.Date(2021, 11, 19, 0, 0, 0, 0, time.UTC),
					Active: true,
				},
				e: nil,
			},
		},
	}

	s := NewService("")
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add budget-micro-service exchange
			ms.Append(tc.exchange)

			ev, e := s.UpdateEvent(tc.event)
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
				t.Errorf("expected event '%v' got '%v'", tc.E.event, ev)
			}
		})
	}
}

func TestService_DeleteEvent(t *testing.T) {
	type E struct  {
		e dutil.Error
		exReqURI string
	}
	tt := []struct {
		name string
		UUID uuid.UUID
		exchange *microtest.Exchange
		E E
	}{
		{
			name: "403 Bad Request",
			UUID: uuid.MustParse("7dee09f0-2ba2-4b10-9c88-0c973ad4ebd0"),
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
				exReqURI: "/event/-?uuid=7dee09f0-2ba2-4b10-9c88-0c973ad4ebd0",
				e: &dutil.Err{
					Status: 403,
					Errors: map[string][]string{
						"permission": {"Please ensure you have permission"},
					},
				},
			},
		},
		{
			name: "200 Success",

			UUID: uuid.MustParse("2c7b7d76-c9e0-49f8-b585-166ac70dba6f"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message":"event deleted successfully",
						"data":{},
						"errors":{}
					}`,
				},
			},
			E: E{
				exReqURI: "/event/-?uuid=2c7b7d76-c9e0-49f8-b585-166ac70dba6f",
				e: nil,
			},
		},
	}

	s := NewService("")
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