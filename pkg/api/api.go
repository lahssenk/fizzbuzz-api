package api

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"

	"github.com/lahssenk/fizzbuzz-api/pkg/service/fizzbuzz"
)

// Exceptions
const (
	// Exception when the request contains invalid data
	ErrInvalidArgument string = "InvalidArgument"
	ErrInternalError   string = "InternalError"
)

// API instance, holding the provided fizzbuzzz service.
type API struct {
	svc *fizzbuzz.Service
}

// NewAPI constructs the API.
func NewAPI() API {
	return API{svc: fizzbuzz.NewService()}
}

// FizzBuzzHandler unmarshals the fizzbuzz input from query parameters and
// returns JSON response with the outcome of fizzbuzz.
// Returns BadRequest if it fails to read from query, or if ComputeFizzBuzzRange fails.
func (a API) FizzBuzzHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error(
			"api: fizzbuzz handler: request.ParseForm failed",
			"err",
			err,
			"path",
			r.URL.Path,
		)
		RespondJSON(
			w,
			ErrorResponse{ErrorCode: ErrInternalError, Message: "something went wrong"},
			http.StatusInternalServerError,
		)
		return
	}

	params, err := parseParamsFromQuery(r.URL.Query())
	if err != nil {
		RespondJSON(
			w,
			ErrorResponse{ErrorCode: ErrInvalidArgument, Message: err.Error()},
			http.StatusBadRequest,
		)
		return
	}

	res, err := a.svc.ComputeFizzBuzzRange(params)
	if err != nil {
		RespondJSON(
			w,
			ErrorResponse{ErrorCode: ErrInvalidArgument, Message: err.Error()},
			http.StatusBadRequest,
		)
		return
	}

	RespondJSON(w, res, http.StatusOK)
}

// ErrorResponse payload
type ErrorResponse struct {
	ErrorCode string `json:"error_code,omitempty"`
	Message   string `json:"message,omitempty"`
	Details   any    `json:"details,omitempty"`
}

func parseParamsFromQuery(values url.Values) (*fizzbuzz.ComputeFizzBuzzRangeParams, error) {
	int1, err := strconv.Atoi(values.Get("int1"))
	if err != nil {
		return nil, errors.Errorf("query parameter int1 must be an integer")
	}

	int2, err := strconv.Atoi(values.Get("int2"))
	if err != nil {
		return nil, errors.Errorf("query parameter int2 must be an integer")
	}

	limit, err := strconv.Atoi(values.Get("limit"))
	if err != nil {
		return nil, errors.Errorf("query parameter limit must be an integer")
	}

	params := fizzbuzz.ComputeFizzBuzzRangeParams{
		String1: values.Get("string1"),
		String2: values.Get("string2"),
		Int1:    int1,
		Int2:    int2,
		Limit:   limit,
	}

	return &params, nil
}

// RespondJSON encodes data as JSON and returns the provided statusCode.
// Returns a default opaque message if fails to encode.
func RespondJSON(w http.ResponseWriter, data any, statusCode int) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)

	if err := enc.Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong"))
		slog.Error("api: respondJSON: failed to encode", "err", err, "data", data)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(buf.Bytes())
}
