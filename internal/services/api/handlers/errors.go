package handlers

import (
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape/problems"
)

// NotFound - return 404 error
func NotFound(msg, field string) *jsonapi.ErrorObject {
	result := problems.NotFound()
	result.Detail = msg
	if field != "" {
		result.Meta = &map[string]interface{}{
			"what": field,
		}
	}
	return result
}
