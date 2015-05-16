package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"

	"github.com/vasuman/twitter-count/db"
	"github.com/vasuman/twitter-count/expander"
	"github.com/vasuman/twitter-count/pipeline"
	"github.com/vasuman/twitter-count/res"
	"github.com/vasuman/twitter-count/routes"
)

var logger *log.Logger

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type allCreds struct {
	Twitter  pipeline.Creds `yaml:"twitter"`
	Expander expander.Creds `yaml:"expander"`
}

var (
	port    int
	dbPath  string
	logPath string
	creds   string
)

func init() {
	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.StringVar(&dbPath, "db", "", "Data source name of MySQL database")
	flag.StringVar(&logPath, "logfile", "", "Path to log file. Defaults to stdout")
	flag.StringVar(&creds, "creds", "", "Path to file containing credentials")
	flag.Parse()
}

func main() {
	var logFile *os.File
	if logPath != "" {
		var err error
		lfs := os.O_CREATE | os.O_APPEND | os.O_WRONLY
		logFile, err = os.OpenFile(logPath, lfs, 0660)
		panicIf(err)
		defer logFile.Close()
	} else {
		logFile = os.Stdout
	}
	lFlags := log.Lshortfile
	logger = log.New(logFile, "", lFlags)
	dbInst, err := sql.Open("mysql", dbPath)
	panicIf(err)
	defer dbInst.Close()
	err = db.Init(dbInst, logger)
	panicIf(err)
	logger.Println("setup database")
	res.Setup()
	logger.Println("loaded resources")
	c := new(allCreds)
	b, err := ioutil.ReadFile(creds)
	panicIf(err)
	err = yaml.Unmarshal(b, c)
	panicIf(err)
	logger.Println("got credentials")
	expander.Init(c.Expander, logger)
	err = pipeline.Init(c.Twitter, logger)
	panicIf(err)
	logger.Println("setup twitter stream")
	addr := fmt.Sprintf(":%d", port)
	handler := routes.GetRootHandler(logger)
	logger.Println("starting server on, ", addr)
	err = http.ListenAndServe(addr, handler)
	logger.Fatal(err)
}
