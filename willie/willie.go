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

func (w *Willie) setupRequest(method string, url string, body interface{}) (*Response, *http.Request) {
	b, _ := json.Marshal(body)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, bytes.NewReader(b))
	return &Response{res}, req
}

func (w *Willie) Get(url string, body interface{}) *Response {
	res, req := w.setupRequest("GET", url, body)
	w.ServeHTTP(res, req)
	return res
}

func (w *Willie) Post(url string, body interface{}) *Response {
	res, req := w.setupRequest("POST", url, body)
	w.ServeHTTP(res, req)
	return res
}

func (w *Willie) Put(url string, body interface{}) *Response {
	res, req := w.setupRequest("PUT", url, body)
	w.ServeHTTP(res, req)
	return res
}

func (w *Willie) Delete(url string, body interface{}) *Response {
	res, req := w.setupRequest("DELETE", url, body)
	w.ServeHTTP(res, req)
	return res
}
