package json

import (
	"errors"
	"os"
	"testing"

	"github.com/MarkusFank/rdfmap2go/internal/mapping"
)

func TestRead(t *testing.T) {
	/*
		Tests whether the data reader can process json files where the array to read is the root
	*/
	testFile := "./testdata/data.json"

	sourceConfig := mapping.JsonSourceConfig{}
	sourceConfig.File = testFile
	sourceConfig.Type = "json"

	testReadInternal(t, sourceConfig)
}

func TestReadWithJsonPath(t *testing.T) {
	/*
		Tests whether the data reader can process json files with a specified json path
	*/
	testFile := "./testdata/data_with_jsonpath.json"

	sourceConfig := mapping.JsonSourceConfig{}
	sourceConfig.File = testFile
	sourceConfig.Type = "json"
	sourceConfig.JsonPath = "here_is_the_data"

	testReadInternal(t, sourceConfig)
}

func testReadInternal(t *testing.T, sourceConfig mapping.JsonSourceConfig) {
	jsonReader := JsonDataReader{}
	err := jsonReader.Init(sourceConfig)

	if err != nil {
		t.Fatal(err)
	}

	readerChannel, err := jsonReader.Read()

	if err != nil {
		t.Fatal(err)
	}

	compareData := getCompareData()

	cnt := 0
	for row := range readerChannel {
		if row.Error != nil {
			t.Fatal(row.Error)
		}

		compareDataForRow := compareData[cnt]

		for k, v := range compareDataForRow {
			actVal, hasActVal := row.Row[k]

			if !hasActVal {
				t.Fatalf("Expected row %d to have field '%s', but it does not exist", cnt, k)
			}

			if actVal != v {
				t.Fatalf("Expected field '%s' in row %d to be %v, but it is %v instead", k, cnt, actVal, v)
			}
		}

		for k := range row.Row {
			_, hasExpectedVal := compareDataForRow[k]

			if !hasExpectedVal {
				t.Fatalf("Row %d contains a field '%s' which was not expected", cnt, k)
			}
		}

		cnt++
	}

	if cnt != 2 {
		t.Fatalf("Expected reader to return 2 rows, but got %d instead", cnt)
	}
}

func TestFileNotFound(t *testing.T) {
	/*
		Tests whether the data reader returns an error when the specified json file does not exist
	*/
	testFile := "/a/path/that/does/not/exist/nonsense.json"

	sourceConfig := mapping.JsonSourceConfig{}
	sourceConfig.File = testFile
	sourceConfig.Type = "json"

	jsonReader := JsonDataReader{}
	err := jsonReader.Init(sourceConfig)

	if err == nil {
		t.Fatal("Expected an error as the file does not exist")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Expected an error that indicates that the error does not exist, instead got %v", err)
	}
}

func TestInvalidOrEmptyJson(t *testing.T) {
	tests := []struct {
		testName     string
		filePath     string
		errorMessage string
	}{
		{testName: "TestEmptyFile", filePath: "./testdata/empty.json", errorMessage: "Expected an error as the file is empty"},
		{testName: "TestInvalidJson", filePath: "./testdata/invalid_json.json", errorMessage: "Expected an error as the file contains invalid JSON"},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			sourceConfig := mapping.JsonSourceConfig{}
			sourceConfig.File = test.filePath
			sourceConfig.Type = "json"

			jsonReader := JsonDataReader{}
			err := jsonReader.Init(sourceConfig)

			if err != nil {
				t.Fatal(err)
			}

			_, err = jsonReader.Read()

			if err == nil {
				t.Fatal(test.errorMessage)
			}
		})
	}
}

func TestMalformedArray(t *testing.T) {
	/*
		Tests whether the data reader returns an error when an item of the array to read is not an object
		The first row of the testdata is malformed, so this should return an error
	*/

	testFile := "./testdata/data_malformed.json"

	sourceConfig := mapping.JsonSourceConfig{}
	sourceConfig.File = testFile
	sourceConfig.Type = "json"

	jsonReader := JsonDataReader{}
	err := jsonReader.Init(sourceConfig)

	if err != nil {
		t.Fatal(err)
	}

	readerChannel, err := jsonReader.Read()

	if err != nil {
		t.Fatal(err)
	}

	cnt := 0
	for row := range readerChannel {

		if cnt == 0 && row.Error == nil {
			t.Fatal("Expected first row of testdata to return an error, but it didn't")
		}

		if cnt > 0 && row.Error != nil {
			t.Fatalf("Expected row %d to be readable, but got error: %v", cnt, row.Error)
		}

		if cnt > 0 && row.Row == nil {
			t.Fatalf("Expected row %d to be readable, but now data row was returned", cnt)
		}

		cnt++
	}
}

func getCompareData() []map[string]any {
	compareData := []map[string]any{}
	compareData = append(compareData, map[string]any{"manufacturer": "Audi", "modelName": "A4"})
	compareData = append(compareData, map[string]any{"manufacturer": "Volkswagen", "modelName": "Golf"})

	return compareData
}
