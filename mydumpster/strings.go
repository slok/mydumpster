package mydumpster

import (
	"database/sql"
	"fmt"
	"strings"
)

const (
	SHOW_TABLE_CREATION_FMT = "SHOW CREATE TABLE %s;"
	DROP_TABLE_FMT          = "DROP TABLE IF EXISTS `%s`;"
	LOCK_TABLE_FMT          = "LOCK TABLES %s;"
	LOCK_READ_FMT           = "`%s` READ"
	LOCK_WRITE_FMT          = "`%s` WRITE"
	UNLOCK_TABLES_FMT       = "UNLOCK TABLES;"
	GET_ONE_ROW_FMT         = "SELECT * FROM %s LIMIT 1;"
	GET_ROWS_FMT            = "SELECT %s from `%s` %s;"
	INSERT_FMT              = "INSERT INTO `%s` (%s) VALUES %s;"
	WHERE_FMT               = "WHERE %s"
	AND_FMT                 = " AND "
)

// Returns the table creanion syntax string
func TableCreationStr(db *sql.DB, tableName string) (string, error) {
	var garbage, result string
	err := db.QueryRow(fmt.Sprintf(
		SHOW_TABLE_CREATION_FMT, tableName)).Scan(&garbage, &result)

	return result + ";", err
}

// Returns the table creanion syntax string
func TableDropStr(tableName string) string {
	return fmt.Sprintf(DROP_TABLE_FMT, tableName)
}

func LockTablesStr(mode string, tableNames ...string) string {

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

func UnlockTablesStr() string {
	return UNLOCK_TABLES_FMT
}

func InsertRowsStr(rowValues [][]string, tableName string, columns []string) string {

	columnStr := strings.Join(columns, ", ")
	strRows := make([]string, 0)

	for _, values := range rowValues {
		strRows = append(strRows, fmt.Sprintf("(%s)", strings.Join(values, ", ")))
	}

	return fmt.Sprintf(INSERT_FMT, tableName, columnStr, strings.Join(strRows, ", "))
}

func filtersStr(filters []string) string {
	return fmt.Sprintf(WHERE_FMT, strings.Join(filters, AND_FMT))

}
