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

	// Prepare our data
	tableName := "modulo_pago_gateway_set"
	filters := []string{"id = 1"}
	censorships := map[string]mydumpster.Censorship{
		"imagen": mydumpster.Censorship{
			Key:    "imagen",
			Suffix: "_after",
			Prefix: "beafore_",
			Blank:  true,
			Null:   true,
			//DefaultValue: "test",
		},
	}

	table := mydumpster.Table{
		Db:          db,
		TableName:   tableName,
		Filters:     filters,
		Censorships: censorships,
		//Triggers:
	}

	// Do row logic
	table.GetColums()
	channel, err := table.GetRows()
	rows := make([][]string, 0)
	for i := range channel {
		rows = append(rows, i)
	}

	// Some other logic and prepare the dump
	paisDrop := mydumpster.TableDropStr(tableName)
	paisCreation, err := mydumpster.TableCreationStr(db, tableName)
	mydumpster.CheckKill(err)

	insertStr := mydumpster.InsertRowsStr(rows, tableName, table.Columns)
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
