package mydumpster

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	NULL = "NULL"
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

func GetTableNames(db *sql.DB) ([]string, error) {

	rows, err := db.Query(ShowTablesStr())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vals := make([]string, 0)

	for rows.Next() {
		var val string
		err = rows.Scan(&val)
		vals = append(vals, val)
	}
	return vals, err
}
