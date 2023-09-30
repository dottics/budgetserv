package budget

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/johannesscr/micro/microtest"
	"testing"
)

func TestService_GetItems(t *testing.T) {
	tests := []struct {
		name      string
		groupUUID uuid.UUID
		exchange  *microtest.Exchange
		uri       string
		items     Items
		e         error
	}{
		{
			name:      "403 permission required",
			groupUUID: uuid.MustParse("40355dba-0923-43a6-83d5-c9b6680edd2e"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   noPermission,
				},
			},
			uri:   "/group/40355dba-0923-43a6-83d5-c9b6680edd2e/items",
			items: Items{},
			e:     errors.New("no permission"),
		},
		{
			name:      "200 successful",
			groupUUID: uuid.MustParse("f27ef50d-f10f-4ff8-b65a-d64b1ebb83c8"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   string(responseGroupItems),
				},
			},
			uri:   "/group/f27ef50d-f10f-4ff8-b65a-d64b1ebb83c8/items",
			items: testGroupItems,
			e:     nil,
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for i, tc := range tests {
		name := fmt.Sprintf("%d %s", i, tc.name)
		t.Run(name, func(t *testing.T) {
			ms.Append(tc.exchange)

			items, e := s.GetItems(tc.groupUUID)

			// ensure errors match as expected
			if NotEqualError(e, tc.e) {
				t.Errorf("expected error %v, got %v", tc.e, e)
			}

			// ensure items match as expected
			if EqualItems(items, tc.items) == false {
				t.Errorf("expected items\n%+v\ngot\n%+v", tc.items, items)
			}

			// ensure the correct uri was called
			if tc.uri != tc.exchange.Request.RequestURI {
				t.Errorf("expected uri %s, got %s", tc.uri, tc.exchange.Request.RequestURI)
			}
		})
	}
}

func TestService_CreateItem(t *testing.T) {
	tests := []struct {
		name     string
		payload  ItemPayload
		exchange *microtest.Exchange
		uri      string
		item     Item
		e        error
	}{
		{
			name: "403 permission required",
			payload: ItemPayload{
				GroupUUID: uuid.MustParse("40355dba-0923-43a6-83d5-c9b6680edd2e"),
				Name:      "test item",
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   noPermission,
				},
			},
			uri:  "/item/",
			item: Item{},
			e:    errors.New("no permission"),
		},
		{
			name: "404 group not found",
			payload: ItemPayload{
				GroupUUID: uuid.MustParse("f27ef50d-f10f-4ff8-b65a-d64b1ebb83c8"),
				Name:      "test item",
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 404,
					Body:   notFound,
				},
			},
			uri:  "/item/",
			item: Item{},
			e:    errors.New("not found"),
		},
		{
			name: "201 successful",
			payload: ItemPayload{
				GroupUUID: uuid.MustParse("f27ef50d-f10f-4ff8-b65a-d64b1ebb83c8"),
				Name:      "test item",
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 201,
					Body:   string(responseItemNew),
				},
			},
			uri:  "/item/",
			item: testItemNew,
			e:    nil,
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for i, tc := range tests {
		name := fmt.Sprintf("%d %s", i, tc.name)
		t.Run(name, func(t *testing.T) {
			ms.Append(tc.exchange)

			item, e := s.CreateItem(tc.payload)

			// ensure errors match as expected
			if NotEqualError(e, tc.e) {
				t.Errorf("expected error %v, got %v", tc.e, e)
			}

			// ensure items match as expected
			if EqualItem(item, tc.item) == false {
				t.Errorf("expected item\n%+v\ngot\n%+v", tc.item, item)
			}

			// ensure the correct uri was called
			if tc.uri != tc.exchange.Request.RequestURI {
				t.Errorf("expected uri %s, got %s", tc.uri, tc.exchange.Request.RequestURI)
			}
		})
	}
}
