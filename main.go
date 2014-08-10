package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/slok/mydumpster/mydumpster"
	"os"
)

const confFileStr = "conf.json.example"
const confDumpStr = "dump.sql"

//TESTING MAIN!!!!!
func main() {

	dumpOut := flag.String("output", confDumpStr, "Dump output file")
	confFile := flag.String("config", confFileStr, "Configuration file")

	flag.Parse()

	conf := mydumpster.LoadConfiguration(*confFile)
	dsn := conf.ConnectionStr()
	db, _ := sql.Open("mysql", dsn)
	defer db.Close()

	tables := conf.GetTables(db)

	// Check connection
	mydumpster.CheckKill(db.Ping())

	// Write to file
	f, err := os.Create(*dumpOut)
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
