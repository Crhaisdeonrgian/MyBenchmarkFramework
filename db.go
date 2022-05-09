package main

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890йцукенгшщзхъфывапролджэёячсмитьбюЙЦУКЕНГШЩЗХЪФЫВАПРОЛДЖЭЁЯЧСМИТЬБЮ"

func RandStringBytes() string {
	b := make([]byte, 1024)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%194]
	}
	return string(b)
}

func CreateDatabaseTable(db *sql.DB) {
	var err error
	_, err = db.Exec(DropTable)
	_, err = db.Exec(CreateTable)
	if err != nil {
		log.Fatal(err)
	}
}
func FillDataBaseTable(db *sql.DB, do dbOptions) {
	var tx *sql.Tx
	var err error
	for i := 0; i < do.RowCount; i++ {
		currentctx, currentcancel := context.WithTimeout(context.Background(), 12*time.Hour)
		defer currentcancel()
		tx, err = db.BeginTx(currentctx, nil)
		if err != nil {
			log.Fatal(err)
		}
		_, err = tx.ExecContext(currentctx, CustomWrite, RandStringBytes(), RandStringBytes(), RandStringBytes(), RandStringBytes())
		if err != nil {
			log.Fatal(err, tx.Rollback())
		}
		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}
		/*if i%100 == 0 {
			fmt.Println(strconv.Itoa(i))
		}*/
	}
}
func connectToDB(bo benchmarkOptions) *sql.DB {
	var err error
	var dbStd *sql.DB
	testMu.Lock()
	benchTestConfig := sqlConfig
	testMu.Unlock()
	if err != nil {
		log.Fatal(err)
	}
	dbStd, err = sql.Open(bo.driverName, benchTestConfig.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	return dbStd
}

func ShowDatabases() {
	var err error
	var rows *sql.Rows
	rows, err = systemdb.Query("show databases")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	second_params := make([]string, 0)
	for rows.Next() {
		var second string
		if err := rows.Scan(&second); err != nil {
			log.Fatal(err)
		}
		second_params = append(second_params, second)
	}
	log.Println("all the bases")
	log.Println(strings.Join(second_params, " "))
}
func CheckRows(db *sql.DB) int {
	var err error
	var rows *sql.Rows
	queryctx, querycancel := context.WithTimeout(context.Background(), 100000000*time.Millisecond)
	defer querycancel()
	rows, err = db.QueryContext(queryctx, "select count(*) from abobd")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var first int
		if err := rows.Scan(&first); err != nil {
			log.Fatal(err)
		}
		return first
	}
	return 0
}
