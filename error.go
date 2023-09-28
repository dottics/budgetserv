package budget

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Err string

func (err Err) Error() string {
	return string(err)
}

func newError(v interface{}) error {
	return Err(strings.ToLower(fmt.Sprintf("%v", v)))
}

// marshalResponse is a helper function which marshals the response into the
// provided interface. If the status code of the response is not the same as the
// provided status code, then an error is returned.
func marshalResponse(status int, res *http.Response, v interface{}) error {
	xb, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = res.Body.Close()
	if err != nil {
		return err
	}
	//var out bytes.Buffer
	//_ = json.Indent(&out, xb, "", "  ")
	//_, _ = out.WriteTo(os.Stdout)

	if res.StatusCode != status {
		resp := struct {
			Detail interface{} `json:"detail"`
		}{}
		err = json.Unmarshal(xb, &resp)
		if err != nil {
			return err
		}
		return newError(resp.Detail)
	}

	if v != nil {
		err = json.Unmarshal(xb, v)
		if err != nil {
			return err
		}
	}

	return nil
}
