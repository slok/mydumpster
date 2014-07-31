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
	LOCK_TABLE_FMT          = "LOCK TABLES %s"
	LOCK_READ_FMT           = "`%s` READ"
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

// Locks the tables for the current session
func LockTables(db *sql.DB, tableNames ...string) error {

	aux := make([]string, 0)
	// Create the table locks
	for _, tn := range tableNames {
		aux = append(aux, fmt.Sprintf(LOCK_READ_FMT, tn))
	}

	query := fmt.Sprintf(LOCK_TABLE_FMT, strings.Join(aux, ", "))

	// apply
	_, err := db.Exec(query)
	return err
}

// Checks and error and the program dies (panic)
func CheckKill(e error) {
	if e != nil {
		panic(e)
	}
}
