package budget

import (
	"bytes"
	"encoding/json"
	"io"
)

// marshalReader is a small wrapper that marshals a struct into json bytes and
// then into an io.Reader and returns an error if there is one.
func marshalReader(v any) (io.Reader, error) {
	xb, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(xb), nil
}
