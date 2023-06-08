// Helper rutines
package postgres

import (
	"database/sql"
	"log"
)

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
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Printf("Error while doing rollback, error was: '%v'\n", rollbackErr)
			}
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

// Interface for things that performs an scan over a given row.
// Actually it is a common interface for https://pkg.go.dev/database/sql#Rows.Scan and https://pkg.go.dev/database/sql#Row.Scan
type RowScanner interface {
	Scan(dest ...interface{}) error
}

// Scans a single row from a given query
type RowScanFunc func(rows RowScanner) (interface{}, error)

// Scans multiples rows using a scanner function in order to build a new "scanable" struct
func ScanMultiples(rows *sql.Rows, rowScanFunc RowScanFunc) ([]interface{}, error) {
	scaneables := []interface{}{}
	for rows.Next() {
		scanable, err := rowScanFunc(rows)
		if scanable == nil {
			return nil, err
		}
		scaneables = append(scaneables, scanable)
	}
	err := rows.Err()
	if err != nil {
		return nil, err
	}
	return scaneables, nil
}
