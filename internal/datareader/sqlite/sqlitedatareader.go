package sqlite

import (
	"database/sql"
	"errors"
	"os"
	"strings"

	"github.com/MarkusFank/rdfmap2go/internal/datareader"
	"github.com/MarkusFank/rdfmap2go/internal/download"
	"github.com/MarkusFank/rdfmap2go/internal/mapping"
	_ "modernc.org/sqlite"
)

type SqliteDataReader struct {
	isInitialized bool
	sourceConfig  *mapping.SqliteSourceConfig
	isRemoteFile  bool
}

func (reader *SqliteDataReader) Init(sourceConfig mapping.SourceConfig) error {
	sqliteSourceConfig, isSqliteSourceConfig := sourceConfig.(mapping.SqliteSourceConfig)

	if !isSqliteSourceConfig {
		return errors.New("Specified source config is not valid")
	}

	if strings.HasPrefix(sqliteSourceConfig.File, "https://") || strings.HasPrefix(sqliteSourceConfig.File, "http://") {
		reader.isRemoteFile = true
	} else {
		_, err := os.Stat(sqliteSourceConfig.File) // check if file exists
		if err != nil {
			return err
		}
	}

	reader.sourceConfig = &sqliteSourceConfig

	reader.isInitialized = true

	return nil
}

func (reader *SqliteDataReader) Read() (<-chan datareader.RowResult, error) {
	if !reader.isInitialized {
		return nil, errors.New("SqliteDataReader has to be initialized before it can be used!")
	}

	var fileName string
	var deleteTempFileInOuterFunction bool

	if !reader.isRemoteFile {
		fileName = reader.sourceConfig.File
	} else {
		tempFile, err := os.CreateTemp(os.TempDir(), "rdf2go-temp-*.sqlite")

		if err != nil {
			return nil, err
		}

		deleteTempFileInOuterFunction = true

		defer func() {
			if deleteTempFileInOuterFunction {
				os.Remove(tempFile.Name())
			}
		}()

		err = download.DownloadFile(reader.sourceConfig.File, tempFile)
		tempFile.Close()

		if err != nil {
			return nil, err
		}

		fileName = tempFile.Name()
	}

	db, err := sql.Open("sqlite", fileName)

	if err != nil {
		return nil, err
	}

	rows, err := db.Query(reader.sourceConfig.Query)

	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()

	if err != nil {
		return nil, err
	}

	// TODO db.Close() and rows.Close() are never called if we return early due to an error

	deleteTempFileInOuterFunction = false

	// columnTypes, err := rows.ColumnTypes()

	// if err != nil {
	// 	return nil, err
	// }

	channel := make(chan datareader.RowResult)

	go func() {
		defer close(channel)
		defer db.Close()
		defer rows.Close()
		defer func(isTempFile bool, file string) {
			if isTempFile {
				os.Remove(file)
			}
		}(reader.isRemoteFile, fileName)

		for {
			hasRow := rows.Next()

			if !hasRow {
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
