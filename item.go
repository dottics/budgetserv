package budget

import (
	"fmt"
	"github.com/google/uuid"
)

// GetItems retrieves all items from a group based on the group's UUID.
func (s *Service) GetItems(GroupUUID uuid.UUID) (Items, error) {
	s.URL.Path = fmt.Sprintf("/group/%s/items", GroupUUID)

	res, e := s.DoRequest("GET", s.URL, nil, nil, nil)
	if e != nil {
		return Items{}, e
	}

	type data struct {
		Items Items `json:"items"`
	}

	resp := struct {
		Message string `json:"message"`
		Data    data   `json:"data"`
	}{}

	err := marshalResponse(200, res, &resp)
	if err != nil {
		return Items{}, err
	}

	return resp.Data.Items, nil
}

// CreateItem creates a new item in a group based on the ItemPayload.
func (s *Service) CreateItem(item ItemPayload) (Item, error) {
	s.URL.Path = "/item/"
	p, err := marshalReader(item)
	if err != nil {
		return Item{}, err
	}

	res, e := s.DoRequest("POST", s.URL, nil, nil, p)
	if e != nil {
		return Item{}, e
	}

	type data struct {
		Item Item `json:"item"`
	}

	resp := struct {
		Message string `json:"message"`
		Data    data   `json:"data"`
	}{}

	err = marshalResponse(201, res, &resp)
	if err != nil {
		return Item{}, err
	}

	return resp.Data.Item, nil
}
