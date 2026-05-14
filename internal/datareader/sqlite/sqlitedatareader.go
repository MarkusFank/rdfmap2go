package sqlite

import (
	"database/sql"

	"github.com/MarkusFank/rdfmap2go/internal/datareader"
	_ "modernc.org/sqlite"
)

type SqliteDataReader struct {
	sqliteFile    string
	isInitialized bool
	db            *sql.DB
	sqlRows       *sql.Rows
	columns       *[]string
	columnTypes   []*sql.ColumnType
}

func (reader *SqliteDataReader) Init(file string, query string) error {

	db, err := sql.Open("sqlite", file)

	if err != nil {
		return err
	}

	reader.db = db

	rows, err := db.Query(query)

	if err != nil {
		return nil
	}

	reader.sqlRows = rows

	columns, err := reader.sqlRows.Columns()

	if err != nil {
		return err
	}

	columnTypes, err := reader.sqlRows.ColumnTypes()

	if err != nil {
		return err
	}

	reader.columns = &columns
	reader.columnTypes = columnTypes

	return nil
}

func (reader *SqliteDataReader) ReadRow() (*datareader.DataRow, error) {
	hasRow := reader.sqlRows.Next()

	if !hasRow {
		return nil, nil
	}

	values := make([]any, len(*reader.columns))
	valuePtrs := make([]any, len(*reader.columns))

	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if err := reader.sqlRows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	dataRow := datareader.DataRow{}
	for i, col := range *reader.columns {
		val := values[i]

		// TODO better type checking (via "columnTypes")
		if s, isString := val.(string); isString {
			dataRow[col] = s
		} else {
			if intVal, isInt := val.(int64); isInt {
				dataRow[col] = intVal
			}

			// TODO handler other types
		}
	}

	return &dataRow, nil
}

func (reader *SqliteDataReader) Close() {
	reader.sqlRows.Close()
	reader.db.Close()
}
