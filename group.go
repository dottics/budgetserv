package budget

import (
	"fmt"
	"github.com/dottics/dutil"
	"github.com/google/uuid"
	"github.com/johannesscr/micro/msp"
)

// GetGroups retrieves all the groups related to a specific budget. The budget
// is identified using the uuid parameter.
//
// uuid (uuid.UUID) for the budget.
func (s *Service) GetGroups(BudgetUUID uuid.UUID) (Groups, dutil.Error) {
	s.URL.Path = fmt.Sprintf("/budget/%s/groups", BudgetUUID.String())

	type data struct {
		Groups Groups `json:"groups"`
	}

	resp := struct {
		Message string `json:"message"`
		Data    data   `json:"data"`
		Errors  map[string][]string
	}{}

	res, e := s.DoRequest("GET", s.URL, nil, nil, nil)
	if e != nil {
		return nil, e
	}
	_, e = msp.Decode(res, &resp)
	if e != nil {
		return nil, e
	}

	if res.StatusCode != 200 {
		e := &dutil.Err{
			Status: res.StatusCode,
			Errors: resp.Errors,
		}
		return nil, e
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
