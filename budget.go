package budget

import (
	"fmt"
	"github.com/google/uuid"
)

// GetBudgets retrieves all the budgets from the budget-microservice.
func (s *Service) GetBudgets() (Budgets, error) {
	s.URL.Path = "/budget/"

	type data struct {
		Budgets Budgets `json:"budgets"`
	}

	resp := struct {
		Message string `json:"message"`
		Data    data   `json:"data"`
	}{}

	res, e := s.DoRequest("GET", s.URL, nil, nil, nil)
	if e != nil {
		return nil, e
	}

	err := marshalResponse(200, res, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data.Budgets, nil
}

// GetEntityBudgets retrieves all the budgets related to a specific entity. The
// entity is identified using the uuid parameter.
func (s *Service) GetEntityBudgets(EntityUUID uuid.UUID) (Budgets, error) {
	s.URL.Path = fmt.Sprintf("/budget/entity/%s", EntityUUID.String())

	type data struct {
		Budgets Budgets `json:"budgets"`
	}

	resp := struct {
		Message string `json:"message"`
		Data    data   `json:"data"`
	}{}

	res, e := s.DoRequest("GET", s.URL, nil, nil, nil)
	if e != nil {
		return nil, e
	}

	err := marshalResponse(200, res, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data.Budgets, nil
}

// GetBudget retrieves a specific budget's data.
func (s *Service) GetBudget(UUID uuid.UUID) (Budget, error) {
	s.URL.Path = fmt.Sprintf("/budget/%s", UUID.String())

	type data struct {
		Budget Budget `json:"budget"`
	}

	resp := struct {
		Message string `json:"message"`
		Data    data   `json:"data"`
	}{}

	res, e := s.DoRequest("GET", s.URL, nil, nil, nil)
	if e != nil {
		return Budget{}, e
	}

	err := marshalResponse(200, res, &resp)
	if err != nil {
		return Budget{}, err
	}

	return resp.Data.Budget, nil
}

// SetupBudget to setup a new budget for an entity.
func (s *Service) SetupBudget(payload BudgetSetupPayload) (Budget, error) {
	s.URL.Path = "/budget/setup"
	p, err := marshalReader(payload)
	if err != nil {
		return Budget{}, nil
	}

	res, e := s.DoRequest("POST", s.URL, nil, nil, p)
	if e != nil {
		return Budget{}, nil
	}

	resp := struct {
		Message string `json:"message"`
		Data    struct {
			Budget Budget `json:"budget"`
		} `json:"data"`
	}{}

	err = marshalResponse(201, res, &resp)
	if err != nil {
		return Budget{}, err
	}

	return resp.Data.Budget, nil
}

// CreateBudget to create a new budget for an entity.
func (s *Service) CreateBudget(budget BudgetCreatePayload) (Budget, error) {
	s.URL.Path = "/budget/"
	p, err := marshalReader(budget)
	if err != nil {
		return Budget{}, nil
	}

	res, e := s.DoRequest("POST", s.URL, nil, nil, p)
	if e != nil {
		return Budget{}, nil
	}

	type data struct {
		Budget Budget `json:"budget"`
	}

	resp := struct {
		Message string `json:"message"`
		Data    data   `json:"data"`
	}{}

	err = marshalResponse(201, res, &resp)
	if err != nil {
		return Budget{}, err
	}

	return resp.Data.Budget, nil
}

// UpdateBudget updates a budget's information
func (s *Service) UpdateBudget(budget BudgetUpdatePayload) (Budget, error) {
	s.URL.Path = fmt.Sprintf("/budget/%s", budget.UUID.String())
	p, err := marshalReader(budget)
	if err != nil {
		return Budget{}, nil
	}

	res, e := s.DoRequest("PUT", s.URL, nil, nil, p)
	if e != nil {
		return Budget{}, nil
	}

	type data struct {
		Budget Budget `json:"budget"`
	}

	resp := struct {
		Data data `json:"data"`
	}{}

	err = marshalResponse(200, res, &resp)
	if err != nil {
		return Budget{}, err
	}

	return resp.Data.Budget, nil
}
