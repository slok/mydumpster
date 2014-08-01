package mydumpster

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

const (
	SHOW_TABLE_CREATION_FMT = "SHOW CREATE TABLE %s"
	DROP_TABLE_FMT          = "DROP TABLE IF EXISTS `%s`;"
	LOCK_TABLE_FMT          = "LOCK TABLES %s;"
	LOCK_READ_FMT           = "`%s` READ"
	LOCK_WRITE_FMT          = "`%s` WRITE"
	UNLOCK_TABLES_FMT       = "UNLOCK TABLES;"
	GET_ONE_ROW_FMT         = "SELECT * FROM %s LIMIT 1;"
)

// Returns the table creanion syntax string
func GetTableCreation(db *sql.DB, tableName string) (string, error) {
	var garbage, result string
	err := db.QueryRow(fmt.Sprintf(
		SHOW_TABLE_CREATION_FMT, tableName)).Scan(&garbage, &result)
	return result, err
}

// Returns the table creanion syntax string
func GetTableDrop(tableName string) string {
	return fmt.Sprintf(DROP_TABLE_FMT, tableName)
}

func GetLockTables(mode string, tableNames ...string) string {

	// default READ
	if mode == "write" {
		mode = LOCK_WRITE_FMT
	} else {
		mode = LOCK_READ_FMT
	}

	aux := make([]string, 0)
	// Create the table locks
	for _, tn := range tableNames {
		aux = append(aux, fmt.Sprintf(mode, tn))
	}

	return fmt.Sprintf(LOCK_TABLE_FMT, strings.Join(aux, ", "))
}

func GetUnlockTables() string {
	return UNLOCK_TABLES_FMT
}

// Locks the tables in read for the current session
func LockTablesRead(db *sql.DB, tableNames ...string) error {
	_, err := db.Exec(GetLockTables("read", tableNames...))
	return err
}

// Locks the tables in write for the current session
func LockTablesWrite(db *sql.DB, tableNames ...string) error {
	_, err := db.Exec(GetLockTables("write", tableNames...))
	return err
}

func UnlockTables(db *sql.DB) error {
	_, err := db.Exec(GetUnlockTables())
	return err
}

func GetColums(db *sql.DB, tableName string) ([]string, error) {

	rows, err := db.Query(fmt.Sprintf(GET_ONE_ROW_FMT, tableName))
	if err != nil {
		return nil, err
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Store the colume names in the list
	vals := make([]string, len(cols))
	for i, col := range cols {
		vals[i] = col
	}

	return vals, err
}

// Checks and error and the program dies (panic)
func CheckKill(e error) {
	if e != nil {
		panic(e)
	}
}
