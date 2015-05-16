package expander

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type Creds struct {
	Bitly string `yaml:"bitly"`
	Googl string `yaml:"googl"`
}

var (
	logger *log.Logger
)

type iface interface {
	Expand(string) (*url.URL, bool)
}

var hostMap map[string]iface

func Undirect(u *url.URL) (*url.URL, bool) {
	exp, ok := hostMap[u.Host]
	if !ok {
		return nil, false
	}
	return exp.Expand(u.String())
}

func Init(c Creds, logInst *log.Logger) {
	logger = logInst
	hostMap = map[string]iface{
		"bit.ly": bitly(c.Bitly),
		"goo.gl": googl(c.Googl),
	}
}

type jsonAPI struct {
	name   string
	path   string
	query  func(string) url.Values
	decode func(*json.Decoder) (string, error)
}

func (e *jsonAPI) Expand(u string) (*url.URL, bool) {
	v := e.query(u)
	r, err := http.Get(e.path + "?" + v.Encode())
	if err != nil {
		logger.Printf("error with %s link (%s) - %v\n", e.name, u, err)
		return nil, false
	}
	if r.StatusCode != http.StatusOK {
		logger.Printf("recieved not ok from %s\n", e.name)
		return nil, false
	}
	dec := json.NewDecoder(r.Body)
	l, err := e.decode(dec)
	if err != nil {
		logger.Printf("error decoding %s response - %v\n", e.name, err)
		return nil, false
	}
	long, err := url.Parse(l)
	if err != nil {
		logger.Println("error parsing URL from %s %v\n", e.name, err)
		return nil, false
	}
	return long, true
}
