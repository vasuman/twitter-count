package db

import (
	"database/sql"
	"log"
)

var (
	logger *log.Logger
	db     *sql.DB
)

const (
	dbSchema = `
CREATE TABLE IF NOT EXISTS Tweet (
	Domain VARCHAR(25),
	Id BIGINT,
	User VARCHAR(15),
	INDEX(Domain)
);
`
	topDomains = `
SELECT Domain, COUNT(Id) AS Count FROM Tweet 
GROUP BY Domain 
ORDER BY Count DESC
LIMIT ? OFFSET ?
`
	insertTweet = `
INSERT INTO Tweet (Domain, Id, User)
VALUES (?, ?, ?)
`
	getTweets = `
SELECT Id, User FROM Tweet WHERE Domain = ?
LIMIT ? OFFSET ?
`
)

var (
	stmtTopDomains  *sql.Stmt
	stmtInsertTweet *sql.Stmt
	stmtGetTweets   *sql.Stmt
)

func Init(dbInst *sql.DB, logInst *log.Logger) error {
	var err error
	prepare := func(s string) (ret *sql.Stmt) {
		if err == nil {
			ret, err = db.Prepare(s)
		}
		return
	}
	db = dbInst
	_, err = db.Exec(dbSchema)
	stmtTopDomains = prepare(topDomains)
	stmtInsertTweet = prepare(insertTweet)
	stmtGetTweets = prepare(getTweets)
	logger = logInst
	return err
}
