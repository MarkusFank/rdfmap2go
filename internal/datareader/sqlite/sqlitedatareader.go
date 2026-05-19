package sqlite

import (
	"database/sql"
	"errors"
	"os"

	"github.com/MarkusFank/rdfmap2go/internal/datareader"
	"github.com/MarkusFank/rdfmap2go/internal/mapping"
	_ "modernc.org/sqlite"
)

type SqliteDataReader struct {
	// sqliteFile    string
	isInitialized bool
	sourceConfig  *mapping.SqliteSourceConfig
	// db            *sql.DB
	// sqlRows       *sql.Rows
	// columns       *[]string
	// columnTypes   []*sql.ColumnType
}

func (reader *SqliteDataReader) Init(sourceConfig mapping.SourceConfig) error {
	sqliteSourceConfig, isSqliteSourceConfig := sourceConfig.(mapping.SqliteSourceConfig)

	if !isSqliteSourceConfig {
		return errors.New("Specified source cofing is not valid")
	}

	_, err := os.Stat(sqliteSourceConfig.File) // check if file exists

	if err != nil {
		return err
	}

	reader.sourceConfig = &sqliteSourceConfig

	// db, err := sql.Open("sqlite", file)

	// if err != nil {
	// 	return err
	// }

	// reader.db = db

	// rows, err := db.Query(query)

	// if err != nil {
	// 	return nil
	// }

	// reader.sqlRows = rows

	// columns, err := reader.sqlRows.Columns()

	// if err != nil {
	// 	return err
	// }

	// columnTypes, err := reader.sqlRows.ColumnTypes()

	// if err != nil {
	// 	return err
	// }

	// reader.columns = &columns
	// reader.columnTypes = columnTypes

	reader.isInitialized = true

	return nil
}

func (reader *SqliteDataReader) Read() (<-chan datareader.RowResult, error) {
	if !reader.isInitialized {
		return nil, errors.New("SqliteDataReader has to be initialized before it can be used!")
	}

	db, err := sql.Open("sqlite", reader.sourceConfig.File)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(reader.sourceConfig.Query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	columns, err := rows.Columns()

	if err != nil {
		return nil, err
	}

	// columnTypes, err := rows.ColumnTypes()

	if err != nil {
		return nil, err
	}

	channel := make(chan datareader.RowResult)

	go func() {
		for {
			hasRow := rows.Next()

			if !hasRow {
				close(channel)
				return
			}

			values := make([]any, len(columns))
			valuePtrs := make([]any, len(columns))

			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				channel <- datareader.RowResult{Error: err}
			}

			dataRow := datareader.DataRow{}
			for i, col := range columns {
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

			channel <- datareader.RowResult{Row: dataRow}
		}
	}()

	return channel, nil
}

func (reader *SqliteDataReader) Close() {
}
