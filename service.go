package budget

import (
	"github.com/johannesscr/micro/msp"
	"net/http"
	"net/url"
)

// Service is the budget-microservice which is used to make requests to the
// budget-service.
type Service struct {
	msp.Service
}

// Config is the configuration for the budget-microservice.
type Config struct {
	UserToken string
	APIKey    string
	Header    http.Header
	Values    url.Values
}

func NewService(config Config) *Service {
	s := &Service{
		Service: *msp.NewService(msp.Config{
			Name:      "budget",
			UserToken: config.UserToken,
			APIKey:    config.APIKey,
			Header:    config.Header,
			Values:    config.Values,
		}),
	}
	return s
}
