package todolist

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHTTPHandler(svc Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(svc)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}
	r.Methods("POST").Path("/tasks").Handler(httptransport.NewServer(
		e.CreateTaskEndpoint,
		DecodeCreateTaskRequest,
		EncodeResponse,
		options...,
	))
	r.Methods("GET").Path("/tasks/{id}").Handler(httptransport.NewServer(
		e.GetTaskEndpoint,
		DecodeGetTaskRequest,
		EncodeResponse,
		options...,
	))
	return r
}

func DecodeCreateTaskRequest(_ context.Context, r *http.Request) (interface{}, error) {
	// decode json data from request body
	var request CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&request.payload); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeGetTaskRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	// get raw id from json
	stringID, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	// convert to unsigned 64 bit
	u64ID, err := strconv.ParseUint(stringID, 10, 32)
	if err != nil {
		msg := fmt.Sprintf("value '%s' is not a valid id", stringID)
		return nil, errors.New(msg)
	}
	// convert to unsigned 32 bit
	u32ID := uint(u64ID)
	return GetTaskRequest{ID: u32ID}, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
