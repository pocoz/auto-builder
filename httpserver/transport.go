package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

)

// Service.CreateBuild encoders/decoders.
func encodeCreateBuildRequest(_ context.Context, r *http.Request, request interface{}) error {
	req := request.(createBuildRequest)
	r.URL.Path = "/api/v1/build"

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(req.Payload); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeCreateBuildRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createBuildRequest
	err := json.NewDecoder(r.Body).Decode(&req.Payload)
	return req, err
}

func encodeCreateBuildResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(createBuildResponse)
	if res.Err != nil {
		return encodeError(w, "", res.Err)
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func decodeCreateBuildResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return createBuildResponse{Err: decodeError(r)}, nil
	}
	res := createBuildResponse{}
	return res, nil
}

// errKindToStatus maps service error kinds to the HTTP response codes.
var errKindToStatus = map[ErrorKind]int{
	ErrBadParams:     http.StatusBadRequest,
	ErrNotFound:      http.StatusNotFound,
	ErrConflict:      http.StatusConflict,
	ErrInternal:      http.StatusInternalServerError,
	ErrUnauthorized:  http.StatusUnauthorized,
	ErrForbidden:     http.StatusForbidden,
	ErrNotAcceptable: http.StatusNotAcceptable,
	ErrNotAllowed:    http.StatusMethodNotAllowed,
}

// encodeError writes a service error to the given http.ResponseWriter.
func encodeError(w http.ResponseWriter, description string, err error) error {
	res := &ErrorResponse{}
	status := http.StatusInternalServerError
	if err, ok := err.(*Error); ok {
		if s, ok := errKindToStatus[err.Kind]; ok {
			status = s
		}
		if err.Kind == ErrInternal {
			res.Error = "internal error"
			res.Description = description
		} else {
			res.Error = err.Message
			res.Description = description
		}
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res)
}

// decodeError reads a service error from the given *http.Response.
func decodeError(r *http.Response) error {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, io.LimitReader(r.Body, 1024)); err != nil {
		return fmt.Errorf("%d: %s", r.StatusCode, http.StatusText(r.StatusCode))
	}
	msg := strings.TrimSpace(buf.String())
	if msg == "" {
		msg = http.StatusText(r.StatusCode)
	}
	for kind, status := range errKindToStatus {
		if status == r.StatusCode {
			return &Error{Kind: kind, Message: msg}
		}
	}
	return fmt.Errorf("%d: %s", r.StatusCode, msg)
}
