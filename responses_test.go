package budget

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// minify removes all whitespace from a string.
func minify(s string) string {
	var r string
	for _, c := range s {
		if c != '\n' && c != '\t' && c != ' ' {
			r += string(c)
		}
	}
	return r
}

var noPermission = `{
	"detail": "No permission"
}`

var notFound = `{
	"detail": "Not found"
}`

func errorResponseDetail(s string) string {
	return `{
	"detail": ` + s + `
}`
}

func response(message, data string) string {
	return `{
	"message": "` + message + `",
	"data": ` + data + `
}`
}

var testBudget = Budget{
	UUID:        uuid.MustParse("f5fca9d0-e308-4ff2-be4e-aff22a4c2a78"),
	EntityUUID:  uuid.MustParse("67b14c0f-b8ea-4f0f-bf07-cadc73cd74d9"),
	Name:        "test budget",
	Description: "test budget description",
}

var responseBudget = json.RawMessage(`{
	"message": "budget retrieved",
	"data": {
		"budget": {
			"uuid": "f5fca9d0-e308-4ff2-be4e-aff22a4c2a78",
			"entity_uuid": "67b14c0f-b8ea-4f0f-bf07-cadc73cd74d9",
			"name": "test budget",
			"description": "test budget description"
		}
	}
}`)

var testBudgetWithGroups = Budget{
	UUID:        uuid.MustParse("f5fca9d0-e308-4ff2-be4e-aff22a4c2a78"),
	EntityUUID:  uuid.MustParse("67b14c0f-b8ea-4f0f-bf07-cadc73cd74d9"),
	Name:        "test budget",
	Description: "test budget description",
	Groups: Groups{
		{
			UUID: uuid.MustParse("52f2c725-2cdc-401a-abdd-66db5fd06789"),
			Name: "income",
			SubGroups: Groups{
				{
					UUID: uuid.MustParse("b8448a78-6417-4fe2-849c-024622bc6106"),
					Name: "base salary",
					Items: Items{
						{
							UUID: uuid.MustParse("2bd06946-c355-4198-8766-949149331e04"),
							Name: "sold old items",
							Events: Events{
								{
									UUID:      uuid.MustParse("c130f4d9-0124-4f0c-8129-298ae60cd9f1"),
									Name:      "go pro",
									Amount:    8000,
									Debit:     true,
									Credit:    false,
									StartDate: time.Date(2023, 3, 12, 0, 0, 0, 0, time.UTC),
									EndDate:   time.Date(2023, 3, 12, 0, 0, 0, 0, time.UTC),
								},
							},
						},
					},
				},
			},
		},
		{
			UUID: uuid.MustParse("eea51d45-c9bd-45e2-bc80-010ecbb7a0d3"),
			Name: "investments",
		},
		{
			UUID: uuid.MustParse("6be3df72-da3d-4a8c-bef6-d0b57120b80a"),
			Name: "expenses",
		},
	},
}

var responseBudgetWithGroups = json.RawMessage(`{
	"message": "budget retrieved",
	"data": {
		"budget": {
			"uuid": "f5fca9d0-e308-4ff2-be4e-aff22a4c2a78",
			"entity_uuid": "67b14c0f-b8ea-4f0f-bf07-cadc73cd74d9",
			"name": "test budget",
			"description": "test budget description",
			"groups":[
				{
					"uuid":"52f2c725-2cdc-401a-abdd-66db5fd06789",
					"name":"income",
					"sub_groups":[
						{
							"uuid":"b8448a78-6417-4fe2-849c-024622bc6106",
							"name":"base salary",
							"sub_groups":[]
						 }
					],
					"items": [
						{
							"uuid": "2bd06946-c355-4198-8766-949149331e04",
							"name": "sold old items",
							"events": [
								{
									"uuid": "c130f4d9-0124-4f0c-8129-298ae60cd9f1",
									"name": "go pro",
									"amount": 8000,
									"debit": true,
									"credit": false,
									"start_date": "2023-03-12T00:00:00Z",
									"end_date": "2023-03-12T00:00:00Z"
								}
							]
						}
					]
				},
				{
					"uuid":"eea51d45-c9bd-45e2-bc80-010ecbb7a0d3",
					"name":"investments",
					"sub_groups":[]
				},
				{
					"uuid":"6be3df72-da3d-4a8c-bef6-d0b57120b80a",
					"name":"expenses",
					"sub_groups":[]
				}
			]
		}
	}
}`)

var responseBudgets = json.RawMessage(`{
	"message": "budgets retrieved",
	"data": {
		"budgets": [
			{
				"uuid":"4ebae4ad-803c-4487-98ea-3f1f926e59e6",
				"entity_uuid":"ecdccfe9-95fe-4c9f-bd86-169ad67c445a",
				"name":"test budget uno",
				"active":true
		    },
		    {
				"uuid":"0d79e5cb-5b26-49bc-a5fa-3b39e2710675",
				"entity_uuid":"ecdccfe9-95fe-4c9f-bd86-169ad67c445a",
				"name":"test budget dos",
				"active":true
		    }
		]
	}	
}`)

var testGroupBasic = Group{
	UUID:        uuid.MustParse("52f2c725-2cdc-401a-abdd-66db5fd06789"),
	Name:        "income",
	Description: "income description",
}

var responseGroupBasic = json.RawMessage(`{
	"message": "group retrieved",
	"data": {
		"group": {
			"uuid": "52f2c725-2cdc-401a-abdd-66db5fd06789",
			"name": "income",
			"description": "income description"
		}
	}
}`)

var testGroup = Group{
	UUID:      uuid.MustParse("52f2c725-2cdc-401a-abdd-66db5fd06789"),
	Name:      "income",
	SubGroups: Groups{},
	Items: Items{
		{
			UUID: uuid.MustParse("2bd06946-c355-4198-8766-949149331e04"),
			Name: "salary",
			Events: Events{
				{
					UUID:      uuid.MustParse("c130f4d9-0124-4f0c-8129-298ae60cd9f1"),
					Name:      "work salary",
					Amount:    30000,
					Debit:     true,
					Credit:    false,
					StartDate: time.Date(2023, 1, 25, 0, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2023, 12, 20, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	},
}

var responseGroup = json.RawMessage(`{
	"message": "group retrieved",
	"data": {
		"group": {
			"uuid": "52f2c725-2cdc-401a-abdd-66db5fd06789",
			"name": "income",
			"sub_groups": [],
			"items": [
				{
					"uuid": "2bd06946-c355-4198-8766-949149331e04",
					"name": "salary",
					"events": [
						{
							"uuid": "c130f4d9-0124-4f0c-8129-298ae60cd9f1",
							"name": "work salary",
							"amount": 30000,
							"debit": true,
							"credit": false,
							"start_date": "2023-01-25T00:00:00Z",
							"end_date": "2023-12-20T00:00:00Z"
						}
					]
				}
			]
		}
	}
}`)

var testGroupItems = Items{
	{
		UUID: uuid.MustParse("2bd06946-c355-4198-8766-949149331e04"),
		Name: "salary",
		Events: Events{
			{
				UUID:      uuid.MustParse("c130f4d9-0124-4f0c-8129-298ae60cd9f1"),
				Name:      "work salary",
				Amount:    30000,
				Debit:     true,
				Credit:    false,
				StartDate: time.Date(2023, 1, 25, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2023, 12, 20, 0, 0, 0, 0, time.UTC),
			},
			{
				UUID:      uuid.MustParse("8feec066-2dfb-44f5-b353-1cb6e75c3084"),
				Name:      "bonus",
				Amount:    5000,
				Debit:     true,
				Credit:    false,
				StartDate: time.Date(2023, 12, 20, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2023, 12, 20, 0, 0, 0, 0, time.UTC),
			},
		},
	},
	{
		UUID: uuid.MustParse("0aab721b-4224-474c-9df8-e77117a31a02"),
		Name: "sold old items",
		Events: Events{
			{
				UUID:      uuid.MustParse("f71ce1d0-0ddd-4a39-9abd-baea8a6d8bbe"),
				Name:      "go pro",
				Amount:    8000,
				Debit:     true,
				Credit:    false,
				StartDate: time.Date(2023, 3, 12, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2023, 3, 12, 0, 0, 0, 0, time.UTC),
			},
		},
	},
}

var responseGroupItems = json.RawMessage(`{
	"message": "items retrieved",
	"data": {
		"items": [
			{
				"uuid": "2bd06946-c355-4198-8766-949149331e04",
				"name": "salary",
				"events": [
					{
						"uuid": "c130f4d9-0124-4f0c-8129-298ae60cd9f1",
						"name": "work salary",
						"amount": 30000,
						"debit": true,
						"credit": false,
						"start_date": "2023-01-25T00:00:00Z",
						"end_date": "2023-12-20T00:00:00Z"
					},
					{
						"uuid": "8feec066-2dfb-44f5-b353-1cb6e75c3084",
						"name": "bonus",
						"amount": 5000,
						"debit": true,
						"credit": false,
						"start_date": "2023-12-20T00:00:00Z",
						"end_date": "2023-12-20T00:00:00Z"
					}
				]
			},
			{
				"uuid": "0aab721b-4224-474c-9df8-e77117a31a02",
				"name": "sold old items",
				"events": [
					{
						"uuid": "f71ce1d0-0ddd-4a39-9abd-baea8a6d8bbe",
						"name": "go pro",
						"amount": 8000,
						"debit": true,
						"credit": false,
						"start_date": "2023-03-12T00:00:00Z",
						"end_date": "2023-03-12T00:00:00Z"
					}
				]
			}
		]
	}
}`)

var testItemNew = Item{
	UUID:        uuid.MustParse("2bd06946-c355-4198-8766-949149331e04"),
	Name:        "sold old items",
	Description: "sold old items description",
	Events:      Events{},
}

var responseItemNew = json.RawMessage(`{
	"message": "item created",
	"data": {
		"item": {
			"uuid": "2bd06946-c355-4198-8766-949149331e04",
			"name": "sold old items",
			"description": "sold old items description",
			"events": []
		}
	}
}`)

var testItem = Item{
	UUID:        uuid.MustParse("2bd06946-c355-4198-8766-949149331e04"),
	Name:        "sold old items",
	Description: "sold old items description",
	Events: Events{
		{
			UUID:      uuid.MustParse("f71ce1d0-0ddd-4a39-9abd-baea8a6d8bbe"),
			Name:      "go pro",
			Amount:    8000,
			Debit:     true,
			Credit:    false,
			StartDate: time.Date(2023, 3, 12, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2023, 3, 12, 0, 0, 0, 0, time.UTC),
		},
		{
			UUID:      uuid.MustParse("8feec066-2dfb-44f5-b353-1cb6e75c3084"),
			Name:      "computer screen",
			Amount:    5000,
			Debit:     true,
			Credit:    false,
			StartDate: time.Date(2023, 5, 6, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2023, 5, 6, 0, 0, 0, 0, time.UTC),
		},
	},
}

var responseItem = json.RawMessage(`{
	"message": "item retrieved",
	"data": {
		"item": {
			"uuid": "2bd06946-c355-4198-8766-949149331e04",
			"name": "sold old items",
			"description": "sold old items description",
			"events": [
				{
					"uuid": "f71ce1d0-0ddd-4a39-9abd-baea8a6d8bbe",
					"name": "go pro",
					"amount": 8000,
					"debit": true,
					"credit": false,
					"start_date": "2023-03-12T00:00:00Z",
					"end_date": "2023-03-12T00:00:00Z"
				},
				{
					"uuid": "8feec066-2dfb-44f5-b353-1cb6e75c3084",
					"name": "computer screen",
					"amount": 5000,
					"debit": true,
					"credit": false,
					"start_date": "2023-05-06T00:00:00Z",
					"end_date": "2023-05-06T00:00:00Z"
				}
			]
		}
	}
}`)

var testCategories = []Category{
	{
		UUID:        uuid.MustParse("f5fca9d0-e308-4ff2-be4e-aff22a4c2a78"),
		Name:        "test category",
		Description: "test category description",
		Norm:        false,
	},
	{
		UUID:        uuid.MustParse("67b14c0f-b8ea-4f0f-bf07-cadc73cd74d9"),
		Name:        "test category 2",
		Description: "test category 2 description",
		Norm:        true,
	},
}

var responseCategories = json.RawMessage(`{
	"message": "categories retrieved",
	"data": {
		"categories": [
			{
				"uuid": "f5fca9d0-e308-4ff2-be4e-aff22a4c2a78",
				"name": "test category",
				"description": "test category description",
				"norm": false
			},
			{
				"uuid": "67b14c0f-b8ea-4f0f-bf07-cadc73cd74d9",
				"name": "test category 2",
				"description": "test category 2 description",
				"norm": true
			}
		]
	}
}`)

var testCategory = Category{
	UUID:        uuid.MustParse("f5fca9d0-e308-4ff2-be4e-aff22a4c2a78"),
	Name:        "test category",
	Description: "test category description",
	Norm:        false,
}

var responseCategory = json.RawMessage(`{
	"message": "category retrieved",
	"data": {
		"category": {
			"uuid": "f5fca9d0-e308-4ff2-be4e-aff22a4c2a78",
			"name": "test category",
			"description": "test category description",
			"norm": false
		}
	}
}`)
