package api_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lahssenk/fizzbuzz-api/pkg/api"
)

func TestAPI_FizzBuzzHandler(t *testing.T) {
	tests := []struct {
		name       string
		req        *http.Request
		statusCode int
		respBody   []byte
	}{
		{
			name: "OK",
			req: httptest.NewRequest(
				"GET",
				"http://localhost:8080/fizzbuzz?string1=f&string2=b&int1=3&int2=5&limit=3",
				nil,
			),
			statusCode: 200,
			respBody:   []byte("{\"data\":[\"1\",\"2\",\"f\"]}\n"),
		},
		{
			name: "invalid argument",
			req: httptest.NewRequest(
				"GET",
				"http://localhost:8080/fizzbuzz?string2=b&int1=3&int2=5&limit=3",
				nil,
			),
			statusCode: 400,
			respBody: []byte(
				"{\"error_code\":\"InvalidArgument\",\"message\":\"string1 required\"}\n",
			),
		},
	}

	apiHandlers := api.NewAPI()

	for ind := range tests {
		test := tests[ind]
		t.Run(test.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			apiHandlers.FizzBuzzHandler(rec, test.req)

			statusCode := rec.Result().StatusCode
			respBody := rec.Result().Body
			if respBody != nil {
				defer respBody.Close()
			}

			data, err := io.ReadAll(respBody)
			if err != nil {
				t.Fatal(err)
			}

			require.Equal(t, test.respBody, data)
			require.Equal(t, test.statusCode, statusCode)
		})
	}
}
