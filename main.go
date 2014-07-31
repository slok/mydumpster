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

	pais_drop := mydumpster.GetTableDrop("pais")
	pais_creation, err := mydumpster.GetTableCreation(db, "pais")
	mydumpster.CheckKill(err)

	mydumpster.CheckKill(mydumpster.LockTablesRead(db, "pais"))
	mydumpster.CheckKill(mydumpster.LockTablesWrite(db, "pais"))
	mydumpster.CheckKill(mydumpster.UnlockTables(db))
	// Write to file
	f, err := os.Create("dump.sql")
	mydumpster.CheckKill(err)
	defer f.Close()

	f.WriteString(pais_drop)
	f.WriteString("\n")
	f.WriteString(pais_creation)
	f.Sync()

}
