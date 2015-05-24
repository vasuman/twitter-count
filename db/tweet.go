package db

import (
	"database/sql"
	"fmt"
)

func Associate(dom string, id int64, handle string) error {
	tx, err := db.Begin()
	execStmt := func(stmt *sql.Stmt, args ...interface{}) {
		if err == nil {
			_, err = tx.Stmt(stmt).Exec(args...)
		}
	}
	execStmt(stmtInsertTweet, dom, id, handle)
	execStmt(stmtUpdateCount, dom)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

type DomainParams struct {
	Start int `json:"start"`
	Num   int `json:"num"`
}

type DomainItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func GetTopDomains(q *DomainParams) ([]DomainItem, error) {
	ret := make([]DomainItem, 0, q.Num)
	rows, err := stmtTopDomains.Query(q.Num, q.Start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			d string
			c int
		)
		err = rows.Scan(&d, &c)
		if err != nil {
			return nil, err
		}
		ret = append(ret, DomainItem{d, c})
	}
	return ret, nil
}

type TweetParams struct {
	Domain string `json:"domain"`
	Idx    int    `json:"idx"`
	Num    int    `json:"num"`
}

type TweetItem struct {
	Id         string `json:"id"`
	UserHandle string `json:"userHandle"`
}

func GetDomainTweets(q *TweetParams) ([]TweetItem, error) {
	ret := make([]TweetItem, 0, q.Num)
	rows, err := stmtGetTweets.Query(q.Domain, q.Num, q.Idx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id     int64
			handle string
		)
		err = rows.Scan(&id, &handle)
		if err != nil {
			return nil, err
		}
		ret = append(ret, TweetItem{fmt.Sprintf("%d", id), handle})
	}
	return ret, nil
}
