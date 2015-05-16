package pipeline

import (
	"log"

	"github.com/ChimeraCoder/anaconda"
)

type Creds struct {
	Client struct {
		Key    string `yaml:"key"`
		Secret string `yaml:"secret"`
	} `yaml:"client"`
	Access struct {
		Token  string `yaml:"token"`
		Secret string `yaml:"secret"`
	} `yaml:"access"`
}

var (
	logger *log.Logger
	api    *anaconda.TwitterApi
)

func Init(c Creds, logInst *log.Logger) error {
	logger = logInst
	anaconda.SetConsumerKey(c.Client.Key)
	anaconda.SetConsumerSecret(c.Client.Secret)
	api = anaconda.NewTwitterApi(c.Access.Token, c.Access.Secret)
	go startIn(0) // start now
	return nil
}
