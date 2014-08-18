package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/slok/mydumpster/mydumpster"
	"os"
	"sync"
)

const confFileStr = "conf_examples/example.conf"
const confDumpStr = "dump.sql"

var log = mydumpster.GetLogger(mydumpster.MydumpsterLogger)

//TESTING MAIN!!!!!
func main() {

	// Configuration-----------------------------------------------------------
	log.Info("Parsing configuration file...")
	dumpOut := flag.String("output", confDumpStr, "Dump output file")
	confFile := flag.String("config", confFileStr, "Configuration file")

	flag.Parse()

	conf := mydumpster.LoadConfiguration(*confFile)
	dsn := conf.ConnectionStr()
	db, _ := sql.Open("mysql", dsn)
	defer db.Close()

	tables := conf.GetTables(db)

	// Check connection
	log.Info("Checking DB connection...")
	mydumpster.CheckKill(db.Ping())

	// Write to file
	f, err := os.Create(*dumpOut)
	mydumpster.CheckKill(err)
	defer f.Close()

	// SQL dump ----------------------------------------------------------------
	tableList := make([]mydumpster.Table, 0)
	for _, v := range tables {
		tableList = append(tableList, v)
	}

	log.Info(fmt.Sprintf("Start dump process for #%d tables", len(tableList)))
	f.WriteString(mydumpster.DumpHeaderStr(tableList))

	var wg sync.WaitGroup
	tasks := make(chan mydumpster.Table, conf.DumpOptions.Parallel)
	finished := make(chan bool, conf.DumpOptions.Parallel)
	counter := 0
	for _, t := range tables {
		// If will be triggered then don't execute directly and let the
		// triggers do the job
		if t.TriggeredBy == nil {
			counter += 1
			tasks <- t
			wg.Add(1)
			go func() {
				defer wg.Done()
				t := <-tasks
				t.WriteRows(f)
				finished <- true

			}()

			for counter >= conf.DumpOptions.Parallel {
				<-finished // Used to control how many are at the same time
				counter -= 1
			}

		}
	}
	wg.Wait()
	close(tasks)
	close(finished)

	f.WriteString(mydumpster.DumpFooterStr(tableList))
	f.Sync()

	log.Info("Bye bye :)")
}
