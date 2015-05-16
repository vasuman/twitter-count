package pipeline

import (
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/vasuman/twitter-count/db"
	"github.com/vasuman/twitter-count/expander"
)

func handleTweet(x *anaconda.Tweet) {
	// TODO: time and profile pic
	for _, eURL := range x.Entities.Urls {
		loc := eURL.Expanded_url
		u, err := url.Parse(loc)
		if err != nil {
			logger.Printf("error parsing (%s) - %v\n", loc, err)
			return
		}
		lu, ok := expander.Undirect(u)
		var host string
		if ok {
			host = lu.Host
		} else {
			host = u.Host
		}
		err = db.Associate(host, x.Id, x.User.ScreenName)
		if err != nil {
			logger.Println("error inserting into db, ", err)
		}
	}
}

const sleepTime = 12

func startIn(seconds time.Duration) {
	logger.Printf("starting pipeline in %d seconds\n", seconds)
	time.Sleep(seconds * time.Second)
	stream := api.PublicStreamSample(nil)
	go pipeline(&stream)
}

func pipeline(stream *anaconda.Stream) {
	logger.Println("started tweet pipeline")
	for {
		select {
		case <-stream.Quit:
			logger.Println("got quit. shutting down pipeline")
			startIn(sleepTime)
			return
		case v := <-stream.C:
			switch x := v.(type) {
			case anaconda.Tweet:
				go handleTweet(&x)
			}
		}
	}
}
