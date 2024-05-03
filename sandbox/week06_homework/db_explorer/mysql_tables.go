package main

import (
	"fmt"
	"strings"
)

type ColumnType interface {
	NewVar() interface{}
	IsValidValue(val interface{}) bool
}

// ColumnType implementation
type IntColumn struct {
	Null bool
}

// ColumnType implementation
type StringColumn struct {
	Null bool
}

func (c IntColumn) NewVar() interface{} {
	if c.Null {
		return new(*int64)
	} else {
		return new(int64)
	}
}
func (c IntColumn) IsValidValue(val interface{}) bool {
	if val == nil {
		return c.Null
	}

	_, ok := val.(int64)
	return ok
}

func (c StringColumn) NewVar() interface{} {
	if c.Null {
		return new(*string)
	} else {
		return new(string)
	}
}
func (c StringColumn) IsValidValue(val interface{}) bool {
	if val == nil {
		return c.Null
	}

	_, ok := val.(string)
	return ok
}

type TableColumn struct {
	Field      string
	Type       ColumnType
	Collation  interface{}
	Null       bool
	Key        string
	Default    interface{}
	Extra      string
	Privileges string
	Comment    string
}

type Table struct {
	Name    string
	Pk      string
	Columns []TableColumn
}

// row as list of values, types known in runtime
type TableRow []interface{}

// row as map field:value
type TableRecord map[string]interface{}

func (t Table) NewRow() TableRow {
	row := make(TableRow, len(t.Columns))
	for i := range row {
		row[i] = t.Columns[i].Type.NewVar() // key moment: value of dynamic type (int, *int, string, *string, ...)
	}
	return row
}

func (t Table) NewRecord(row TableRow) TableRecord {
	record := TableRecord{}
	for i, c := range t.Columns {
		record[c.Field] = row[i]
	}
	return record
}

func (srv MysqlExplorerHttpHandlers) GetTables() (map[string]Table, error) {
	tableNameList, err := srv.GetTableNames()
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %s", err)
	}

	tables := make(map[string]Table, len(tableNameList))
	for _, tname := range tableNameList {
		columns, err := srv.GetTableColumns(tname)
		if err != nil {
			return nil, fmt.Errorf("failed to get table's %s columns: %s", tname, err)
		}

		table := Table{
			Name:    tname,
			Columns: columns,
		}

		for _, col := range columns {
			if col.Key == "PRI" {
				table.Pk = col.Field
				break
			}
		}

		tables[tname] = table
	}

	return tables, nil
}

func (srv MysqlExplorerHttpHandlers) GetTableNames() (tables []string, err error) {
	rows, err := srv.DB.Query("SHOW TABLES")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch table names: %s", err)
	}
	defer rows.Close()

	tables = make([]string, 0, 16)
	var t string
	for rows.Next() {
		rows.Scan(&t)
		tables = append(tables, t)
	}

	return tables, nil
}

func (srv MysqlExplorerHttpHandlers) GetTableColumns(table string) (columns []TableColumn, err error) {
	rows, err := srv.DB.Query("SHOW FULL COLUMNS FROM " + table)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch columns for table '%s': %s", table, err)
	}
	defer rows.Close()

	var (
		colType   string
		sNullable string
		bNullable bool
	)

	columns = make([]TableColumn, 0, 16)
	for rows.Next() {
		col := TableColumn{}
		rows.Scan(
			&col.Field,
			&colType,
			&col.Collation,
			&sNullable,
			&col.Key,
			&col.Default,
			&col.Extra,
			&col.Privileges,
			&col.Comment,
		)
		bNullable = (sNullable == "YES")

		if strings.Contains(colType, "int") {
			col.Type = IntColumn{bNullable}
		} else {
			col.Type = StringColumn{bNullable}
		}

		columns = append(columns, col)
	}

	return columns, nil
}
