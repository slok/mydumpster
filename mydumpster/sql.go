package mydumpster

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	NULL = "NULL"
)

// Locks the tables in read for the current session
func LockTablesRead(db *sql.DB, tableNames ...string) error {
	log.Warning(fmt.Sprintf("Locking tables '%v' in READ mode", tableNames))
	_, err := db.Exec(LockTablesStr("read", tableNames...))
	if err != nil {
		log.Error(fmt.Sprintf("Error locking tables '%v' in READ mode", tableNames))
	}
	return err
}

// Locks the tables in write for the current session
func LockTablesWrite(db *sql.DB, tableNames ...string) error {
	log.Warning(fmt.Sprintf("Locking tables '%v' in WRITE mode", tableNames))
	_, err := db.Exec(LockTablesStr("write", tableNames...))
	if err != nil {
		log.Error(fmt.Sprintf("Error locking tables '%v' in WRITE mode", tableNames))
	}
	return err
}

func UnlockTables(db *sql.DB) error {
	log.Warning(fmt.Sprintf("Unlocking all tables"))
	_, err := db.Exec(UnlockTablesStr())
	if err != nil {
		log.Error(fmt.Sprintf("Error unlocking all tables"))
	}
	return err
}

func GetTableNames(db *sql.DB) ([]string, error) {
	log.Debug(fmt.Sprintf("Getting table names"))

	rows, err := db.Query(ShowTablesStr())
	if err != nil {
		log.Error(fmt.Sprintf("Error Executing query"))
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
