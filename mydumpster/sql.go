package mydumpster

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

const (
	WHERE_FMT = "WHERE %s"
	AND_FMT   = " AND "
	NULL      = "NULL"
)

// Locks the tables in read for the current session
func LockTablesRead(db *sql.DB, tableNames ...string) error {
	_, err := db.Exec(LockTablesStr("read", tableNames...))
	return err
}

// Locks the tables in write for the current session
func LockTablesWrite(db *sql.DB, tableNames ...string) error {
	_, err := db.Exec(LockTablesStr("write", tableNames...))
	return err
}

func UnlockTables(db *sql.DB) error {
	_, err := db.Exec(UnlockTablesStr())
	return err
}

// FIXME: Change in the future to be lazy
func GetRows(db *sql.DB, tableName string, columns []string, filters []string) (chan []string, error) {

	// Create the select string
	columnStr := strings.Join(columns, ", ")

	// Apply wheres if needed
	wheres := ""
	if filters != nil {
		wheres = fmt.Sprintf(WHERE_FMT, strings.Join(filters, AND_FMT))
	}
	selectStr := fmt.Sprintf(GET_ROWS_FMT, columnStr, tableName, wheres)

	rows, err := db.Query(selectStr)

	// Create the channel to be lazy
	channel := make(chan []string)
	go func() {
		defer rows.Close()
		// For each row...
		for rows.Next() {
			// Create the slice to save the rawbytes
			scanArgs := make([]interface{}, len(columns))
			scanArgsCopy := make([]string, len(columns))

			// Initialize our "abstract" list
			for i := range columns { // use columns as a lenth loop only
				scanArgs[i] = new(sql.NullString)
			}

			//FIXME: for now channels don't send errors
			err = rows.Scan(scanArgs...)
			var argValue sql.NullString
			for i, v := range scanArgs {
				argValue = (*(v.(*sql.NullString)))

				// Check if is NULL before doing anything
				if argValue.Valid {
					// Scape before surrounding by ''(apostrophes)
					scapedString := ReplaceCharacters(
						fmt.Sprintf("%s", argValue.String))
					scanArgsCopy[i] = fmt.Sprintf("'%s'", scapedString)

				} else {
					scanArgsCopy[i] = NULL
				}
			}
			// Send lazily
			channel <- scanArgsCopy
		}
		// We are done here
		close(channel)
	}()
	return channel, err
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
