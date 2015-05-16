package routes

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/vasuman/twitter-count/res"
)

var logger *log.Logger

func wrongMethod(w http.ResponseWriter, meth string) {
	code := http.StatusMethodNotAllowed
	txt := fmt.Sprintf("%s not allowed", meth)
	http.Error(w, txt, code)
}

func badRequest(w http.ResponseWriter, txt string) {
	code := http.StatusBadRequest
	http.Error(w, "bad request: "+txt, code)
}

func internalError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	txt := fmt.Sprintf("internal error - %v", err)
	http.Error(w, txt, code)
}

func serveHTML(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(res.HTML[path])
	}
}

func serveStatic(prefix, cType string, m map[string][]byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			wrongMethod(w, r.Method)
			return
		}
		path := strings.TrimPrefix(r.URL.Path, prefix)
		b, ok := m[path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Add("Content-Type", cType)
		w.Write(b)
	}
}

func execTemplate(w http.ResponseWriter, name string, ctx interface{}) {
	var b bytes.Buffer
	err := res.Template.ExecuteTemplate(&b, name, ctx)
	if err != nil {
		logger.Printf("error in template %s - %v\n", name, err)
		internalError(w, err)
		return
	}
	io.Copy(w, &b)
}

func exactWrap(path, method string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			http.NotFound(w, r)
			return
		}
		if r.Method != method {
			wrongMethod(w, r.Method)
			return
		}
		fn(w, r)
	}
}

func GetRootHandler(logInst *log.Logger) http.Handler {
	logger = logInst
	r := http.NewServeMux()
	register := func(path, meth string, fn http.HandlerFunc) {
		r.HandleFunc(path, exactWrap(path, meth, fn))
	}
	static := func(prefix, contentType string, m map[string][]byte) {
		r.HandleFunc(prefix, serveStatic(prefix, contentType, m))
	}
	register("/", "GET", serveHTML("dashboard.html"))
	register("/api/listDomains", "POST", listDomains)
	register("/api/getTweets", "POST", getTweets)
	register("/domain", "GET", showDomain)
	static("/scripts/", "text/javascript", res.Scripts)
	static("/styles/", "text/css", res.Styles)
	return r
}

func showDomain(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	name := v.Get("name")
	if name == "" {
		badRequest(w, "need a name")
		return
	}
	execTemplate(w, "domain", name)
}
