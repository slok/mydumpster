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
	Triggers    []*Trigger
	TriggeredBy *Table
	DumpAll     bool
}

// Loads the column data of the table
func (t *Table) GetColums() error {
	log.Debug(fmt.Sprintf("Get '%s' table columns", t.TableName))
	rows, err := t.Db.Query(fmt.Sprintf(GET_ONE_ROW_FMT, t.TableName))
	if err != nil {
		log.Error(fmt.Sprintf("Error getting '%s' table columns", t.TableName))
		return err
	}
	defer rows.Close()

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

	log.Debug(fmt.Sprintf("Start getting '%s' table rows...", t.TableName))
	// Create the select string
	columnStr := strings.Join(t.Columns, ", ")

	// Apply wheres if needed
	log.Debug(fmt.Sprintf("Applying filters..."))
	wheres := ""
	if t.Filters != nil && len(t.Filters) > 0 {
		wheres = FiltersStr(t.Filters)
	}
	selectStr := fmt.Sprintf(GET_ROWS_FMT, columnStr, t.TableName, wheres)
	rows, err := t.Db.Query(selectStr)

	if err != nil {
		log.Error(fmt.Sprintf("Error executing query for table rows"))
	}

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
			if err != nil {
				log.Error(fmt.Sprintf("Error Scanning args"))
			}

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
						//log.Debug(fmt.Sprintf("Censoring '%s' field from '%s' table", t.Columns[i], t.TableName))
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
	log.Info(fmt.Sprintf("Start Dump process for table '%s'...", t.TableName))
	// First lock
	CheckKill(LockTablesRead(t.Db, t.TableName))

	// Do row logic
	t.GetColums()
	channel, err := t.getRows()
	rows := make([][]string, 0)

	for i := range channel {
		rows = append(rows, i)
	}

	insertStr := InsertRowsStr(rows, t.TableName, t.Columns, "replace")

	// Get triggers (For now one level)
	for _, tr := range t.Triggers {
		if tr != nil {
			// Only get the ids of te arent related rows, so we set this as a filter
			if !tr.TableDst.DumpAll {
				tr.TableDst.Filters = append(
					tr.TableDst.Filters, tr.SelectQueryFromRowsStr(rows, t.Columns))
			} else {
				log.Warning(fmt.Sprintf("'%s' table will be totally dumped", t.TableName))
			}
			log.Debug(fmt.Sprintf("'%s' table triggered '%s' table dump", t.TableName, tr.TableDst.TableName))
			tr.TableDst.WriteRows(w)
		}
	}

	// Save in the file
	t.WriteTableHeader(w)
	fmt.Fprintln(w, insertStr)
	t.WriteTableFooter(w)
	CheckKill(UnlockTables(t.Db))
	log.Info(fmt.Sprintf("Finish Dump process for table '%s'...", t.TableName))
	return err
}

func (t *Table) WriteTableHeader(w io.Writer) {
	log.Debug(fmt.Sprintf("Writting table header for table '%s'...", t.TableName))
	fmt.Fprintln(w,
		sqlComment(fmt.Sprintf("Start table `%s` dump", t.TableName)))
	fmt.Fprintln(w, LockTablesStr("write", t.TableName))
}

func (t *Table) WriteTableFooter(w io.Writer) {
	log.Debug(fmt.Sprintf("Writting table footer for table '%s'...", t.TableName))
	fmt.Fprintln(w, UnlockTablesStr())
	fmt.Fprintln(w,
		sqlComment(fmt.Sprintf("End table `%s` dump", t.TableName)))
}
