package budget

import (
	"errors"
	"github.com/google/uuid"
	"github.com/johannesscr/micro/microtest"
	"testing"
)

func TestService_GetBudgetCategories(t *testing.T) {
	tt := []struct {
		name       string
		groupUUID  uuid.UUID
		exchange   *microtest.Exchange
		uri        string
		categories []Category
		e          error
	}{
		{
			name:      "403 Permission Required",
			groupUUID: uuid.MustParse("ae9f5130-81fe-4526-9573-f7e892cc2e01"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   noPermission,
				},
			},
			uri:        "/category/?budget_uuid=ae9f5130-81fe-4526-9573-f7e892cc2e01",
			categories: []Category{},
			e:          errors.New("no permission"),
		},
		{
			name:      "200 Successful",
			groupUUID: uuid.MustParse("3df56eb1-b90f-4f1a-a734-9cc0e75f89ae"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   string(responseCategories),
				},
			},
			uri:        "/category/?budget_uuid=3df56eb1-b90f-4f1a-a734-9cc0e75f89ae",
			categories: testCategories,
			e:          nil,
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add the budget-micro-service exchange
			ms.Append(tc.exchange)

			categories, e := s.GetBudgetCategories(tc.groupUUID)
			// test the error response
			if NotEqualError(tc.e, e) {
				t.Errorf("expected error '%v' got '%v'", tc.e, e)
			}
			// test the exchange request uri
			if tc.uri != tc.exchange.Request.RequestURI {
				t.Errorf("expected uri '%v' got '%v'", tc.uri, tc.exchange.Request.RequestURI)
			}
			// test the response
			for i, c := range categories {
				if c != tc.categories[i] {
					t.Errorf("expected category '%v' got '%v'", tc.categories[i], c)
				}
			}
		})
	}
}

func TestService_CreateCategory(t *testing.T) {
	tt := []struct {
		name     string
		payload  CategoryCreatePayload
		exchange *microtest.Exchange
		uri      string
		category Category
		e        error
	}{
		{
			name: "403 Permission Required",
			payload: CategoryCreatePayload{
				BudgetUUID: uuid.MustParse("40355dba-0923-43a6-83d5-c9b6680edd2e"),
				Name:       "test category",
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   noPermission,
				},
			},
			uri:      "/category/",
			category: Category{},
			e:        errors.New("no permission"),
		},
		{
			name: "200 Successful",
			payload: CategoryCreatePayload{
				BudgetUUID: uuid.MustParse("3df56eb1-b90f-4f1a-a734-9cc0e75f89ae"),
				Name:       "test category",
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 201,
					Body:   string(responseCategory),
				},
			},
			uri:      "/category/",
			category: testCategory,
			e:        nil,
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add the budget-micro-service exchange
			ms.Append(tc.exchange)

			category, e := s.CreateCategory(tc.payload)
			// test the error response
			if NotEqualError(tc.e, e) {
				t.Errorf("expected error '%v' got '%v'", tc.e, e)
			}
			// test the exchange request uri
			if tc.uri != tc.exchange.Request.RequestURI {
				t.Errorf("expected uri '%v' got '%v'", tc.uri, tc.exchange.Request.RequestURI)
			}
			// test the response
			if category != tc.category {
				t.Errorf("expected category '%v' got '%v'", tc.category, category)
			}
		})
	}
}

func TestService_UpdateCategory(t *testing.T) {
	tt := []struct {
		name     string
		payload  CategoryUpdatePayload
		exchange *microtest.Exchange
		uri      string
		category Category
		e        error
	}{
		{
			name: "403 Permission Required",
			payload: CategoryUpdatePayload{
				UUID: uuid.MustParse("40355dba-0923-43a6-83d5-c9b6680edd2e"),
				Name: "test category",
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   noPermission,
				},
			},
			uri:      "/category/40355dba-0923-43a6-83d5-c9b6680edd2e",
			category: Category{},
			e:        errors.New("no permission"),
		},
		{
			name: "200 Successful",
			payload: CategoryUpdatePayload{
				UUID: uuid.MustParse("3df56eb1-b90f-4f1a-a734-9cc0e75f89ae"),
				Name: "test category",
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   string(responseCategory),
				},
			},
			uri:      "/category/3df56eb1-b90f-4f1a-a734-9cc0e75f89ae",
			category: testCategory,
			e:        nil,
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add the budget-micro-service exchange
			ms.Append(tc.exchange)

			category, e := s.UpdateCategory(tc.payload)
			// test the error response
			if NotEqualError(tc.e, e) {
				t.Errorf("expected error '%v' got '%v'", tc.e, e)
			}
			// test the exchange request uri
			if tc.uri != tc.exchange.Request.RequestURI {
				t.Errorf("expected uri '%v' got '%v'", tc.uri, tc.exchange.Request.RequestURI)
			}
			// test the response
			if category != tc.category {
				t.Errorf("expected category '%v' got '%v'", tc.category, category)
			}
		})
	}
}

func TestService_DeleteCategory(t *testing.T) {
	tt := []struct {
		name     string
		UUID     uuid.UUID
		exchange *microtest.Exchange
		uri      string
		e        error
	}{
		{
			name: "403 Permission Required",
			UUID: uuid.MustParse("40355dba-0923-43a6-83d5-c9b6680edd2e"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   noPermission,
				},
			},
			uri: "/category/40355dba-0923-43a6-83d5-c9b6680edd2e",
			e:   errors.New("no permission"),
		},
		{
			name: "200 Successful",
			UUID: uuid.MustParse("3df56eb1-b90f-4f1a-a734-9cc0e75f89ae"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message":"category deleted successfully"
					}`,
				},
			},
			uri: "/category/3df56eb1-b90f-4f1a-a734-9cc0e75f89ae",
			e:   nil,
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add the budget-micro-service exchange
			ms.Append(tc.exchange)

			e := s.DeleteCategory(tc.UUID)
			// test the error response
			if NotEqualError(tc.e, e) {
				t.Errorf("expected error '%v' got '%v'", tc.e, e)
			}
			// test the exchange request uri
			if tc.uri != tc.exchange.Request.RequestURI {
				t.Errorf("expected uri '%v' got '%v'", tc.uri, tc.exchange.Request.RequestURI)
			}
		})
	}
}
