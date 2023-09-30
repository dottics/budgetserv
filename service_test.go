package budget

import (
	"github.com/dottics/dutil"
	"github.com/johannesscr/micro/microtest"
	"os"
	"testing"
)

func TestNewService(t *testing.T) {
	type E struct {
		scheme string
		host   string
		token  string
	}
	tt := []struct {
		name            string
		budgetSchemeEnv string
		budgetHostEnv   string
		token           string
		E               E
	}{
		{
			name: "default",
			E: E{
				scheme: "",
				host:   "",
				token:  "",
			},
		},
		{
			name:            "env vars",
			budgetSchemeEnv: "https",
			budgetHostEnv:   "budget.ms.dottics.com",
			E: E{
				scheme: "https",
				host:   "budget.ms.dottics.com",
				token:  "",
			},
		},
		{
			name:  "token",
			token: "my-test-token",
			E: E{
				scheme: "",
				host:   "",
				token:  "my-test-token",
			},
		},
		{
			name:            "token and env vars",
			budgetSchemeEnv: "https",
			budgetHostEnv:   "budget.ms.dottics.com",
			token:           "my-test-token",
			E: E{
				scheme: "https",
				host:   "budget.ms.dottics.com",
				token:  "my-test-token",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := os.Setenv("BUDGET_SERVICE_SCHEME", tc.budgetSchemeEnv)
			err = os.Setenv("BUDGET_SERVICE_HOST", tc.budgetHostEnv)
			if err != nil {
				t.Errorf("unexpected error before: %v", err)
			}

			s := NewService(Config{UserToken: tc.token})
			xut := s.Header.Get("X-User-Token")
			if tc.E.token != xut {
				t.Errorf("expected '%v' got '%v'", tc.E.token, xut)
			}
			if tc.E.scheme != s.URL.Scheme {
				t.Errorf("expected '%v' got '%v'", tc.E.scheme, s.URL.Scheme)
			}
			if tc.E.host != s.URL.Host {
				t.Errorf("expected '%v' got '%v'", tc.E.host, s.URL.Host)
			}

			// reset to blank
			err = os.Setenv("BUDGET_SERVICE_SCHEME", "")
			err = os.Setenv("BUDGET_SERVICE_HOST", "")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestService_SetURL(t *testing.T) {
	s := NewService(Config{})
	s.SetURL("http", "budget.ms.test.dottics.com")
	if s.URL.Scheme != "http" {
		t.Errorf("expected '%v' got '%v'", "http", s.URL.Scheme)
	}
	if s.URL.Host != "budget.ms.test.dottics.com" {
		t.Errorf("expected '%v' got '%v'", "budget.ms.test.dottics.com", s.URL.Host)
	}
}

func TestService_SetEnv(t *testing.T) {
	s := NewService(Config{})
	s.URL.Scheme = "https"
	s.URL.Host = "test.host.com"
	err := s.SetEnv()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	x1 := os.Getenv("BUDGET_SERVICE_SCHEME")
	x2 := os.Getenv("BUDGET_SERVICE_HOST")
	if x1 != "https" {
		t.Errorf("expected '%v' got '%v'", "https", x1)
	}
	if x2 != "test.host.com" {
		t.Errorf("expected '%v' got '%v'", "test.host.com", x2)
	}
}

func TestGetHome(t *testing.T) {
	type E struct {
		alive bool
		e     dutil.Err
	}
	tt := []struct {
		name     string
		exchange *microtest.Exchange
		E        E
	}{
		{
			name: "decode error",
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   `{"message":"Welcome to the budget micro-service","data":{},"errors":{"internal_server_error":"server down for some reason"]}}`,
				},
			},
			E: E{
				alive: false,
				e: dutil.Err{
					Status: 500,
					Errors: map[string][]string{
						"unmarshal": {"invalid character ']' after object key:value pair"},
					},
				},
			},
		},
		{
			name: "500 internal server error",
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 500,
					Body:   `{"message":"Welcome to the budget micro-service","data":{},"errors":{"internal_server_error":["server down for some reason"]}}`,
				},
			},
			E: E{
				alive: false,
				e: dutil.Err{
					Status: 500,
					Errors: map[string][]string{
						"internal_server_error": {"server down for some reason"},
					},
				},
			},
		},
		{
			name: "200 server alive",
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   `{"message":"Welcome to the budget micro-service","data":{"alive":true},"errors":{}}`,
				},
			},
			E: E{
				alive: true,
				e:     dutil.Err{},
			},
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// append the exchange for the test
			ms.Append(tc.exchange)

			alive, e := s.HealthCheck()
			if tc.E.alive != alive {
				t.Errorf("expected '%v' got '%v'", tc.E.alive, alive)
			}
			if tc.E.e.Status != 0 {
				if tc.E.e.Error() != e.Error() {
					t.Errorf("expected '%v' got '%v'", tc.E.e.Error(), e.Error())
				}
			}
		})
	}
}
