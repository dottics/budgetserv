package budget

import (
	"fmt"
	"github.com/dottics/dutil"
	"github.com/google/uuid"
	"github.com/johannesscr/micro/msp"
)

// GetItems retrieves all the group's items of the specific group. The group
// is identified using the uuid parameter.
//
// uuid (uuid.UUID) for the group
func (s *Service) GetItems(GroupUUID uuid.UUID) (Items, dutil.Error) {
	s.URL.Path = fmt.Sprintf("/group/%s/items", GroupUUID.String())

	type data struct {
		Items Items `json:"items"`
	}
	resp := struct {
		Message string              `json:"message"`
		Data    data                `json:"data"`
		Errors  map[string][]string `json:"errors"`
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

	return resp.Data.Items, nil
}
