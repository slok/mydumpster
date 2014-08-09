package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/slok/mydumpster/mydumpster"
	"os"
)

const confFile = "conf.json.example"

//TESTING MAIN!!!!!
func main() {

	conf := mydumpster.LoadConfiguration(confFile)
	dsn := conf.ConnectionStr()
	db, _ := sql.Open("mysql", dsn)
	defer db.Close()

	tables := conf.GetTables(db)

	// Check connection
	mydumpster.CheckKill(db.Ping())

	//mydumpster.CheckKill(mydumpster.LockTablesRead(db, "pais"))
	//mydumpster.CheckKill(mydumpster.LockTablesWrite(db, "pais"))
	//mydumpster.CheckKill(mydumpster.UnlockTables(db))

	// Write to file
	f, err := os.Create("dump.sql")
	mydumpster.CheckKill(err)
	defer f.Close()

	tableList := make([]mydumpster.Table, 0)
	for _, v := range tables {
		tableList = append(tableList, v)
	}

	f.WriteString(mydumpster.DumpHeaderStr(tableList))
	gw := tables["modulo_pago_gateway"]
	gw.WriteRows(f)
	f.WriteString(mydumpster.DumpFooterStr(tableList))

	f.Sync()

}
