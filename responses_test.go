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

var noAuth = `{
	"detail": "No token"
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
