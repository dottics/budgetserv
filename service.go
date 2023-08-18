package budget

import (
	"encoding/json"
	"github.com/dottics/dutil"
	"github.com/johannesscr/micro/msp"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// Service is the budget-microservice which is used to make requests to the
// budget-service.
type Service struct {
	//Header http.Header
	//URL    url.URL
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

// NewRequest consistently maps and executes requests to the requirements
// for the service and returns the response.
func (s *Service) newRequest(method string, url string, headers map[string][]string, payload io.Reader) (*http.Response, dutil.Error) {
	client := http.Client{}
	req, _ := http.NewRequest(method, url, payload)
	// set the default headers from the service
	req.Header = s.Header
	// set/override additional header iff necessary
	for key, values := range headers {
		req.Header.Set(key, values[0])
	}
	res, err := client.Do(req)
	log.Printf("- budget-service -> [ %v  %v ] <- %d",
		req.Method, req.URL.String(), res.StatusCode)
	if err != nil {
		e := dutil.NewErr(500, "request", []string{err.Error()})
		return nil, e
	}
	return res, nil
}

// decode is a function that decodes a body into a slice of bytes and error of
// there is one. If the interface pointer value is given then unmarshal the
// response body into the value pointed to by the interface.
func (s *Service) decode(res *http.Response, v interface{}) ([]byte, dutil.Error) {
	xb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		e := dutil.NewErr(500, "read", []string{err.Error()})
		return nil, e
	}
	err = res.Body.Close()
	if err != nil {
		e := dutil.NewErr(500, "decode", []string{err.Error()})
		return nil, e
	}
	//log.Printf("BUDGET SERVICE DECODE: %s", string(xb))

	if v != nil {
		err = json.Unmarshal(xb, v)
		if err != nil {
			e := dutil.NewErr(500, "unmarshal", []string{err.Error()})
			return nil, e
		}
	}

	return xb, nil
}

// GetHome is the health-check function which makes a request to the
// budget-service to check that the service is still up and running.
func (s *Service) GetHome() (bool, dutil.Error) {
	s.URL.Path = "/"

	resp := struct {
		Message string              `json:"message"`
		Data    interface{}         `json:"data"`
		Errors  map[string][]string `json:"errors"`
	}{}

	res, e := s.newRequest("GET", s.URL.String(), nil, nil)
	if e != nil {
		return false, e
	}
	_, e = s.decode(res, &resp)
	if e != nil {
		return false, e
	}

	if res.StatusCode != 200 {
		e := &dutil.Err{
			Status: res.StatusCode,
			Errors: resp.Errors,
		}
		return false, e
	}
	return true, nil
}
