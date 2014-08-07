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

	wheres := ""
	if filters != nil {
		wheres = fmt.Sprintf(WHERE_FMT, strings.Join(filters, AND_FMT))
	}

	selectStr := fmt.Sprintf(GET_ROWS_FMT, columnStr, tableName, wheres)
	// Create the channel
	channel := make(chan []string)

	rows, err := db.Query(selectStr)
	// This will make the thing lazy
	go func() {
		defer rows.Close()
		// For each row...
		i := 0
		for rows.Next() {
			i = i + 1

			// Create the slice to save the rawbytes
			scanArgs := make([]interface{}, len(columns))
			scanArgsCopy := make([]string, len(columns))

			// Initialize our "abstract" list
			for i := range columns { // use columns as a lenth loop only
				scanArgs[i] = new(sql.RawBytes)
			}

			//FIXME: for now channels don't send errors
			err = rows.Scan(scanArgs...)

			for i, v := range scanArgs {
				if v != nil {
					// Scape before surrounding by ''(apostrophes)
					scapedString := ReplaceCharacters(
						fmt.Sprintf("%s", *(v.(*sql.RawBytes))))
					scanArgsCopy[i] = fmt.Sprintf("'%s'", scapedString)

				} else {
					scanArgsCopy[i] = "NULL"
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
