package json

import (
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/MarkusFank/rdfmap2go/internal/datareader"
	"github.com/tidwall/gjson"
)

type JsonDataReader struct {
	currentRow int
	rows       []any
}

func (reader *JsonDataReader) Init(file, jsonPath string) error {
	bytes, err := os.ReadFile(file) // TODO do not read entire file content at once

	if err != nil {
		return err
	}

	var res gjson.Result

	if len(strings.TrimSpace(jsonPath)) == 0 {
		res = gjson.ParseBytes(bytes) // TODO support json paths from config
	} else {
		fmt.Printf("JSON path: %s\n", jsonPath)
		res = gjson.GetBytes(bytes, jsonPath)
	}

	array, isArray := res.Value().([]any)

	if !isArray {
		return fmt.Errorf("Unable to process JSON data for file '%s'", file)
	}

	reader.rows = array

	// var targetObj any
	// err = json.Unmarshal(bytes, &targetObj)

	// if err != nil {
	// 	return err
	// }

	// if targetObj == nil {
	// 	return errors.New("Unable to unmarshall json data")
	// }

	// switch typed := targetObj.(type) {
	// case []any:
	// 	reader.rows = typed
	// case map[string]any:
	// 	// TODO handle complex structure
	// }

	reader.currentRow = -1

	return nil
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
