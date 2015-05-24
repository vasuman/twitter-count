package expander

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func googl(key string) *jsonAPI {
	e := new(jsonAPI)
	e.name = "goo.gl"
	e.path = "https://www.googleapis.com/urlshortener/v1/url"
	e.query = func(u string) url.Values {
		v := make(url.Values, 2)
		v.Add("key", key)
		v.Add("shortUrl", u)
		return v
	}
	e.decode = func(dec *json.Decoder) (string, error) {
		var r struct {
			Long   string `json:"longUrl"`
			Status string `json:"status"`
		}
		err := dec.Decode(&r)
		if err != nil {
			return "", err
		}
		if r.Status != "OK" {
			return "", fmt.Errorf("status %s", r.Status)
		}
		return r.Long, nil
	}
	return e
}
