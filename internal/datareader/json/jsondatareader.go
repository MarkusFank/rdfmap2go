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
	sourceConfig  *mapping.JsonSourceConfig
	isInitialized bool
}

func (reader *JsonDataReader) Init(sourceConfig mapping.SourceConfig) error {
	jsonSourceConfig, isJsonSourceConfig := sourceConfig.(mapping.JsonSourceConfig)

	if !isJsonSourceConfig {
		return errors.New("Specified source config is not valid")
	}

	_, err := os.Stat(jsonSourceConfig.File) // check if file exists

	if err != nil {
		return err
	}

	reader.sourceConfig = &jsonSourceConfig

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
		currentRow := 0
		res.ForEach(func(key, val gjson.Result) bool {
			row := val.Value()

			// TODO ATM, we assume that row only contains "primitives" like string, number, etc. We handle complex objects later

			objMap, isMap := row.(map[string]any)

			if !isMap {
				channel <- datareader.RowResult{Error: fmt.Errorf("Unable to handle row %d", currentRow)}
			}

			dataRow := datareader.DataRow{}
			maps.Copy(dataRow, objMap)

			channel <- datareader.RowResult{Row: dataRow}

			currentRow++

			return true
		})

		close(channel)
	}()

	return channel, nil
}
