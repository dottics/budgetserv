package budget

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/johannesscr/micro/microtest"
	"testing"
)

func TestService_GetBudgets(t *testing.T) {
	type E struct {
		status int
		len    int
		e      error
	}

	tt := []struct {
		name     string
		exchange *microtest.Exchange
		E        E
	}{
		{
			name: "403 Forbidden",
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   noPermission,
				},
			},
			E: E{
				status: 403,
				len:    0,
				e:      errors.New("no permission"),
			},
		},
		{
			name: "200 Successful",
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   string(responseBudgets),
				},
			},
			E: E{
				status: 200,
				len:    2,
				e:      nil,
			},
		},
	}

	s := NewService(Config{})
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
		e      error
	}

	tt := []struct {
		name     string
		uuid     uuid.UUID
		exchange *microtest.Exchange
		E        E
	}{
		{
			name: "403 Forbidden",
			uuid: uuid.New(),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   noPermission,
				},
			},
			E: E{
				budget: Budget{},
				e:      errors.New("no permission"),
			},
		},
		{
			name: "404 Not Found",
			uuid: uuid.New(),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 404,
					Body:   errorResponseDetail(`"budget not found"`),
				},
			},
			E: E{
				budget: Budget{},
				e:      errors.New("budget not found"),
			},
		},
		{
			name: "500 Unmarshal Error",
			uuid: uuid.New(),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   `{"message":"error""}`,
				},
			},
			E: E{
				budget: Budget{},
				e:      errors.New(`invalid character '"' after object key:value pair`),
			},
		},
		{
			name: "200 Successful",
			uuid: uuid.New(),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   string(responseBudget),
				},
			},
			E: E{
				budget: testBudget,
				e:      nil,
			},
		},
		{
			name: "200 Successful with Groups",
			uuid: uuid.New(),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   string(responseBudgetWithGroups),
				},
			},
			E: E{
				budget: testBudgetWithGroups,
				e:      nil,
			},
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add the mock exchanges to budget-micro-service
			ms.Append(tc.exchange)

			budget, e := s.GetBudget(tc.uuid)

			if NotEqualError(tc.E.e, e) {
				t.Errorf("expected '%+v' got '%+v'", tc.E.e, e)
			}

			if !EqualBudget(tc.E.budget, budget) {
				t.Errorf("expected\n'%+v'\ngot\n'%+v'", tc.E.budget, budget)
			}
		})
	}
}

func TestService_CreateBudget(t *testing.T) {
	tests := []struct {
		name      string
		exchange  *microtest.Exchange
		budget    BudgetCreatePayload
		e         error
		resBudget Budget
	}{
		{
			name: "403 permission required",
			budget: BudgetCreatePayload{
				Name:        "new budget",
				Description: "new budget desc",
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   noPermission,
				},
			},
			resBudget: Budget{},
			e:         errors.New("no permission"),
		},
		{
			name: "200 successful",
			budget: BudgetCreatePayload{
				Name:        "new budget",
				Description: "new budget desc",
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 201,
					Body: `{
						"message": "budget created",
						"data": {
							"budget": {
								"uuid": "8ace0389-f7a9-4e54-b4c8-83c2e88b1a23",
								"name": "new budget",
								"description": "new budget desc"
						   }
						}
					}`,
				},
			},
			resBudget: Budget{
				UUID:        uuid.MustParse("8ace0389-f7a9-4e54-b4c8-83c2e88b1a23"),
				Name:        "new budget",
				Description: "new budget desc",
				Groups:      Groups{},
			},
			e: nil,
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for i, tc := range tests {
		name := fmt.Sprintf("%d %s", i, tc.name)
		t.Run(name, func(t *testing.T) {
			// add the service response to the mocked service
			ms.Append(tc.exchange)

			budget, err := s.CreateBudget(tc.budget)
			if NotEqualError(tc.e, err) {
				t.Errorf("expected error '%v' got '%v'", tc.e, err)
			}

			if !EqualBudget(tc.resBudget, budget) {
				t.Errorf("expected budget\n'%+v'\ngot\n'%+v'", tc.resBudget, budget)
			}
		})
	}
}

func TestService_UpdateBudget(t *testing.T) {
	tests := []struct {
		name     string
		payload  BudgetUpdatePayload
		exchange *microtest.Exchange
		uri      string
		budget   Budget
		e        error
	}{
		{
			name: "403 permission required",
			payload: BudgetUpdatePayload{
				UUID:        uuid.MustParse("8ace0389-f7a9-4e54-b4c8-83c2e88b1a23"),
				Name:        "new budget",
				Description: "new budget desc",
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   noPermission,
				},
			},
			uri:    "/budget/8ace0389-f7a9-4e54-b4c8-83c2e88b1a23",
			budget: Budget{},
			e:      errors.New("no permission"),
		},
		{
			name: "200 successful",
			payload: BudgetUpdatePayload{
				UUID:        uuid.MustParse("c9521d38-c6cd-401e-b968-832457a31217"),
				Name:        "new budget",
				Description: "new budget desc",
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message": "budget updated",
						"data": {
							"budget": {
								"uuid": "c9521d38-c6cd-401e-b968-832457a31217",
								"name": "new budget",
								"description": "new budget desc",
								"groups": []
							}
						}
					}`,
				},
			},
			uri: "/budget/c9521d38-c6cd-401e-b968-832457a31217",
			budget: Budget{
				UUID:        uuid.MustParse("c9521d38-c6cd-401e-b968-832457a31217"),
				Name:        "new budget",
				Description: "new budget desc",
				Groups:      Groups{},
			},
			e: nil,
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for i, tc := range tests {
		name := fmt.Sprintf("%d %s", i, tc.name)
		t.Run(name, func(t *testing.T) {
			ms.Append(tc.exchange)

			budget, err := s.UpdateBudget(tc.payload)
			if NotEqualError(tc.e, err) {
				t.Errorf("expected error '%v' got '%v'", tc.e, err)
			}

			if !EqualBudget(tc.budget, budget) {
				t.Errorf("expected budget\n'%+v'\ngot\n'%+v'", tc.budget, budget)
			}

			if tc.uri != tc.exchange.Request.RequestURI {
				t.Errorf("expected uri '%s' got '%s'", tc.uri, tc.exchange.Request.RequestURI)
			}
		})
	}
}
