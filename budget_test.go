package budget

import (
	"github.com/dottics/dutil"
	"github.com/google/uuid"
	"github.com/johannesscr/micro/microtest"
	"testing"
)


func TestService_GetBudgets(t *testing.T) {
	type E struct {
		status int
		len int
		e dutil.Error
	}

	tt := []struct{
		name string
		exchange *microtest.Exchange
		E E
	}{
		{
			name: "403 Forbidden",
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body: `{
						"message":"Forbidden: unable to process request",
						"data":{},
						"errors": {
							"auth": ["Please ensure you have permissions"]
						}
					}`,
				},
			},
			E: E{
				status: 403,
				len: 0,
				e: &dutil.Err{
					Status: 403,
					Errors: map[string][]string{
						"auth": {"Please ensure you have permissions"},
					},
				},
			},
		},
		{
			name: "500 Internal Server Error",
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 500,
					Body: `{
						"message":"InternalServerError",
						"data":{},
						"errors":{
							"internal_server_error": ["some unexpected error"]
						}
					}`,
				},
			},
			E: E{
				status: 500,
				len: 0,
				e: &dutil.Err{
					Status: 500,
					Errors: map[string][]string{
						"internal_server_error": {"some unexpected error"},
					},
				},
			},
		},
		{
			name: "500 Unmarshal Error",
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 500,
					Body: `{
						"message":"",
						"data":{
							"budgets":[
								{
									"uuid":"117d4612-0c95-40b0-8544-c84c1af5407e",
									"name":"budget name"
									"active":"missing comma above"
								}
							]
						},
						"errors":{}
					}`,
				},
			},
			E: E{
				status: 500,
				len: 0,
				e: &dutil.Err{
					Status: 500,
					Errors: map[string][]string{
						"unmarshal": {"invalid character '\"' after object key:value pair"},
					},
				},
			},
		},
		{
			name: "200 Successful",
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message":"budgets found successfully",
						"data":{
							"budgets":[
								{
									"uuid":"4ebae4ad-803c-4487-98ea-3f1f926e59e6",
									"user_uuid":"ecdccfe9-95fe-4c9f-bd86-169ad67c445a",
									"organisation_uuid":null,
									"name":"test budget uno",
									"active":true
								},
								{
									"uuid":"0d79e5cb-5b26-49bc-a5fa-3b39e2710675",
									"user_uuid":"ecdccfe9-95fe-4c9f-bd86-169ad67c445a",
									"organisation_uuid":null,
									"name":"test budget dos",
									"active":true
								}
							]
						},
						"errors":{}
					}`,
				},
			},
			E: E{
				status: 200,
				len: 2,
				e: nil,
			},
		},
	}

	s := NewService("")
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// set micro-service mock responses
			ms.Append(tc.exchange)

			xb, e := s.GetBudgets()

			// test errors
			if tc.E.e != nil {
				if tc.E.e.Error() != e.Error() {
					t.Errorf("expected '%v' got '%v'", tc.E.e.Error(), e.Error())
				}
			} else if e != nil {
				t.Errorf("unexpected err: %s", e.Error())
			}

			if tc.E.len != len(xb) {
				t.Errorf("expected %d got %d", tc.E.len, len(xb))
			}
		})
	}
}

func TestService_GetBudget(t *testing.T) {
	type E struct {
		budget Budget
		e dutil.Error
	}

	tt := []struct{
		name string
		uuid uuid.UUID
		exchange *microtest.Exchange
		E E
	}{
		{
			name: "403 Forbidden",
			uuid: uuid.New(),
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
				budget: Budget{},
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
			uuid: uuid.New(),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 404,
					Body: `{
						"message":"NotFound: unable to process request",
						"data":{},
						"errors":{
							"budget":["not found"]
						}
					}`,
				},
			},
			E: E{
				budget: Budget{},
				e: &dutil.Err{
					Status: 404,
					Errors: map[string][]string{
						"budget": {"not found"},
					},
				},
			},
		},
		{
			name: "500 Internal Server Error",
			uuid: uuid.New(),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 500,
					Body: `{
						"message":"InternalServerError: unable to process request",
						"data":{},
						"errors":{
							"internal_server_error":["some unexpected error"]
						}
					}`,
				},
			},
			E: E{
				budget: Budget{},
				e: &dutil.Err{
					Status: 500,
					Errors: map[string][]string{
						"internal_server_error": {"some unexpected error"},
					},
				},
			},
		},
		{
			name: "500 Unmarshal Error",
			uuid: uuid.New(),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message":"budget found successfully",
						"data":{
							"budget":{
								"uuid":"f5fca9d0-e308-4ff2-be4e-aff22a4c2a78",
								"user_uuid":"67b14c0f-b8ea-4f0f-bf07-cadc73cd74d9",
								"organisation_uuid":null,
								"name":"test budget"
								"active":true
							}
						},
						"errors":{}
					}`,
				},
			},
			E: E{
				budget: Budget{},
				e: &dutil.Err{
					Status: 500,
					Errors: map[string][]string{
						"unmarshal": {"invalid character '\"' after object key:value pair"},
					},
				},
			},
		},
		{
			name: "200 Successful",
			uuid: uuid.New(),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message":"budget found successfully",
						"data":{
							"budget":{
								"uuid":"f5fca9d0-e308-4ff2-be4e-aff22a4c2a78",
								"user_uuid":"67b14c0f-b8ea-4f0f-bf07-cadc73cd74d9",
								"organisation_uuid":null,
								"name":"test budget",
								"active":true
							}
						},
						"errors":{}
					}`,
				},
			},
			E: E{
				budget: Budget{
					UUID: uuid.MustParse("f5fca9d0-e308-4ff2-be4e-aff22a4c2a78"),
					UserUUID: uuid.MustParse("67b14c0f-b8ea-4f0f-bf07-cadc73cd74d9"),
					Name: "test budget",
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
			// add the mock exchanges to budget-micro-service
			ms.Append(tc.exchange)

			budget, e := s.GetBudget(tc.uuid)
			if tc.E.e != nil {
				if tc.E.e.Error() != e.Error() {
					t.Errorf("expected '%v' got '%v'", tc.E.e.Error(), e.Error())
				}
			} else if e != nil {
				t.Errorf("unexpected error: %s", e.Error())
			}

			if tc.E.budget != (Budget{}) {
				if budget.UUID != tc.E.budget.UUID {
					t.Errorf("expected '%v' got '%v'", tc.E.budget.UUID, budget.UUID)
				}
				if budget.Name != tc.E.budget.Name {
					t.Errorf("expected '%v' got '%v'", tc.E.budget.Name, budget.Name)
				}
				if budget.Active != tc.E.budget.Active {
					t.Errorf("expected '%v' got '%v'", tc.E.budget.Active, budget.Active)
				}
				if budget.UserUUID != tc.E.budget.UserUUID {
					t.Errorf("expected '%v' got '%v'", tc.E.budget.UserUUID, budget.UserUUID)
				}
				if budget.OrganisationUUID != tc.E.budget.OrganisationUUID {
					t.Errorf("expected '%v' got '%v'", tc.E.budget.OrganisationUUID, budget.OrganisationUUID)
				}
			}
		})
	}
}