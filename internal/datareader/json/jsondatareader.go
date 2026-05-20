package json

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/MarkusFank/rdfmap2go/internal/datareader"
	"github.com/MarkusFank/rdfmap2go/internal/mapping"
	"github.com/tidwall/gjson"
)

type JsonDataReader struct {
	// currentRow int
	// rows       []any
	sourceConfig  *mapping.JsonSourceConfig
	isInitialized bool
}

func (reader *JsonDataReader) Init(sourceConfig mapping.SourceConfig) error {
	jsonSourceConfig, isJsonSourceConfig := sourceConfig.(mapping.JsonSourceConfig)

	if !isJsonSourceConfig {
		return errors.New("Specified source cofing is not valid")
	}

	_, err := os.Stat(jsonSourceConfig.File) // check if file exists

	if err != nil {
		return err
	}

	reader.sourceConfig = &jsonSourceConfig

	// bytes, err := os.ReadFile(file) // TODO do not read entire file content at once

	// if err != nil {
	// 	return err
	// }

	// var res gjson.Result

	// if len(strings.TrimSpace(jsonPath)) == 0 {
	// 	res = gjson.ParseBytes(bytes) // TODO support json paths from config
	// } else {
	// 	fmt.Printf("JSON path: %s\n", jsonPath)
	// 	res = gjson.GetBytes(bytes, jsonPath)
	// }

	array, isArray := res.Value().([]any)

	if !isArray {
		return fmt.Errorf("Unable to process JSON data for file '%s'", file)
	}

	// reader.rows = array

	// reader.currentRow = -1

	reader.isInitialized = true

	return nil
}

func (reader *JsonDataReader) Read() (<-chan datareader.RowResult, error) {

	if !reader.isInitialized {
		return nil, errors.New("JsonDataReader has to be initialized before it can be used!")
	}

	bytes, err := os.ReadFile(reader.sourceConfig.File) // TODO do not read entire file content at once

	if err != nil {
		return nil, err
	}

	var res gjson.Result

	jsonPath := reader.sourceConfig.JsonPath
	if len(strings.TrimSpace(jsonPath)) == 0 {
		res = gjson.ParseBytes(bytes)
	} else {
		res = gjson.GetBytes(bytes, jsonPath)
	}

	channel := make(chan datareader.RowResult)

	go func() {
		res.ForEach(func(key, val gjson.Result) bool {

			return true
		})
	}()

	return channel, nil
}

func (reader *JsonDataReader) ReadRow() (*datareader.DataRow, error) {

	reader.currentRow++

	// if reader.currentRow < len(reader.result) {
	// 	row := reader.result[reader.currentRow]

	// 	dataRow := datareader.DataRow{}
	// 	row.ForEach(func(key, val gjson.Result) bool {
	// 		keyStr := key.String() // we assume the key is a string (property name)

	// 		// TODO ATM, we assume that row only contains "primitives" like string, number, etc. We handle complex objects later
	// 		dataRow[keyStr] = val.Value()

	// 		return true
	// 	})

	// 	return &dataRow, nil
	// }

	if reader.currentRow < len(reader.rows) {
		row := reader.rows[reader.currentRow]

		// TODO ATM, we assume that row only contains "primitives" like string, number, etc. We handle complex objects later

		objMap, isMap := row.(map[string]any)

		if !isMap {
			return nil, fmt.Errorf("Unable to handle row %d", reader.currentRow)
		}

		dataRow := datareader.DataRow{}
		maps.Copy(dataRow, objMap)

		return &dataRow, nil
	}

	return nil, nil
}

func (reader *JsonDataReader) Close() {}
