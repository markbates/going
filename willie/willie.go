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

func (r *Response) URL() string {
	return r.Header().Get("Location")
}

func (w *Willie) JSONPerform(method string, url string, body interface{}) *Response {
	b, _ := json.Marshal(body)
	res, req := w.SetupRequest(method, url, bytes.NewReader(b))
	w.ServeHTTP(res, req)
	w.Cookies = res.Header().Get("Set-Cookie")
	return res
}

func (w *Willie) Perform(method string, u string, body body) *Response {
	if body == nil {
		body = url.Values{}
	}
	res, req := w.SetupRequest(method, u, strings.NewReader(body.Encode()))
	if method == "POST" {
		req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	}
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
	return w.Perform("GET", url, body)
}

func (w *Willie) Post(url string, body body) *Response {
	return w.Perform("POST", url, body)
}

func (w *Willie) Put(url string, body body) *Response {
	return w.Perform("PUT", url, body)
}

func (w *Willie) Delete(url string, body body) *Response {
	return w.Perform("DELETE", url, body)
}

func (w *Willie) JSONGet(url string, body interface{}) *Response {
	return w.JSONPerform("GET", url, body)
}

func (w *Willie) JSONPost(url string, body interface{}) *Response {
	return w.JSONPerform("POST", url, body)
}

func (w *Willie) JSONPut(url string, body interface{}) *Response {
	return w.JSONPerform("PUT", url, body)
}

func (w *Willie) JSONDelete(url string, body interface{}) *Response {
	return w.JSONPerform("DELETE", url, body)
}
