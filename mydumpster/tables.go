package mydumpster

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"strings"
)

type Table struct {
	Db          *sql.DB
	TableName   string
	Filters     []string
	Columns     []string
	Censorships map[string]Censorship
	//Triggers    []Trigger
}

// Loads the column data of the table
func (t *Table) GetColums() error {

	rows, err := t.Db.Query(fmt.Sprintf(GET_ONE_ROW_FMT, t.TableName))
	if err != nil {
		return err
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	// Store the colume names in the list
	vals := make([]string, len(cols))
	for i, col := range cols {
		vals[i] = col
	}
	t.Columns = vals

	return err
}

// Gets the rows of a table censored if neccesary
func (t *Table) getRows() (chan []string, error) {

	// Create the select string
	columnStr := strings.Join(t.Columns, ", ")

	// Apply wheres if needed
	wheres := ""
	if t.Filters != nil {
		wheres = filtersStr(t.Filters)
	}
	selectStr := fmt.Sprintf(GET_ROWS_FMT, columnStr, t.TableName, wheres)

	rows, err := t.Db.Query(selectStr)

	// Create the channel to be lazy
	channel := make(chan []string)
	go func() {
		defer rows.Close()
		// For each row...
		for rows.Next() {
			// Create the slice to save the rawbytes
			scanArgs := make([]interface{}, len(t.Columns))
			scanArgsCopy := make([]string, len(t.Columns))

			// Initialize our "abstract" list
			for i := range t.Columns { // use columns as a lenth loop only
				scanArgs[i] = new(sql.NullString)
			}

			//FIXME: for now channels don't send errors
			err = rows.Scan(scanArgs...)
			var argValue sql.NullString

			for i, v := range scanArgs {
				argValue = (*(v.(*sql.NullString)))

				setToNull := !argValue.Valid

				// Check if is NULL before doing anything
				if !setToNull {
					// Scape before surrounding by ''(apostrophes)
					scapedString := ReplaceCharacters(
						fmt.Sprintf("%s", argValue.String))

					// Censore the string only if necessary
					censoreship, ok := t.Censorships[t.Columns[i]]
					if ok {
						scapedString, setToNull = censoreship.censore(scapedString)
					}
					scanArgsCopy[i] = fmt.Sprintf("'%s'", scapedString)
				}

				// Use this style instead of else because the censor could set
				// to NULL after entering in the string logic
				if setToNull {
					scanArgsCopy[i] = NULL
				}
			}

			// Finished, so send lazily
			channel <- scanArgsCopy
		}
		// We are done here
		close(channel)
	}()
	return channel, err
}

// Gets a table (and its triggers) and writes to the writer passed
func (t *Table) WriteRows(w io.Writer) error {

	// Do row logic
	t.GetColums()
	channel, err := t.getRows()
	rows := make([][]string, 0)

	for i := range channel {
		rows = append(rows, i)
	}
	insertStr := InsertRowsStr(rows, t.TableName, t.Columns)

	// Get the triggers

	fmt.Fprint(w, insertStr)
	return err
}
