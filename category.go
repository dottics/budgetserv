package budget

import (
	"fmt"
	"github.com/google/uuid"
	"net/url"
)

// GetBudgetCategories retrieves all categories from a budget based on the
// budget's UUID.
func (s *Service) GetBudgetCategories(BudgetUUID uuid.UUID) ([]Category, error) {
	// set URL Path
	s.URL.Path = "/category/"
	// remove any previous query parameters
	s.Values.Del("budget_uuid")
	qs := url.Values{
		"budget_uuid": []string{BudgetUUID.String()},
	}

	// make request
	res, e := s.DoRequest("GET", s.URL, qs, nil, nil)
	if e != nil {
		return []Category{}, e
	}

	// define response structure
	type data struct {
		Categories []Category `json:"categories"`
	}
	resp := struct {
		Message string `json:"message"`
		Data    data   `json:"data"`
	}{}

	// unmarshal response
	err := marshalResponse(200, res, &resp)
	if err != nil {
		return []Category{}, err
	}

	return resp.Data.Categories, nil
}

// CreateCategory creates a new category in a budget based on the
// CategoryCreatePayload.
func (s *Service) CreateCategory(payload CategoryCreatePayload) (Category, error) {
	// set the URL Path
	s.URL.Path = "/category/"
	// marshal payload
	p, err := marshalReader(payload)
	if err != nil {
		return Category{}, err
	}
	// make the request
	res, e := s.DoRequest("POST", s.URL, nil, nil, p)
	if e != nil {
		return Category{}, e
	}
	// define response structure
	type data struct {
		Category Category `json:"category"`
	}
	resp := struct {
		Message string `json:"message"`
		Data    data   `json:"data"`
	}{}
	// unmarshal response
	err = marshalResponse(201, res, &resp)
	if err != nil {
		return Category{}, err
	}
	return resp.Data.Category, nil
}

// UpdateCategory makes the request to the budget-microservice to update a
// category.
func (s *Service) UpdateCategory(payload CategoryUpdatePayload) (Category, error) {
	// set the URL Path
	s.URL.Path = fmt.Sprintf("/category/%s", payload.UUID.String())
	// marshal the payload
	p, err := marshalReader(payload)
	if err != nil {
		return Category{}, err
	}
	// make the request
	res, e := s.DoRequest("PUT", s.URL, nil, nil, p)
	if e != nil {
		return Category{}, e
	}
	// define response structure
	type data struct {
		Category Category `json:"category"`
	}
	resp := struct {
		Message string `json:"message"`
		Data    data   `json:"data"`
	}{}
	// unmarshal response
	err = marshalResponse(200, res, &resp)
	if err != nil {
		return Category{}, err
	}
	return resp.Data.Category, nil
}

// DeleteCategory makes the request to the budget-microservice to delete a
// category.
func (s *Service) DeleteCategory(CategoryUUID uuid.UUID) error {
	// set the URL Path
	s.URL.Path = fmt.Sprintf("/category/%s", CategoryUUID.String())
	// make the request
	res, e := s.DoRequest("DELETE", s.URL, nil, nil, nil)
	if e != nil {
		return e
	}
	// define response structure
	resp := struct {
		Message string `json:"message"`
	}{}
	// unmarshal response
	return marshalResponse(200, res, &resp)
}
