package budget

import (
	"github.com/dottics/dutil"
	"github.com/google/uuid"
	"net/url"
)

// GetGroups retrieves all the groups related to a specific budget. The budget
// is identified using the uuid parameter.
//
// uuid (uuid.UUID) for the budget.
func (s * Service) GetGroups(UUID uuid.UUID) (Groups, dutil.Error) {
	s.URL.Path = "/budget/-/group"
	// uuid denotes the budget's uuid
	q := url.Values{
		"uuid": {UUID.String()},
	}
	s.URL.RawQuery = q.Encode()

	type data struct {
		Groups Groups `json:"groups"`
	}

	resp := struct {
		Message string `json:"message"`
		Data data `json:"data"`
		Errors map[string][]string
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
	return resp.Data.Groups, nil
}
