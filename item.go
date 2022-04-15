package budget

import (
	"github.com/dottics/dutil"
	"github.com/google/uuid"
	"net/url"
)

// GetItems retrieves all the group's items of the specific group. The group
// is identified using the uuid parameter.
//
// uuid (uuid.UUID) for the group
func (s *Service) GetItems(UUID uuid.UUID) (Items, dutil.Error) {
	s.URL.Path = "/budget/group/-/item"
	q := url.Values{
		"uuid": {UUID.String()},
	}
	s.URL.RawQuery = q.Encode()

	type data struct {
		Items Items `json:"items"`
	}
	resp := struct{
		Message string `json:"message"`
		Data data `json:"data"`
		Errors map[string][]string `json:"errors"`
	}{}

	res, e := s.newRequest("GET", s.URL.String(), nil, nil)
	if e != nil {
		return nil, e
	}
	_, e = s.decode(res, &resp)
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
