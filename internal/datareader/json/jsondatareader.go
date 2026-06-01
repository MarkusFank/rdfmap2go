package json

import (
	"bytes"
	"errors"
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/MarkusFank/rdfmap2go/internal/datareader"
	"github.com/MarkusFank/rdfmap2go/internal/download"
	"github.com/MarkusFank/rdfmap2go/internal/mapping"
	"github.com/tidwall/gjson"
)

type JsonDataReader struct {
	sourceConfig  *mapping.JsonSourceConfig
	isRemoteFile  bool
	isInitialized bool
}

func (reader *JsonDataReader) Init(sourceConfig mapping.SourceConfig) error {
	jsonSourceConfig, isJsonSourceConfig := sourceConfig.(mapping.JsonSourceConfig)

	if !isJsonSourceConfig {
		return errors.New("Specified source config is not valid")
	}

	if strings.HasPrefix(jsonSourceConfig.File, "https://") || strings.HasPrefix(jsonSourceConfig.File, "http://") {
		reader.isRemoteFile = true
	} else {
		_, err := os.Stat(jsonSourceConfig.File) // check if file exists

		if err != nil {
			return err
		}
	}

	reader.sourceConfig = &jsonSourceConfig

	reader.isInitialized = true

	return nil
}

func (reader *JsonDataReader) Read() (<-chan datareader.RowResult, error) {

	if !reader.isInitialized {
		return nil, errors.New("JsonDataReader has to be initialized before it can be used!")
	}

	var bytesArr []byte
	var err error

	if !reader.isRemoteFile {
		bytesArr, err = os.ReadFile(reader.sourceConfig.File) // TODO do not read entire file content at once
	} else {
		tempFile, err := os.CreateTemp(os.TempDir(), "rdf2go-temp-*.json")

		if err != nil {
			return nil, err
		}

		defer tempFile.Close()
		defer os.Remove(tempFile.Name())

		writer := bytes.Buffer{}
		err = download.DownloadFile(reader.sourceConfig.File, &writer)

		if err != nil {
			return nil, err
		}

		bytesArr = writer.Bytes()
	}

	if err != nil {
		return nil, err
	}

	var res gjson.Result

	jsonPath := reader.sourceConfig.JsonPath
	if len(strings.TrimSpace(jsonPath)) == 0 {
		res = gjson.ParseBytes(bytesArr)
	} else {
		res = gjson.GetBytes(bytesArr, jsonPath)
	}

	if res.Type == gjson.Null {
		return nil, errors.New("JSON file is invalid or empty!")
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
