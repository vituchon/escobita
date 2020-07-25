// Routines for dealing with any postgres dbms bureaucracy

package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Connection = *sql.DB

type Tx = *sql.Tx

// Common interface for sql.Tx and sql.Db regarding Exec, // see : https://github.com/golang/go/issues/14468
type DbTx interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Server struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	DatabaseName string `json:"databaseName"`
}

type Config struct {
	Credentials Credentials `json:"credentials"`
	Server      Server      `json:"server"`
}

func OpenConnection(config Config) (conn Connection, err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Server.Host, config.Server.Port, config.Credentials.Username, config.Credentials.Password, config.Server.DatabaseName)
	conn, err = sql.Open("postgres", psqlInfo)
	if err == nil {
		_, err = conn.Exec("SELECT 1")
	}
	return
}

// Represents the block of code intended to be executed within a given database's transacton
type TransactionFunc = func(*sql.Tx) (interface{}, error)

// Takes care of executing a transaction's block handling (and hiding) all the transaction specific commands.
func ExecuteTransactionFunc(conn Connection, txFunc TransactionFunc) (result interface{}, err error) {
	tx, err := conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			log.Println("Doing Rollback")
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Printf("Error while doing rollback, error was: '%v'\n", rollbackErr)
			}
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()
	result, err = txFunc(tx)
	return
}
