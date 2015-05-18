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
	createTweet = `
CREATE TABLE IF NOT EXISTS Tweet (
	Domain VARCHAR(25),
	Id BIGINT,
	User VARCHAR(15),
	INDEX(Domain)
);
`
	createCnt = `
CREATE TABLE IF NOT EXISTS Cnt (
	Domain VARCHAR(25) PRIMARY KEY,
	Count BIGINT,
	INDEX(Count)
);

`
	topDomains = `
SELECT Domain, Count FROM Cnt
ORDER BY Count DESC
LIMIT ? OFFSET ?
`
	insertTweet = `
INSERT INTO Tweet (Domain, Id, User)
VALUES (?, ?, ?);
`
	updateCount = `
INSERT INTO Cnt (Domain, Count) 
VALUES (?, 1)
ON DUPLICATE KEY 
UPDATE Count = Count + 1;
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
	stmtUpdateCount *sql.Stmt
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
	_, err = db.Exec(createTweet)
	_, err = db.Exec(createCnt)
	stmtTopDomains = prepare(topDomains)
	stmtInsertTweet = prepare(insertTweet)
	stmtGetTweets = prepare(getTweets)
	stmtUpdateCount = prepare(updateCount)
	logger = logInst
	return err
}
