package willie

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

type body interface {
	Encode() string
}

type Willie struct {
	http.Handler
	Cookies string
}

func New(a http.Handler) *Willie {
	return &Willie{Handler: a, Cookies: ""}
}

type Response struct {
	*httptest.ResponseRecorder
}

func (r *Response) Bind(x interface{}) {
	json.NewDecoder(r.Body).Decode(&x)
}

func (w *Willie) jperform(method string, url string, body interface{}) *Response {
	b, _ := json.Marshal(body)
	return w.perform(method, url, bytes.NewReader(b))
}

func (w *Willie) xperform(method string, u string, body body) *Response {
	if body == nil {
		body = url.Values{}
	}
	return w.perform(method, u, strings.NewReader(body.Encode()))
}

func (w *Willie) perform(method string, url string, body io.Reader) *Response {
	res, req := w.SetupRequest(method, url, body)
	w.ServeHTTP(res, req)
	w.Cookies = res.Header().Get("Set-Cookie")
	return res
}

func (w *Willie) SetupRequest(method string, url string, body io.Reader) (*Response, *http.Request) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Cookie", w.Cookies)
	return &Response{res}, req
}

func (w *Willie) Get(url string, body body) *Response {
	return w.xperform("GET", url, body)
}

func (w *Willie) Post(url string, body body) *Response {
	return w.xperform("POST", url, body)
}

func (w *Willie) Put(url string, body body) *Response {
	return w.xperform("PUT", url, body)
}

func (w *Willie) Delete(url string, body body) *Response {
	return w.xperform("DELETE", url, body)
}

func (w *Willie) JSONGet(url string, body interface{}) *Response {
	return w.jperform("GET", url, body)
}

func (w *Willie) JSONPost(url string, body interface{}) *Response {
	return w.jperform("POST", url, body)
}

func (w *Willie) JSONPut(url string, body interface{}) *Response {
	return w.jperform("PUT", url, body)
}

func (w *Willie) JSONDelete(url string, body interface{}) *Response {
	return w.jperform("DELETE", url, body)
}
