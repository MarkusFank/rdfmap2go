package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"

	"github.com/MarkusFank/rdfmap2go/internal/datareader"
)

type JsonDataReader struct {
	rows       []any
	currentRow int
}

func (reader *JsonDataReader) Init(file string) error {
	bytes, err := os.ReadFile(file) // TODO do not read entire file content at once

	if err != nil {
		return err
	}

	var targetObj any
	err = json.Unmarshal(bytes, &targetObj)

	if err != nil {
		return err
	}

	if targetObj == nil {
		return errors.New("Unable to unmarshall json data")
	}

	switch typed := targetObj.(type) {
	case []any:
		reader.rows = typed
	case map[string]any:
		// TODO handle complex structure
	}

	reader.currentRow = -1

	return nil
}

func (reader *JsonDataReader) ReadRow() (*datareader.DataRow, error) {
	reader.currentRow++

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
