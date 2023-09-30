package budget

import (
	"errors"
	"github.com/dottics/dutil"
	"github.com/google/uuid"
	"github.com/johannesscr/micro/microtest"
	"testing"
)

func TestService_GetGroups(t *testing.T) {
	type E struct {
		groups   Groups
		e        dutil.Error
		exReqURI string
	}
	tt := []struct {
		name       string
		budgetUUID uuid.UUID
		exchange   *microtest.Exchange
		E          E
	}{
		{
			name:       "403 Permission Required",
			budgetUUID: uuid.MustParse("2520f807-915e-41f6-9557-84500e1aebcc"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body: `{
						"message":"Forbidden: unable to process request",
						"data":{},
						"errors":{
							"auth":["Please ensure you have permission"]
						}
					}`,
				},
			},
			E: E{
				exReqURI: "/budget/2520f807-915e-41f6-9557-84500e1aebcc/groups",
				e: &dutil.Err{
					Status: 403,
					Errors: map[string][]string{
						"auth": {"Please ensure you have permission"},
					},
				},
				groups: Groups{},
			},
		},
		{
			name:       "404 Not Found",
			budgetUUID: uuid.MustParse("7cb47f06-0d96-494b-a847-a472e2c04d9d"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 404,
					Body: `{
						"message":"NotFound: unable to process request",
						"data":{},
						"errors":{
							"groups":["not found"]
						}
					}`,
				},
			},
			E: E{
				exReqURI: "/budget/7cb47f06-0d96-494b-a847-a472e2c04d9d/groups",
				e: &dutil.Err{
					Status: 404,
					Errors: map[string][]string{
						"groups": {"not found"},
					},
				},
				groups: Groups{},
			},
		},
		{
			name:       "500 Unmarshal Error",
			budgetUUID: uuid.MustParse("1a252f35-c84f-4937-9c5b-f5deb19b5b10"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 500,
					Body: `{
						"message":"groups found successfully",
						"data":{
							"groups":[
								{
									"uuid":"6be3df72-da3d-4a8c-bef6-d0b57120b80a",
									"name":"income",
									"active":true,
									"sub_groups":[
										{
											"uuid":"ef8e63d5-b26c-4824-8a1a-a729bb8574d3",
											"name":"salary",
											"active":true,
											"sub_groups":[],
										}
									]
								},
								{
									"uuid":"52f2c725-2cdc-401a-abdd-66db5fd06789",
									"name":"Investments",
									"active":true,
									"sub_groups":[]
								}
							]
						},
						"errors":{}
					}`,
				},
			},
			E: E{
				exReqURI: "/budget/1a252f35-c84f-4937-9c5b-f5deb19b5b10/groups",
				e: &dutil.Err{
					Status: 500,
					Errors: map[string][]string{
						"unmarshal": {"invalid character '}' looking for beginning of object key string"},
					},
				},
				groups: Groups{},
			},
		},
		{
			name:       "500 Internal Server Error",
			budgetUUID: uuid.MustParse("3a8113f2-af77-4430-a59e-519a8ad0819d"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 500,
					Body: `{
						"message":"Forbidden: unable to process request",
						"data":{},
						"errors":{
							"internal_server_error":["some unexpected error"]
						}
					}`,
				},
			},
			E: E{
				exReqURI: "/budget/3a8113f2-af77-4430-a59e-519a8ad0819d/groups",
				e: &dutil.Err{
					Status: 500,
					Errors: map[string][]string{
						"internal_server_error": {"some unexpected error"},
					},
				},
				groups: Groups{},
			},
		},
		{
			name:       "200 Successful",
			budgetUUID: uuid.MustParse("b440353e-cc26-449c-a470-e0e36a2919a6"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body: `{
						"message":"groups found successfully",
						"data":{
							"groups":[
								{
									"uuid":"52f2c725-2cdc-401a-abdd-66db5fd06789",
									"name":"income",
									"active":true,
									"sub_groups":[
										{
											"uuid":"b8448a78-6417-4fe2-849c-024622bc6106",
											"name":"base salary",
											"active":true,
											"sub_groups":[]
										}
									]
								},
								{
									"uuid":"eea51d45-c9bd-45e2-bc80-010ecbb7a0d3",
									"name":"investments",
									"active":true,
									"sub_groups":[]
								},
								{
									"uuid":"6be3df72-da3d-4a8c-bef6-d0b57120b80a",
									"name":"expenses",
									"active":true,
									"sub_groups":[]
								}
							]
						},
						"errors":{}
					}`,
				},
			},
			E: E{
				exReqURI: "/budget/b440353e-cc26-449c-a470-e0e36a2919a6/groups",
				e:        nil,
				groups: Groups{
					Group{
						UUID: uuid.MustParse("52f2c725-2cdc-401a-abdd-66db5fd06789"),
						Name: "income",
						SubGroups: Groups{
							Group{
								UUID: uuid.MustParse("b8448a78-6417-4fe2-849c-024622bc6106"),
								Name: "base salary",
							},
						},
					},
					Group{
						UUID: uuid.MustParse("eea51d45-c9bd-45e2-bc80-010ecbb7a0d3"),
						Name: "investments",
					},
					Group{
						UUID: uuid.MustParse("6be3df72-da3d-4a8c-bef6-d0b57120b80a"),
						Name: "expenses",
					},
				},
			},
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add the budget-micro-service exchange
			ms.Append(tc.exchange)

			xg, e := s.GetGroups(tc.budgetUUID)
			// test the error response
			if NotEqualError(tc.E.e, e) {
				t.Errorf("expected error '%v' got '%v'", tc.E.e, e)
			}

			// test the groups structure returned
			if EqualGroups(xg, tc.E.groups) == false {
				t.Errorf("expected groups\n'%+v'\ngot\n'%+v'", tc.E.groups, xg)
			}

			// test the exchange request URI
			if tc.exchange.Request.RequestURI != tc.E.exReqURI {
				t.Errorf("expected uri '%v' got '%v'", tc.E.exReqURI, tc.exchange.Request.RequestURI)
			}
		})
	}
}

func TestService_GetGroup(t *testing.T) {
	tt := []struct {
		name      string
		groupUUID uuid.UUID
		exchange  *microtest.Exchange
		uri       string
		group     Group
		e         error
	}{
		{
			name:      "403 Permission Required",
			groupUUID: uuid.MustParse("2520f807-915e-41f6-9557-84500e1aebcc"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   noPermission,
				},
			},
			uri:   "/group/2520f807-915e-41f6-9557-84500e1aebcc",
			group: Group{},
			e:     errors.New("no permission"),
		},
		{
			name:      "404 Not Found",
			groupUUID: uuid.MustParse("7cb47f06-0d96-494b-a847-a472e2c04d9d"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 404,
					Body:   notFound,
				},
			},
			uri:   "/group/7cb47f06-0d96-494b-a847-a472e2c04d9d",
			group: Group{},
			e:     errors.New("not found"),
		},
		{
			name:      "200 Successful",
			groupUUID: uuid.MustParse("b440353e-cc26-449c-a470-e0e36a2919a6"),
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   string(responseGroup),
				},
			},
			uri:   "/group/b440353e-cc26-449c-a470-e0e36a2919a6",
			group: testGroup,
			e:     nil,
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add the budget-micro-service exchange
			ms.Append(tc.exchange)

			xg, e := s.GetGroup(tc.groupUUID)
			// test the error response
			if NotEqualError(tc.e, e) {
				t.Errorf("expected error '%v' got '%v'", tc.e, e)
			}

			// test the group structure returned
			if EqualGroup(xg, tc.group) == false {
				t.Errorf("expected group\n'%+v'\ngot\n'%+v'", tc.group, xg)
			}
			// test the exchange request URI
			if tc.exchange.Request.RequestURI != tc.uri {
				t.Errorf("expected uri '%v' got '%v'", tc.uri, tc.exchange.Request.RequestURI)
			}
		})
	}
}

func TestService_CreateGroup(t *testing.T) {
	tt := []struct {
		name     string
		payload  GroupCreatePayload
		exchange *microtest.Exchange
		uri      string
		group    Group
		e        error
	}{
		{
			name: "403 Permission Required",
			payload: GroupCreatePayload{
				Name: "test group",
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   noPermission,
				},
			},
			uri:   "/group/",
			group: Group{},
			e:     errors.New("no permission"),
		},
		{
			name: "201 Created",
			payload: GroupCreatePayload{
				Name: "test group",
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 201,
					Body:   string(responseGroupBasic),
				},
			},
			uri:   "/group/",
			group: testGroupBasic,
			e:     nil,
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add the budget-micro-service exchange
			ms.Append(tc.exchange)

			g, e := s.CreateGroup(tc.payload)
			// test the error response
			if NotEqualError(tc.e, e) {
				t.Errorf("expected error '%v' got '%v'", tc.e, e)
			}

			// test the group structure returned
			if EqualGroup(g, tc.group) == false {
				t.Errorf("expected group\n'%+v'\ngot\n'%+v'", tc.group, g)
			}
		})
	}
}

func TestService_UpdateGroup(t *testing.T) {
	tt := []struct {
		name     string
		payload  GroupUpdatePayload
		exchange *microtest.Exchange
		uri      string
		group    Group
		e        error
	}{
		{
			name: "403 Permission Required",
			payload: GroupUpdatePayload{
				UUID: testGroup.UUID,
				Name: "test group",
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 403,
					Body:   noPermission,
				},
			},
			uri:   "/group/52f2c725-2cdc-401a-abdd-66db5fd06789",
			group: Group{},
			e:     errors.New("no permission"),
		},
		{
			name: "200 Successful",
			payload: GroupUpdatePayload{
				UUID: testGroup.UUID,
				Name: "test group",
			},
			exchange: &microtest.Exchange{
				Response: microtest.Response{
					Status: 200,
					Body:   string(responseGroup),
				},
			},
			uri:   "/group/52f2c725-2cdc-401a-abdd-66db5fd06789",
			group: testGroup,
			e:     nil,
		},
	}

	s := NewService(Config{})
	ms := microtest.MockServer(s)
	defer ms.Server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// add the budget-micro-service exchange
			ms.Append(tc.exchange)

			g, e := s.UpdateGroup(tc.payload)
			// test the error response
			if NotEqualError(tc.e, e) {
				t.Errorf("expected error '%v' got '%v'", tc.e, e)
			}

			// test the group structure returned
			if EqualGroup(g, tc.group) == false {
				t.Errorf("expected group\n'%+v'\ngot\n'%+v'", tc.group, g)
			}

			// test the exchange request URI
			if tc.exchange.Request.RequestURI != tc.uri {
				t.Errorf("expected uri '%v' got '%v'", tc.uri, tc.exchange.Request.RequestURI)
			}
		})
	}
}
