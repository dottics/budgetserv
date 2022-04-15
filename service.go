package budget

import (
	"encoding/json"
	"github.com/dottics/dutil"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Service struct {
	Header http.Header
	URL    url.URL
}

func NewService(token string) *Service {
	s := &Service{
		URL: url.URL{
			Scheme: os.Getenv("BUDGET_SERVICE_SCHEME"),
			Host:   os.Getenv("BUDGET_SERVICE_HOST"),
		},
		Header: make(http.Header),
	}
	// default budget-micro-service headers
	(*s).Header.Set("Content-Type", "application/json")
	(*s).Header.Set("X-User-Token", token)

	return s
}

// SetURL sets the URL for the budget-micro-service to point to the service.
//
// SetURL is also the function which makes the service a mock service
// interface.
func (s *Service) SetURL(scheme string, host string) {
	s.URL.Scheme = scheme
	s.URL.Host = host
}

// SetEnv is used for testing, when the dynamic micro-service is created
// then SetEnv is used to dynamically set the env vars for the temporary
// service.
func (s *Service) SetEnv() error {
	err := os.Setenv("BUDGET_SERVICE_SCHEME", s.URL.Scheme)
	if err != nil {
		return err
	}
	err = os.Setenv("BUDGET_SERVICE_HOST", s.URL.Host)
	if err != nil {
		return err
	}
	return nil
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
