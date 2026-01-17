package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/crazydw4rf/oil-bank-backend/internal/delivery/http/response"
)

func TestHTTPResponse(t *testing.T) {
	type Bio struct {
		Name string `json:"name"`
	}

	names := make([]*Bio, 0)
	names = append(names, &Bio{Name: "Dika"})
	names = append(names, &Bio{Name: "Ucup"})
	names = append(names, &Bio{Name: "Budi"})

	var res response.HTTPResponse[[]*Bio]
	res.Data = names
	res.Code = http.StatusTeapot
	res.Error = response.HTTPError{Message: "Sybau", KnownError: true}

	b, _ := json.Marshal(res)
	fmt.Printf("%#v\n", string(b))
}
