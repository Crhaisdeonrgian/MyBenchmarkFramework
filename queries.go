package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

var fakeRows *sql.Rows
var fakeResult sql.Result

/*
Fast read is about 10ms
Medium read is about 5s
Slow read is about 3min
*/
const (
	FastRead   string = "SELECT * FROM tbl WHERE col_a"
	MediumRead        = "SELECT * FROM tbl WHERE id<110000 ORDER BY col_b DESC, col_a ASC"
	SlowRead          = "SELECT * FROM tbl first JOIN tbl second ON second.id<5 " +
		"WHERE first.col_a LIKE '%a%' ORDER BY first.col_b DESC, first.col_a ASC"
	CustomRead  = "SELECT * FROM tbl WHERE id<?"
	CustomWrite = "INSERT INTO tbl(col_a, col_b, col_c, col_d) VALUES (?,?,?,?)"
	DropTable   = "DROP TABLE tbl"
	CreateTable = "CREATE TABLE IF NOT EXISTS tbl  (id int AUTO_INCREMENT PRIMARY KEY, a_col nvarchar(1025), b_col nvarchar(1025), c_col nvarchar(1025), d_col nvarchar(1025) )"
)

func executeReadWithTimeout(dbStd *sql.DB, readQuery string, timeout time.Duration) {
	var err error
	queryctx, querycancel := context.WithTimeout(context.Background(), timeout)
	defer querycancel()
	fakeRows, err = dbStd.QueryContext(queryctx, readQuery)
	if err != nil && err != context.DeadlineExceeded {
		log.Fatal("got error in read with timeout: ", err)
	}
	fmt.Println("read with timeout done")
}
func executeRead(dbStd *sql.DB, readQuery string) {
	var err error
	fakeRows, err = dbStd.Query(readQuery)
	if err != nil {
		log.Fatal("got error in read: ", err)
	}
}
func executeWriteWithTimeout(dbStd *sql.DB, writeQuery string, timeout time.Duration) {
	var err error
	queryctx, querycancel := context.WithTimeout(context.Background(), timeout)
	defer querycancel()
	fakeResult, err = dbStd.ExecContext(queryctx, writeQuery)
	if err != nil && err != context.DeadlineExceeded {
		log.Fatal("got error in read with timeout: ", err)
	}
	fmt.Println("read with timeout done")
}
func executeWrite(dbStd *sql.DB, writeQuery string) {
	var err error
	fakeResult, err = dbStd.Exec(writeQuery)
	if err != nil {
		log.Fatal("got error in read: ", err)
	}
}
