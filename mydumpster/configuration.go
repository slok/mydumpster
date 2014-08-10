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
	DumpAll      bool   `json:"dump_all"`
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
}

type ConfDatabase struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Db       string `json:"db"`
}

type Configuration struct {
	Tables   map[string]*ConfTable `json:"tables"`
	Database *ConfDatabase         `json:"database"`
}

func LoadConfiguration(filePath string) *Configuration {

	file, err := os.Open(filePath)
	CheckKill(err)

	decoder := json.NewDecoder(file)
	configuration := Configuration{}

	err = decoder.Decode(&configuration)
	CheckKill(err)

	return &configuration
}

func (c *Configuration) ConnectionStr() string {

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		c.Database.User, c.Database.Password,
		c.Database.Host, c.Database.Port,
		c.Database.Db)
}

//FIXME:
//  - If any table dump all then create next triggers(Optimization)
func (c *Configuration) GetTables(db *sql.DB) map[string]Table {

	// Create the containers
	var tables = make(map[string]Table)

	// Init our map (We need al the table structs created so
	// we can refer from a table to another while we parse)
	for k, v := range c.Tables {
		tables[k] = Table{
			Db:          db,
			TableName:   k,
			Filters:     make([]string, len(v.Filters)),
			Censorships: make(map[string]Censorship),
			Triggers:    make([]Trigger, len(v.Triggers)),
		}
	}

	for k, v := range c.Tables {
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

		// Create the censorships
		if v.Triggers != nil {
			for tk, tv := range v.Triggers {
				aux, ok := tables[tv.TableDstName]

				// Table not declared (We create)
				if !ok {
					//CheckKill(errors.New("Not table in map"))
					tables[tv.TableDstName] = Table{
						Db:          db,
						TableName:   tv.TableDstName,
						Filters:     make([]string, 0),
						Censorships: make(map[string]Censorship),
						Triggers:    make([]Trigger, 0),
					}
					aux = tables[tv.TableDstName]
				}

				t.Triggers[tk] = Trigger{
					TableDst:      aux,
					TableSrcName:  k,
					TableSrcField: tv.FieldSrcName,
					TableDstField: tv.FieldDstName,
					DumpAll:       tv.DumpAll,
				}
			}
		}

	}
	return tables
}

func (c *Configuration) PrintConfiguration() {
	v0 := c.Database
	fmt.Println("Database ")
	fmt.Println("-----------")
	fmt.Println(fmt.Sprintf("  -Host: %s", v0.Host))
	fmt.Println(fmt.Sprintf("  -Port: %d", v0.Port))
	fmt.Println(fmt.Sprintf("  -Passwords: %s", v0.Password))
	fmt.Println(fmt.Sprintf("  -User: %s", v0.User))
	fmt.Println(fmt.Sprintf("  -Db: %s", v0.Db))

	fmt.Println("")

	for k, v := range c.Tables {
		fmt.Println("Table " + k)
		fmt.Println("-----------")
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
			fmt.Println(fmt.Sprintf("  -Dump all: %t", v3.DumpAll))
			fmt.Println("")
		}

		fmt.Println("=============================")
	}
}
