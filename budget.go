package budget

import (
	"github.com/dottics/dutil"
	"github.com/google/uuid"
	"net/url"
)

func (s *Service) GetBudgets() (Budgets, dutil.Error) {
	s.URL.Path = "/budget"

	type data struct {
		Budgets Budgets `json:"budgets"`
	}

	resp := struct {
		Message string              `json:"message"`
		Data    data                `json:"data"`
		Errors  map[string][]string `json:"errors"`
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

	return resp.Data.Budgets, nil
}

func (s *Service) GetBudget(UUID uuid.UUID) (Budget, dutil.Error) {
	s.URL.Path = "/budget/-"
	q := url.Values{}
	q.Add("uuid", UUID.String())
	s.URL.RawQuery = q.Encode()

	type data struct {
		Budget Budget `json:"budget"`
	}

	resp := struct {
		Message string              `json:"message"`
		Data    data                `json:"data"`
		Errors  map[string][]string `json:"errors"`
	}{}

	res, e := s.newRequest("GET", s.URL.String(), nil, nil)
	if e != nil {
		return Budget{}, e
	}
	_, e = s.decode(res, &resp)
	if e != nil {
		return Budget{}, e
	}

	if res.StatusCode != 200 {
		e := &dutil.Err{
			Status: res.StatusCode,
			Errors: resp.Errors,
		}
		return Budget{}, e
	}

	return resp.Data.Budget, nil
}
