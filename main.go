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

//TESTING MAIN!!!!!
func main() {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", dbUser, dbPass, dbHost, dbPort, dbName)

	db, _ := sql.Open("mysql", dsn)
	defer db.Close()

	// Check connection
	mydumpster.CheckKill(db.Ping())

	tableName := "modulo_pago_gateway"
	//filters := []string{"id < 15"}
	censorships := map[string]mydumpster.Censorship{
		"imagen": mydumpster.Censorship{
			Key:    "imagen",
			Suffix: "_after",
			//Prefix:       "beafore_",
			//Blank:        true,
			//Null:         true,
			//DefaultValue: "test",
		},
	}
	setTable := mydumpster.Table{
		Db:        db,
		TableName: "modulo_pago_gateway_set",
	}

	gatewayTable := mydumpster.Table{
		Db:        db,
		TableName: tableName,
		//Filters:     filters,
		Censorships: censorships,
		Triggers: []mydumpster.Trigger{
			mydumpster.Trigger{
				TableDst:      setTable,
				TableSrcName:  "modulo_pago_gateway",
				TableSrcField: "gateway_set_id",
				TableDstField: "id",
				//DumpAll:       true,
			},
			mydumpster.Trigger{
				TableDst:      setTable,
				TableSrcName:  "modulo_pago_gateway",
				TableSrcField: "parent_set_id",
				TableDstField: "id",
			},
		},
	}

	tables := []mydumpster.Table{setTable, gatewayTable}

	//mydumpster.CheckKill(mydumpster.LockTablesRead(db, "pais"))
	//mydumpster.CheckKill(mydumpster.LockTablesWrite(db, "pais"))
	//mydumpster.CheckKill(mydumpster.UnlockTables(db))

	// Write to file
	f, err := os.Create("dump.sql")
	mydumpster.CheckKill(err)
	defer f.Close()

	f.WriteString(mydumpster.DumpHeaderStr(tables))
	gatewayTable.WriteRows(f)
	f.WriteString(mydumpster.DumpFooterStr(tables))

	f.Sync()

}
