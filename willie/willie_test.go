package willie_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/gorilla/pat"
	"github.com/markbates/going/willie"
	"github.com/stretchr/testify/require"
)

type record struct {
	Method string
	Body   []byte
}

func app() http.Handler {
	p := pat.New()

	f := func(res http.ResponseWriter, req *http.Request) {
		m := record{Method: req.Method}
		m.Body, _ = ioutil.ReadAll(req.Body)
		json.NewEncoder(res).Encode(&m)
	}
	p.Get("/get", f)
	p.Post("/post", f)
	p.Put("/put", f)
	p.Delete("/delete", f)
	return p
}

type methodHandler func(string, interface{}) *willie.Response

func Test_Willie(t *testing.T) {
	a := require.New(t)
	w := willie.New(app())

	m := map[string]methodHandler{
		"get":    w.Get,
		"post":   w.Post,
		"put":    w.Put,
		"delete": w.Delete,
	}
	for k, h := range m {
		res := h("/"+k, []string{"a", "b"})
		r := &record{}
		res.Bind(r)
		a.Equal(strings.ToUpper(k), r.Method)
		a.Equal(`["a","b"]`, string(r.Body))
	}
}
