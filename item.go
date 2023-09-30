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
