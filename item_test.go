package budget

import (
	"github.com/dottics/dutil"
	"github.com/google/uuid"
	"github.com/johannesscr/micro/microtest"
	"testing"
)

func TestService_GetItems(t *testing.T) {
	type E struct {
		e        dutil.Error
		items    Items
		exReqURI string
	}
	tt := []struct {
		name     string
		uuid     uuid.UUID
		exchange *microtest.Exchange
		E        E
	}{
		{
			name: "403 Permission Required",
			uuid: uuid.MustParse("a0e09cfc-414d-4b42-9661-333090390a16"),
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
				exReqURI: "/budget/group/-/item?uuid=a0e09cfc-414d-4b42-9661-333090390a16",
				items:    Items{},
				e: &dutil.Err{
					Status: 403,
					Errors: map[string][]string{
						"auth": {"Please ensure you have permission"},
					},
				},
			},
		},
		{
			name: "404 Not Found",
			uuid: uuid.MustParse("0db30884-59ab-4214-8ae1-d3a1a3ae81c9"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 404,
					Body: `{
						"message":"NotFound: unable to process request",
						"data":{},
						"errors":{
							"group":["not found"]
						}
					}`,
				},
			},
			E: E{
				exReqURI: "/budget/group/-/item?uuid=0db30884-59ab-4214-8ae1-d3a1a3ae81c9",
				items:    Items{},
				e: &dutil.Err{
					Status: 401,
					Errors: map[string][]string{
						"group": {"not found"},
					},
				},
			},
		},
		{
			name: "200 Successful Empty",
			uuid: uuid.MustParse("7fa5252e-faaf-4486-ba68-ca6d00c203cf"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message":"NotFound: unable to process request",
						"data":{"items":[]},
						"errors":{}
					}`,
				},
			},
			E: E{
				exReqURI: "/budget/group/-/item?uuid=7fa5252e-faaf-4486-ba68-ca6d00c203cf",
				items:    Items{},
				e:        nil,
			},
		},
		{
			name: "200 Successful",
			uuid: uuid.MustParse("4694620e-a67e-418e-8b41-44a2413a6450"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message":"NotFound: unable to process request",
						"data":{
							"items":[
								{
									"uuid":"d2b64e51-8b31-4cbd-be90-439ddb33c3b7",
									"name":"item one",
									"active":true
								},
								{
									"uuid":"c1d396ab-4e32-4d1e-9baf-48a10529cf80",
									"name":"item two",
									"active":true
								},
								{
									"uuid":"8e5ec8c2-f89b-464f-ba66-06f9365ebb2b",
									"name":"item three",
									"active":true
								}
							]
						},
						"errors":{}
					}`,
				},
			},
			E: E{
				exReqURI: "/budget/group/-/item?uuid=4694620e-a67e-418e-8b41-44a2413a6450",
				items: Items{
					Item{
						UUID:   uuid.MustParse("d2b64e51-8b31-4cbd-be90-439ddb33c3b7"),
						Name:   "item one",
						Active: true,
					},
					Item{
						UUID:   uuid.MustParse("c1d396ab-4e32-4d1e-9baf-48a10529cf80"),
						Name:   "item two",
						Active: true,
					},
					Item{
						UUID:   uuid.MustParse("8e5ec8c2-f89b-464f-ba66-06f9365ebb2b"),
						Name:   "item three",
						Active: true,
					},
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
			// add budget-mirco-service exchange
			ms.Append(tc.exchange)

			xi, e := s.GetItems(tc.uuid)
			// test errors
			if tc.E.e != nil {
				if tc.E.e.Error() != e.Error() {
					t.Errorf("expected error '%v' got '%v'", tc.E.e.Error(), e.Error())
				}
			} else if e != nil {
				t.Errorf("unexpected err: %s", e.Error())
			}

			// test items
			if len(tc.E.items) != len(xi) {
				t.Errorf("expected items length %d got %d", len(tc.E.items), len(xi))
			}
			for i, item := range tc.E.items {
				j := xi[i]
				if item.UUID != j.UUID {
					t.Errorf("expected uuid '%v' got '%v'", item.UUID, j.UUID)
				}
				if item.Name != j.Name {
					t.Errorf("expected name '%v' got '%v'", item.Name, j.Name)
				}
				if item.Active != j.Active {
					t.Errorf("expected active '%v' got '%v'", item.Active, j.Active)
				}
			}

			// test exchange request
			if tc.exchange.Request.RequestURI != tc.E.exReqURI {
				t.Errorf("expected URI '%v' got '%v'", tc.E.exReqURI, tc.exchange.Request.RequestURI)
			}
		})
	}
}
