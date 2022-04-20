package dcrpg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bobo-boom/dcrnsugar/db/dbtypes"
	"github.com/bobo-boom/dcrnsugar/db/internal"
	"log"
)

// closeRows closes the input sql.Rows, logging any error.
func closeRows(rows *sql.Rows) {
	if e := rows.Close(); e != nil {

		log.Fatalf("Close of Query failed: %v", e)
	}
}

// SqlExecutor is implemented by both sql.DB and sql.Tx.
type SqlExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// SqlQueryer is implemented by both sql.DB and sql.Tx.
type SqlQueryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// SqlExecQueryer is implemented by both sql.DB and sql.Tx.
type SqlExecQueryer interface {
	SqlExecutor
	SqlQueryer
}

// sqlExec executes the SQL statement string with any optional arguments, and
// returns the number of rows affected.
func sqlExec(db SqlExecutor, stmt, execErrPrefix string, args ...interface{}) (int64, error) {
	res, err := db.Exec(stmt, args...)
	if err != nil {
		return 0, fmt.Errorf("%v: %w", execErrPrefix, err)
	}
	if res == nil {
		return 0, nil
	}

	var N int64
	N, err = res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error in RowsAffected: %w", err)
	}
	return N, err
}

// sqlExecStmt executes the prepared SQL statement with any optional arguments,
// and returns the number of rows affected.
func sqlExecStmt(stmt *sql.Stmt, execErrPrefix string, args ...interface{}) (int64, error) {
	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, fmt.Errorf("%v: %w", execErrPrefix, err)
	}
	if res == nil {
		return 0, nil
	}

	var N int64
	N, err = res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error in RowsAffected: %w", err)
	}
	return N, err
}

// TableExists checks if the specified table exists.
func TableExists(db *sql.DB, tableName string) (bool, error) {
	rows, err := db.Query(`select relname from pg_class where relname = $1`,
		tableName)
	if err != nil {
		return false, err
	}

	defer func() {
		if e := rows.Close(); e != nil {
			log.Fatalf("Close of Query failed: %v", e)
		}
	}()
	return rows.Next(), nil
}

// CreateTable creates a table with the given name using the provided SQL
// statement, if it does not already exist.
func CreateTable(db *sql.DB, tableName, stmt string) error {
	exists, err := TableExists(db, tableName)
	if err != nil {
		return err
	}

	if !exists {
		log.Printf(`Creating the "%s" table.`, tableName)
		_, err = db.Exec(stmt)
		if err != nil {
			return err
		}
	} else {
		//log.Tracef(`Table "%s" exists.`, tableName)
		log.Printf(`Table "%s" exists.`, tableName)
	}

	return err
}

func DropTable(db SqlExecutor, tableName string) error {
	_, err := db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, tableName))
	return err
}

func RetrieveAddresses(ctx context.Context, db *sql.DB, start int64, end int64) (addresses []string, ids []int64, err error) {
	var rows *sql.Rows
	rows, err = db.QueryContext(ctx, internal.SelectAddressRows, start, end)
	if err != nil {
		return nil, nil, err
	}

	defer closeRows(rows)

	for rows.Next() {
		var address string
		var id int64
		err = rows.Scan(&address, &id)
		if err != nil {
			return
		}
		addresses = append(addresses, address)
		ids = append(ids, id)
	}
	err = rows.Err()

	return
}

func RetrieveAddress(ctx context.Context, db *sql.DB, id int64) (addresses []string, ids []int64, err error) {
	var rows *sql.Rows
	rows, err = db.QueryContext(ctx, internal.SelectAddressRow, id)
	if err != nil {
		return nil, nil, err
	}

	defer closeRows(rows)

	for rows.Next() {
		var address string
		var id int64
		err = rows.Scan(&address, &id)
		if err != nil {
			return
		}
		addresses = append(addresses, address)
		ids = append(ids, id)

	}
	err = rows.Err()

	return
}

//// DeleteDuplicateAgendas deletes rows in agendas with duplicate names leaving
//// the one row with the lowest id.
//func DeleteDuplicateAgendas(db *sql.DB) (int64, error) {
//	if isuniq, err := IsUniqueIndex(db, "uix_agendas_name"); err != nil && err != sql.ErrNoRows {
//		return 0, err
//	} else if isuniq {
//		return 0, nil
//	}
//	execErrPrefix := "failed to delete duplicate agendas: "
//	return sqlExec(db, internal.DeleteAgendasDuplicateRows, execErrPrefix)
//}

func InsertAddsBalances(db *sql.DB, balances []*dbtypes.BalanceInfo) error {
	dbtx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin database balance: %w", err)
	}

	stmt, err := dbtx.Prepare(internal.UpsertBalanceRow)
	if err != nil {
		log.Fatalf("Ticket INSERT prepare: %v", err)
		_ = dbtx.Rollback() // try, but we want the Prepare error back
		return err
	}

	for _, b := range balances {

		_, err = stmt.Exec(b.Address, b.Balance, b.Index, b.Flag)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			_ = stmt.Close() // try, but we want the QueryRow error back
			if errRoll := dbtx.Rollback(); errRoll != nil {
				log.Fatalf("Rollback failed: %v", errRoll)
			}
			return err
		}
	}

	// Close prepared statement. Ignore errors as we'll Commit regardless.
	_ = stmt.Close()

	return dbtx.Commit()
}

func InsertAddsBalance(db *sql.DB, balance *dbtypes.BalanceInfo) error {

	stmt, err := db.Prepare(internal.UpsertBalanceRow)
	if err != nil {
		log.Fatalf("Ticket INSERT prepare: %v", err)
		return err
	}
	_, err = stmt.Exec(balance.Address, balance.Balance, balance.Index, balance.Flag)

	return err
}
