package routes

import (
	"encoding/json"
	"net/http"

	"github.com/vasuman/twitter-count/db"
)

func listDomains(w http.ResponseWriter, r *http.Request) {
	req := new(db.DomainParams)
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(req)
	if err != nil {
		internalError(w, err)
		return
	}
	ds, err := db.GetTopDomains(req)
	if err != nil {
		internalError(w, err)
		return
	}
	enc := json.NewEncoder(w)
	enc.Encode(ds)
}

func getTweets(w http.ResponseWriter, r *http.Request) {
	req := new(db.TweetParams)
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(req)
	if err != nil {
		internalError(w, err)
		return
	}
	ts, err := db.GetDomainTweets(req)
	if err != nil {
		internalError(w, err)
		return
	}
	enc := json.NewEncoder(w)
	enc.Encode(ts)
}
