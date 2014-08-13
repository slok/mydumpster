package mydumpster

import (
	"fmt"
	"testing"
)

// Unit tests --------------------------
func TestTableCreationStr(t *testing.T) {
}

func TestTableDropStr(t *testing.T) {
	var testPairs = []struct {
		table  string
		result string
	}{
		{"author", "DROP TABLE IF EXISTS `author`;"},
		{"book", "DROP TABLE IF EXISTS `book`;"},
		{"category", "DROP TABLE IF EXISTS `category`;"},
	}

	for _, test := range testPairs {
		v := TableDropStr(test.table)
		if test.result != v {
			t.Error(
				"For", fmt.Sprintf("'%s'", test.table),
				"expected", fmt.Sprintf("'%s'", test.result),
				"got", fmt.Sprintf("'%s'", v),
			)
		}
	}

}

func TestLockTablesStr(t *testing.T) {

	var testPairs = []struct {
		mode   string
		tables []string
		result string
	}{
		{"write", []string{"author"}, "LOCK TABLES `author` WRITE;"},
		{"read", []string{"book"}, "LOCK TABLES `book` READ;"},
		{"write", []string{"book", "author", "category"}, "LOCK TABLES `book` WRITE, `author` WRITE, `category` WRITE;"},
		{"read", []string{"book", "author", "category"}, "LOCK TABLES `book` READ, `author` READ, `category` READ;"},
	}

	for _, test := range testPairs {
		v := LockTablesStr(test.mode, test.tables...)
		if test.result != v {
			t.Error(
				"For", fmt.Sprintf("'%s' - '%s'", test.tables, test.mode),
				"expected", fmt.Sprintf("'%s'", test.result),
				"got", fmt.Sprintf("'%s'", v),
			)
		}
	}
}

//Shit test xD
func TestUnlockTablesStr(t *testing.T) {
	result := "UNLOCK TABLES;"
	if v := UnlockTablesStr(); v != result {
		t.Error(
			"For", fmt.Sprintf("'%s'", ""),
			"expected", fmt.Sprintf("'%s'", result),
			"got", fmt.Sprintf("'%s'", v),
		)
	}
}

func TestInsertRowsStr(t *testing.T) {

	columns := []string{"name", "last_name", "country"}

	var testPairs = []struct {
		rowValues [][]string
		tableName string
		columns   []string
		mode      string
		result    string
	}{
		{
			[][]string{
				[]string{"'Doge'", "'wow'", "'1'"}, []string{"'grumpy'", "'cat'", "'2'"}, []string{"'shibe'", "'much'", "'1'"}},
			"meme",
			columns,
			"",
			"INSERT INTO `meme` (name, last_name, country) VALUES ('Doge', 'wow', '1'), ('grumpy', 'cat', '2'), ('shibe', 'much', '1');",
		},
		{
			[][]string{
				[]string{"'Doge'", "'wow'", "'1'"}, []string{"'grumpy'", "'cat'", "'2'"}, []string{"'shibe'", "'much'", "'1'"}},
			"meme2",
			columns,
			"replace",
			"REPLACE INTO `meme2` (name, last_name, country) VALUES ('Doge', 'wow', '1'), ('grumpy', 'cat', '2'), ('shibe', 'much', '1');",
		},
		{
			[][]string{
				[]string{"'Doge'", "'wow'", "'1'"}, []string{"'grumpy'", "'cat'", "'2'"}, []string{"'shibe'", "'much'", "'1'"}},
			"meme3",
			columns,
			"ignore",
			"INSERT INTO IGNORE `meme3` (name, last_name, country) VALUES ('Doge', 'wow', '1'), ('grumpy', 'cat', '2'), ('shibe', 'much', '1');",
		},
	}

	for _, test := range testPairs {
		v := InsertRowsStr(test.rowValues, test.tableName, test.columns, test.mode)
		if test.result != v {
			t.Error(
				"For", fmt.Sprintf("'%v' - '%s' - '%v' - '%s' - '%s'", test.rowValues, test.tableName, test.columns, test.mode),
				"expected", fmt.Sprintf("'%s'", test.result),
				"got", fmt.Sprintf("'%s'", v),
			)
		}
	}
}

func TestFiltersStr(t *testing.T) {
	var testPairs = []struct {
		filters []string
		result  string
	}{
		{[]string{}, ""},
		{[]string{"id < 10"}, "WHERE id < 10"},
		{[]string{"id < 10", "id > 2", "name like '%something%'"}, "WHERE id < 10 AND id > 2 AND name like '%something%'"},
	}

	for _, test := range testPairs {
		v := FiltersStr(test.filters)
		if test.result != v {
			t.Error(
				"For", fmt.Sprintf("'%v'", test.filters),
				"expected", fmt.Sprintf("'%s'", test.result),
				"got", fmt.Sprintf("'%s'", v),
			)
		}
	}

}

func TestForeignCheckStr(t *testing.T) {

	var testPairs = []struct {
		value  bool
		result string
	}{
		{true, "SET FOREIGN_KEY_CHECKS=1;"},
		{false, "SET FOREIGN_KEY_CHECKS=0;"},
	}

	for _, test := range testPairs {
		v := ForeignCheckStr(test.value)
		if test.result != v {
			t.Error(
				"For", fmt.Sprintf("'%t'", test.value),
				"expected", fmt.Sprintf("'%s'", test.result),
				"got", fmt.Sprintf("'%s'", v),
			)
		}
	}

}

func TestDumpHeaderStr(t *testing.T) {

}

func TestDumpFooterStr(t *testing.T) {
	result := "SET FOREIGN_KEY_CHECKS=1;"
	if v := DumpFooterStr(nil); v != result {
		t.Error(
			"For", fmt.Sprintf("'%t'", nil),
			"expected", fmt.Sprintf("'%s'", result),
			"got", fmt.Sprintf("'%s'", v),
		)

	}
}

func TestShowTablesStr(t *testing.T) {
	result := "SHOW TABLES;"
	if v := ShowTablesStr(); v != result {
		t.Error(
			"expected", fmt.Sprintf("'%s'", result),
			"got", fmt.Sprintf("'%s'", v),
		)

	}
}
