package willie

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

type Willie struct {
	http.Handler
}

func New(a http.Handler) *Willie {
	return &Willie{a}
}

type Response struct {
	*httptest.ResponseRecorder
}

func (r *Response) Bind(x interface{}) {
	json.NewDecoder(r.Body).Decode(&x)
}

func (w *Willie) perform(method string, url string, body interface{}) *Response {
	res, req := w.setupRequest(method, url, body)
	w.ServeHTTP(res, req)
	return res
}

func (w *Willie) setupRequest(method string, url string, body interface{}) (*Response, *http.Request) {
	b, _ := json.Marshal(body)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, bytes.NewReader(b))
	return &Response{res}, req
}

func (w *Willie) Get(url string, body interface{}) *Response {
	return w.perform("GET", url, body)
}

func (w *Willie) Post(url string, body interface{}) *Response {
	return w.perform("POST", url, body)
}

func (w *Willie) Put(url string, body interface{}) *Response {
	return w.perform("PUT", url, body)
}

func (w *Willie) Delete(url string, body interface{}) *Response {
	return w.perform("DELETE", url, body)
}
