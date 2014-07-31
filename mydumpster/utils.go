package mydumpster

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

const (
	SHOW_TABLE_CREATION_FMT = "SHOW CREATE TABLE %s"
	DROP_TABLE_FMT          = "DROP TABLE IF EXISTS `%s`;"
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

// Checks and error and the program dies (panic)
func CheckKill(e error) {
	if e != nil {
		panic(e)
	}
}
