package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	writer http.ResponseWriter
	status int
	body   interface{}
}

func New(w http.ResponseWriter) *Response {
	return &Response{
		writer: w,
		status: http.StatusOK,
	}
}

func (r *Response) Status(status int) *Response {
	r.status = status
	return r
}
func Ok(r *Response) {
	r.writer.WriteHeader(r.status) // need to validate that a correct status is set
	json.NewEncoder(r.writer).Encode(r.body)
}

func (r *Response) Json(body interface{}) *Response {
	r.body = body
	r.writer.Header().Set("Content-Type", "application/json")
	Ok(r)
	return r
}
