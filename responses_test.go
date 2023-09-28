package budget

import (
	"encoding/json"
	"github.com/google/uuid"
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
