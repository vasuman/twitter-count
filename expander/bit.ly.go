package expander

import (
	"encoding/json"
	"errors"
	"net/url"
)

func bitly(token string) *jsonAPI {
	e := new(jsonAPI)
	e.name = "bit.ly"
	e.path = "https://api-ssl.bitly.com/v3/expand"
	e.query = func(u string) url.Values {
		v := make(url.Values, 2)
		v.Add("access_token", token)
		v.Add("shortUrl", u)
		return v
	}
	e.decode = func(dec *json.Decoder) (string, error) {
		type response struct {
			Data struct {
				Expand []struct {
					Long string `json:"long_url"`
				} `json:"expand"`
			} `json:"data"`
		}
		resp := new(response)
		err := dec.Decode(resp)
		if err != nil {
			return "", err
		}
		long := resp.Data.Expand[0].Long
		if long == "" {
			return "", errors.New("got empty response")
		}
		return long, nil
	}
	return e
}
