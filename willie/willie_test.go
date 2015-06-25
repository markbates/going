package willie_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/gorilla/context"
	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/markbates/going/willie"
	"github.com/stretchr/testify/require"
)

var Store sessions.Store = sessions.NewCookieStore([]byte("something-very-secret"))

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

func Test_WillieSessions(t *testing.T) {
	r := require.New(t)
	w := willie.New(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		defer context.Clear(req)
		sess, _ := Store.Get(req, "my-session")
		t := sess.Values["foo"]
		fmt.Printf("t: %s\n", t)
		if t != nil {
			res.WriteHeader(200)
			fmt.Fprint(res, t)
		} else {
			sess.Values["foo"] = "bar"
			sess.Save(req, res)
			res.WriteHeader(201)
			fmt.Fprint(res, "setting session")
		}
	}))

	res := w.Get("/", nil)
	r.Equal(201, res.Code)
	r.Equal("setting session", res.Body.String())

	res = w.Get("/", nil)
	r.Equal(200, res.Code)
	r.Equal("bar", res.Body.String())
}
