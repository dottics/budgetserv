package budget

import (
	"errors"
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
				t.Errorf("expected '%+v' got '%+v'", tc.E.budget, budget)
			}
		})
	}
}
