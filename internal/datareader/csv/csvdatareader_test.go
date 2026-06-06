package csv

import (
	"testing"

	"github.com/MarkusFank/rdfmap2go/internal/mapping"
)

func TestRead(t *testing.T) {
	/*
		Tests whether the data reader can process csv files
	*/
	testFile := "./testdata/data.csv"

	sourceConfig := mapping.CsvSourceConfig{}
	sourceConfig.File = testFile
	sourceConfig.Type = "csv"

	csvReader := CsvDataReader{}
	err := csvReader.Init(sourceConfig)

	if err != nil {
		t.Fatal(err)
	}

	readerChannel, err := csvReader.Read()

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
}

func TestInvalidSourceConfig(t *testing.T) {
	/*
		Tests whether the data reader makes sure it receives a correct source config
	*/

	sourceConfig := mapping.JsonSourceConfig{}
	sourceConfig.File = "asdf"
	sourceConfig.Type = "json"

	csvReader := CsvDataReader{}
	err := csvReader.Init(sourceConfig)

	if err == nil {
		t.Fatal("Expected data reader to return an error as the source config is not a valid csv source config")
	}
}

func TestUninitializedReader(t *testing.T) {
	/*
		Tests whether the data reader makes sure it was initialized before reading
	*/

	sourceConfig := mapping.CsvSourceConfig{}
	sourceConfig.File = "asdf"
	sourceConfig.Type = "csv"

	csvReader := CsvDataReader{}

	_, err := csvReader.Read()

	if err == nil {
		t.Fatal("Expected data reader to return an error as the Init function was not called")
	}
}

func TestFileNotFound(t *testing.T) {
	/*
		Tests whether the data reader handles files that do not exist correctly
	*/

	sourceConfig := mapping.CsvSourceConfig{}
	sourceConfig.File = "/does/not/exist.csv"
	sourceConfig.Type = "csv"

	csvReader := CsvDataReader{}
	err := csvReader.Init(sourceConfig)

	if err == nil {
		t.Fatal("Expected data reader to return an error as the specified csv file does not exist")
	}
}

func TestEmptyFile(t *testing.T) {
	/*
		Tests whether the data reader handles empty csv files correctly
	*/
	sourceConfig := mapping.CsvSourceConfig{}
	sourceConfig.File = "./testdata/empty.csv"
	sourceConfig.Type = "csv"

	csvReader := CsvDataReader{}
	err := csvReader.Init(sourceConfig)

	if err != nil {
		t.Fatal(err)
	}

	_, err = csvReader.Read()

	if err == nil {
		t.Fatal("Expected data reader to return an error as the csv file is empty")
	}
}

func TestIncompleteData(t *testing.T) {
	/*
		Tests whether the data reader can process incomplete csv files
	*/
	testFile := "./testdata/data_incomplete.csv"

	sourceConfig := mapping.CsvSourceConfig{}
	sourceConfig.File = testFile
	sourceConfig.Type = "csv"

	csvReader := CsvDataReader{}
	err := csvReader.Init(sourceConfig)

	if err != nil {
		t.Fatal(err)
	}

	readerChannel, err := csvReader.Read()

	if err != nil {
		t.Fatal(err)
	}

	cnt := 0

	for row := range readerChannel {

		if cnt == 2 && row.Error == nil {
			t.Fatal("Expected data reader to return an error on row 2 as it is incomplete")
		}

		if cnt < 2 && row.Error != nil {
			t.Fatal(row.Error)
		}

		cnt++
	}
}

func getCompareData() []map[string]any {
	compareData := []map[string]any{}
	compareData = append(compareData, map[string]any{"id": "1", "firstName": "Michael", "lastName": "Scott"})
	compareData = append(compareData, map[string]any{"id": "2", "firstName": "Dwight", "lastName": "Schrute"})
	compareData = append(compareData, map[string]any{"id": "3", "firstName": "Jim", "lastName": "Halpert"})

	return compareData
}
