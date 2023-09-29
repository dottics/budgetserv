package budget

import (
	"fmt"
	"github.com/dottics/dutil"
	"github.com/google/uuid"
	"github.com/johannesscr/micro/microtest"
	"testing"
)

func TestService_GetGroups(t *testing.T) {
	type E struct {
		groups   Groups
		e        dutil.Error
		exReqURI string
	}
	tt := []struct {
		name       string
		budgetUUID uuid.UUID
		exchange   *microtest.Exchange
		E          E
	}{
		{
			name:       "403 Permission Required",
			budgetUUID: uuid.MustParse("2520f807-915e-41f6-9557-84500e1aebcc"),
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
				exReqURI: "/budget/2520f807-915e-41f6-9557-84500e1aebcc/groups",
				e: &dutil.Err{
					Status: 403,
					Errors: map[string][]string{
						"auth": {"Please ensure you have permission"},
					},
				},
				groups: Groups{},
			},
		},
		{
			name:       "404 Not Found",
			budgetUUID: uuid.MustParse("7cb47f06-0d96-494b-a847-a472e2c04d9d"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 404,
					Body: `{
						"message":"NotFound: unable to process request",
						"data":{},
						"errors":{
							"groups":["not found"]
						}
					}`,
				},
			},
			E: E{
				exReqURI: "/budget/7cb47f06-0d96-494b-a847-a472e2c04d9d/groups",
				e: &dutil.Err{
					Status: 404,
					Errors: map[string][]string{
						"groups": {"not found"},
					},
				},
				groups: Groups{},
			},
		},
		{
			name:       "500 Unmarshal Error",
			budgetUUID: uuid.MustParse("1a252f35-c84f-4937-9c5b-f5deb19b5b10"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 500,
					Body: `{
						"message":"groups found successfully",
						"data":{
							"groups":[
								{
									"uuid":"6be3df72-da3d-4a8c-bef6-d0b57120b80a",
									"name":"income",
									"active":true,
									"sub_groups":[
										{
											"uuid":"ef8e63d5-b26c-4824-8a1a-a729bb8574d3",
											"name":"salary",
											"active":true,
											"sub_groups":[],
										}
									]
								},
								{
									"uuid":"52f2c725-2cdc-401a-abdd-66db5fd06789",
									"name":"Investments",
									"active":true,
									"sub_groups":[]
								}
							]
						},
						"errors":{}
					}`,
				},
			},
			E: E{
				exReqURI: "/budget/1a252f35-c84f-4937-9c5b-f5deb19b5b10/groups",
				e: &dutil.Err{
					Status: 500,
					Errors: map[string][]string{
						"unmarshal": {"invalid character '}' looking for beginning of object key string"},
					},
				},
				groups: Groups{},
			},
		},
		{
			name:       "500 Internal Server Error",
			budgetUUID: uuid.MustParse("3a8113f2-af77-4430-a59e-519a8ad0819d"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 500,
					Body: `{
						"message":"Forbidden: unable to process request",
						"data":{},
						"errors":{
							"internal_server_error":["some unexpected error"]
						}
					}`,
				},
			},
			E: E{
				exReqURI: "/budget/3a8113f2-af77-4430-a59e-519a8ad0819d/groups",
				e: &dutil.Err{
					Status: 500,
					Errors: map[string][]string{
						"internal_server_error": {"some unexpected error"},
					},
				},
				groups: Groups{},
			},
		},
		{
			name:       "200 Successful",
			budgetUUID: uuid.MustParse("b440353e-cc26-449c-a470-e0e36a2919a6"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message":"groups found successfully",
						"data":{
							"groups":[
								{
									"uuid":"52f2c725-2cdc-401a-abdd-66db5fd06789",
									"name":"income",
									"active":true,
									"sub_groups":[
										{
											"uuid":"b8448a78-6417-4fe2-849c-024622bc6106",
											"name":"base salary",
											"active":true,
											"sub_groups":[]
										}
									]
								},
								{
									"uuid":"eea51d45-c9bd-45e2-bc80-010ecbb7a0d3",
									"name":"investments",
									"active":true,
									"sub_groups":[]
								},
								{
									"uuid":"6be3df72-da3d-4a8c-bef6-d0b57120b80a",
									"name":"expenses",
									"active":true,
									"sub_groups":[]
								}
							]
						},
						"errors":{}
					}`,
				},
			},
			E: E{
				exReqURI: "/budget/b440353e-cc26-449c-a470-e0e36a2919a6/groups",
				e:        nil,
				groups: Groups{
					Group{
						UUID: uuid.MustParse("52f2c725-2cdc-401a-abdd-66db5fd06789"),
						Name: "income",
						SubGroups: Groups{
							Group{
								UUID: uuid.MustParse("b8448a78-6417-4fe2-849c-024622bc6106"),
								Name: "base salary",
							},
						},
					},
					Group{
						UUID: uuid.MustParse("eea51d45-c9bd-45e2-bc80-010ecbb7a0d3"),
						Name: "investments",
					},
					Group{
						UUID: uuid.MustParse("6be3df72-da3d-4a8c-bef6-d0b57120b80a"),
						Name: "expenses",
					},
				},
			},
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add the budget-micro-service exchange
			ms.Append(tc.exchange)

			xg, e := s.GetGroups(tc.budgetUUID)
			// test the error response
			if tc.E.e != nil {
				if tc.E.e.Error() != e.Error() {
					t.Errorf("expected error '%v' got '%v'", tc.E.e.Error(), e.Error())
				}
			} else if e != nil {
				t.Errorf("unexpected error: %s", e.Error())
			}

			// test the groups structure returned
			if len(xg) != len(tc.E.groups) {
				t.Errorf("expected len groups %d got %d", len(tc.E.groups), len(xg))
			}
			seg := fmt.Sprintf("%v", tc.E.groups)
			sxg := fmt.Sprintf("%v", xg)
			if seg != sxg {
				t.Errorf("expected groups '%v' got '%v'", seg, sxg)
			}
			// test the exchange request URI
			if tc.exchange.Request.RequestURI != tc.E.exReqURI {
				t.Errorf("expected uri '%v' got '%v'", tc.E.exReqURI, tc.exchange.Request.RequestURI)
			}
		})
	}
}
