package sqlite

import (
	"testing"

	"github.com/MarkusFank/rdfmap2go/internal/mapping"
)

func TestRead(t *testing.T) {
	/*
		Tests whether the data reader can process sqlite files
	*/

	testFile := "./testdata/depts.sqlite"

	sourceConfig := mapping.SqliteSourceConfig{}
	sourceConfig.File = testFile
	sourceConfig.Type = "sqlite"
	sourceConfig.Query = "select * from departments"

	sqliteReader := SqliteDataReader{}
	err := sqliteReader.Init(sourceConfig)

	if err != nil {
		t.Fatal(err)
	}

	readerChannel, err := sqliteReader.Read()

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

	if cnt != 3 {
		t.Fatalf("Expected reader to return 3 rows, but got %d rows instead", cnt)
	}
}

func TestInvalidSourceConfig(t *testing.T) {
	/*
		Tests whether the data reader makes sure it receives a correct source config
	*/

	testFile := "./testdata/depts.sqlite"

	sourceConfig := mapping.CsvSourceConfig{}
	sourceConfig.File = testFile
	sourceConfig.Type = "csv"

	sqliteReader := SqliteDataReader{}
	err := sqliteReader.Init(sourceConfig)

	if err == nil {
		t.Fatal("Expected data reader to return an error as the source config is not a valid sqlite source config")
	}
}

func TestUninitializedReader(t *testing.T) {
	/*
		Tests whether the data reader makes sure it was initialized before reading
	*/

	sourceConfig := mapping.CsvSourceConfig{}
	sourceConfig.File = "asdf"
	sourceConfig.Type = "csv"

	sqliteReader := SqliteDataReader{}

	_, err := sqliteReader.Read()

	if err == nil {
		t.Fatal("Expected data reader to return an error as the Init function was not called")
	}
}

func TestFileNotFound(t *testing.T) {
	/*
		Tests whether the data reader handles files that do not exist correctly
	*/

	sourceConfig := mapping.CsvSourceConfig{}
	sourceConfig.File = "/does/not/exist.sqlite"
	sourceConfig.Type = "csv"

	sqliteReader := SqliteDataReader{}
	err := sqliteReader.Init(sourceConfig)

	if err == nil {
		t.Fatal("Expected data reader to return an error as the specified sqlite file does not exist")
	}
}

func getCompareData() []map[string]any {
	compareData := []map[string]any{}
	compareData = append(compareData, map[string]any{"depId": int64(1), "departmentName": "Sales", "code": "SAL"})
	compareData = append(compareData, map[string]any{"depId": int64(2), "departmentName": "Accounting", "code": "ACC"})
	compareData = append(compareData, map[string]any{"depId": int64(3), "departmentName": "Marketing", "code": "MAR"})

	return compareData
}
