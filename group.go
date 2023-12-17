package budget

import (
	"fmt"
	"github.com/google/uuid"
)

// GetGroups retrieves all the groups related to a specific budget. The budget
// is identified using the uuid parameter.
//
// uuid (uuid.UUID) for the budget.
func (s *Service) GetGroups(BudgetUUID uuid.UUID) (Groups, error) {
	s.URL.Path = fmt.Sprintf("/budget/%s/groups", BudgetUUID.String())

	res, e := s.DoRequest("GET", s.URL, nil, nil, nil)
	if e != nil {
		return Groups{}, e
	}

	type data struct {
		Groups Groups `json:"groups"`
	}

	resp := struct {
		Message string `json:"message"`
		Data    data   `json:"data"`
	}{}

	err := marshalResponse(200, res, &resp)
	if err != nil {
		return Groups{}, err
	}
	return resp.Data.Groups, nil
}

// GetGroup retrieves a specific group. The group is identified using the uuid
// parameter.
func (s *Service) GetGroup(UUID uuid.UUID) (Group, error) {
	s.URL.Path = fmt.Sprintf("/group/%s", UUID.String())

	res, e := s.DoRequest("GET", s.URL, nil, nil, nil)
	if e != nil {
		return Group{}, e
	}

	type data struct {
		Group Group `json:"group"`
	}

	resp := struct {
		Message string `json:"message"`
		Data    data   `json:"data"`
	}{}

	err := marshalResponse(200, res, &resp)
	if err != nil {
		return Group{}, err
	}
	return resp.Data.Group, nil
}

// CreateGroup creates a new group based on the GroupCreatePayload which is
// passed as a parameter.
func (s *Service) CreateGroup(payload GroupCreatePayload) (Group, error) {
	s.URL.Path = "/group/"
	p, err := marshalReader(payload)
	if err != nil {
		return Group{}, err
	}

	res, e := s.DoRequest("POST", s.URL, nil, nil, p)
	if e != nil {
		return Group{}, e
	}

	type data struct {
		Group Group `json:"group"`
	}

	resp := struct {
		Message string `json:"message"`
		Data    data   `json:"data"`
	}{}

	err = marshalResponse(201, res, &resp)
	if err != nil {
		return Group{}, err
	}
	return resp.Data.Group, nil
}

// UpdateGroup updates a specific group based on the GroupUpdatePayload which
// is passed as a parameter.
func (s *Service) UpdateGroup(payload GroupUpdatePayload) (Group, error) {
	s.URL.Path = fmt.Sprintf("/group/%s", payload.UUID.String())
	p, err := marshalReader(payload)
	if err != nil {
		return Group{}, err
	}

	res, e := s.DoRequest("PUT", s.URL, nil, nil, p)
	if e != nil {
		return Group{}, e
	}

	type data struct {
		Group Group `json:"group"`
	}

	resp := struct {
		Message string `json:"message"`
		Data    data   `json:"data"`
	}{}

	err = marshalResponse(200, res, &resp)
	if err != nil {
		return Group{}, err
	}
	return resp.Data.Group, nil
}

// DeleteGroup deletes a specific group. The group is identified using the uuid
// parameter. The methods also deletes all subgroups, items and events related
// to the group.
func (s *Service) DeleteGroup(UUID uuid.UUID) ([]Category, error) {
	s.URL.Path = fmt.Sprintf("/group/%s", UUID.String())

	res, e := s.DoRequest("DELETE", s.URL, nil, nil, nil)
	if e != nil {
		return []Category{}, e
	}

	type data struct {
		RemovedCategories []Category `json:"removed_categories"`
	}
	resp := struct {
		Message string `json:"message"`
		Data    data   `json:"data"`
	}{}

	err := marshalResponse(200, res, &resp)
	if err != nil {
		return []Category{}, err
	}
	return resp.Data.RemovedCategories, err
}
