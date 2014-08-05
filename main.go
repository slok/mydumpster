package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/slok/mydumpster/mydumpster"
	"os"
)

const (
	dbUser = "root"
	dbPass = ""
	dbHost = "172.17.0.2"
	dbPort = 3306
	dbName = "ticketbis_dev"
)

func main() {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", dbUser, dbPass, dbHost, dbPort, dbName)

	db, _ := sql.Open("mysql", dsn)
	defer db.Close()

	// Check connection
	mydumpster.CheckKill(db.Ping())

	table := "historico"
	//filters := []string{"id >= 1", "id < 30"}

	paisDrop := mydumpster.GetTableDrop(table)
	paisCreation, err := mydumpster.GetTableCreation(db, table)
	mydumpster.CheckKill(err)

	columns, err := mydumpster.GetColums(db, table)
	rows, err := mydumpster.GetRows(db, table, columns, nil)
	insertStr := mydumpster.GetInsertStrFromRows(rows, table, columns)
	mydumpster.CheckKill(err)

	//mydumpster.CheckKill(mydumpster.LockTablesRead(db, "pais"))
	//mydumpster.CheckKill(mydumpster.LockTablesWrite(db, "pais"))
	//mydumpster.CheckKill(mydumpster.UnlockTables(db))

	// Write to file
	f, err := os.Create("dump.sql")
	mydumpster.CheckKill(err)
	defer f.Close()

	f.WriteString(paisDrop)
	f.WriteString("\n")
	f.WriteString(paisCreation)
	f.WriteString("\n")
	f.WriteString(insertStr)
	f.Sync()

}
