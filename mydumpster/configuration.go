package mydumpster

import (
	"database/sql"
	"encoding/json"
	//"errors"
	"fmt"
	"os"
)

type ConfTrigger struct {
	TableDstName string `json:"dst_table_name"`
	FieldSrcName string `json:"src_field_name"`
	FieldDstName string `json:"dst_field_name"`
}

type ConfCensorship struct {
	Prefix  string `json:"prefix"`
	Suffix  string `json:"suffix"`
	Blank   bool   `json:"blank"`
	Null    bool   `json:"null"`
	Default string `json:"default"`
}

type ConfTable struct {
	Filters    []string                   `json:"filters"`
	Censorship map[string]*ConfCensorship `json:"censorship"`
	Triggers   []*ConfTrigger             `json:"triggers"`
	Exclude    bool                       `json:"exclude"`
	DumpAll    bool                       `json:"dump_all"`
}

type ConfDatabase struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Db       string `json:"db"`
}

type ConfDump struct {
	AllTables bool `json:"all_tables"`
	Parallel  int  `json:"parallel"`
}

type Configuration struct {
	Tables      map[string]*ConfTable `json:"tables"`
	Database    *ConfDatabase         `json:"database"`
	DumpOptions *ConfDump             `json:"dump"`
}

func LoadConfiguration(filePath string) *Configuration {

	file, err := os.Open(filePath)
	CheckKill(err)

	decoder := json.NewDecoder(file)
	configuration := Configuration{}

	err = decoder.Decode(&configuration)
	CheckKill(err)
	//configuration.PrintConfiguration()

	return &configuration
}

func (c *Configuration) ConnectionStr() string {

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		c.Database.User, c.Database.Password,
		c.Database.Host, c.Database.Port,
		c.Database.Db)
}

//FIXME:
//  - If any table dump all then dont create next triggers(Optimization)
func (c *Configuration) GetTables(db *sql.DB) map[string]Table {

	// Create the containers
	var tables = make(map[string]Table)

	// If 'all_tables' is activated then we need all the tables from the database
	if c.DumpOptions.AllTables {
		tableNames, _ := GetTableNames(db)
		for _, i := range tableNames {
			// Exclude tables
			if !c.Tables[i].Exclude {
				tables[i] = Table{
					Db:          db,
					TableName:   i,
					Filters:     make([]string, 0),
					Censorships: make(map[string]Censorship),
					Triggers:    make([]*Trigger, 0),
					TriggeredBy: nil,
					DumpAll:     false,
				}
			}
		}
	}

	// Init our map (We need al the table structs created so
	// we can refer from a table to another while we parse)
	for k, v := range c.Tables {
		// Exclude table
		if !v.Exclude {
			tables[k] = Table{
				Db:          db,
				TableName:   k,
				Filters:     make([]string, len(v.Filters)),
				Censorships: make(map[string]Censorship),
				Triggers:    make([]*Trigger, len(v.Triggers)),
				TriggeredBy: nil,
				DumpAll:     v.DumpAll,
			}
		} else {
			log.Warning("Excludig '%s' table and pointing triggers ", k)
		}
	}

	for k, v := range c.Tables {
		// Exclude table
		if !v.Exclude {

			t := tables[k]

			// Create the filters
			for k, f := range v.Filters {
				t.Filters[k] = f
			}

			// Create the censorships
			if v.Censorship != nil {
				for ck, cv := range v.Censorship {
					t.Censorships[ck] = Censorship{
						Key:          ck,
						Suffix:       cv.Suffix,
						Prefix:       cv.Prefix,
						Blank:        cv.Blank,
						Null:         cv.Null,
						DefaultValue: cv.Default,
					}
				}
			}

			// Create the triggers
			if v.Triggers != nil {
				for tk, tv := range v.Triggers {

					// Check if triggers an excluded table
					// (if the table doesn't exist thn there isn't configuration of exclude)
					auxConfTabl, ok := c.Tables[tv.TableDstName]
					if !ok || ok && !auxConfTabl.Exclude {

						aux, ok := tables[tv.TableDstName]

						// Table not declared (We create)
						if !ok {
							//CheckKill(errors.New("Not table in map"))
							tables[tv.TableDstName] = Table{
								Db:          db,
								TableName:   tv.TableDstName,
								Filters:     make([]string, 0),
								Censorships: make(map[string]Censorship),
								Triggers:    make([]*Trigger, 0),
								TriggeredBy: &t,
								DumpAll:     false,
							}
							aux = tables[tv.TableDstName]
						}

						t.Triggers[tk] = &Trigger{
							TableDst:      aux,
							TableSrcName:  k,
							TableSrcField: tv.FieldSrcName,
							TableDstField: tv.FieldDstName,
						}
					} else { // Remove the element
						// Difference between slice 1 and 2 applied to the key
						//auxTk := tk - (len(v.Triggers) - len(t.Triggers))
						//t.Triggers = append(t.Triggers[:auxTk], t.Triggers[auxTk+1:]...)
						t.Triggers[tk] = nil
					}
				}
			}
			fmt.Println(len(t.Triggers))
		}
	}
	return tables
}

func (c *Configuration) PrintConfiguration() {
	d := c.Database
	fmt.Println("Database")
	fmt.Println("-----------")
	fmt.Println(fmt.Sprintf("  -Host: %s", d.Host))
	fmt.Println(fmt.Sprintf("  -Port: %d", d.Port))
	fmt.Println(fmt.Sprintf("  -Passwords: %s", d.Password))
	fmt.Println(fmt.Sprintf("  -User: %s", d.User))
	fmt.Println(fmt.Sprintf("  -Db: %s", d.Db))

	fmt.Println("")

	do := c.DumpOptions
	fmt.Println("Dump options")
	fmt.Println("-------------")
	fmt.Println(fmt.Sprintf("  - All tables: %t", do.AllTables))
	fmt.Println(fmt.Sprintf("  - Parallel: %d", do.Parallel))

	fmt.Println("")

	for k, v := range c.Tables {
		fmt.Println("Table " + k)
		fmt.Println("-----------")
		fmt.Println(fmt.Sprintf("Exclude: %t", v.Exclude))
		fmt.Println(fmt.Sprintf("Dump all: %t", v.DumpAll))

		fmt.Println("Filters:")
		for _, f := range v.Filters {
			fmt.Println("  -" + f)
		}

		fmt.Println("Censore:")
		for k2, v2 := range v.Censorship {
			fmt.Println("  -" + k2)
			fmt.Println(fmt.Sprintf("    -Prefix: %s", v2.Prefix))
			fmt.Println(fmt.Sprintf("    -Suffix: %s", v2.Suffix))
			fmt.Println(fmt.Sprintf("    -Blank: %t", v2.Blank))
			fmt.Println(fmt.Sprintf("    -Null: %t", v2.Null))
			fmt.Println(fmt.Sprintf("    -Default: %s", v2.Default))

		}

		fmt.Println("Triggers:")
		for _, v3 := range v.Triggers {
			fmt.Println(fmt.Sprintf("  -Src field name: %s", v3.FieldSrcName))
			fmt.Println(fmt.Sprintf("  -Dst field name: %s", v3.FieldDstName))
			fmt.Println(fmt.Sprintf("  -Dst Table name: %s", v3.TableDstName))
			fmt.Println("")
		}

		fmt.Println("=============================")
	}
}
