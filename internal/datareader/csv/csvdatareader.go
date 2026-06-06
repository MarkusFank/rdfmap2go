package csv

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/MarkusFank/rdfmap2go/internal/datareader"
	"github.com/MarkusFank/rdfmap2go/internal/download"
	"github.com/MarkusFank/rdfmap2go/internal/mapping"
)

type CsvDataReader struct {
	isInitialized bool
	sourceConfig  *mapping.CsvSourceConfig
	isRemoteFile  bool
}

func (r *CsvDataReader) Init(sourceConfig mapping.SourceConfig) error {

	csvSourceConfig, isCsvSourceConfig := sourceConfig.(mapping.CsvSourceConfig)

	if !isCsvSourceConfig {
		return errors.New("Specified source config is not valid")
	}

	if strings.HasPrefix(csvSourceConfig.File, "https://") || strings.HasPrefix(csvSourceConfig.File, "http://") {
		r.isRemoteFile = true
	} else {
		_, err := os.Stat(csvSourceConfig.File) // check if file exists

		if err != nil {
			return err
		}
	}
	r.sourceConfig = &csvSourceConfig

	r.isInitialized = true
	return nil
}

func (r *CsvDataReader) Read() (<-chan datareader.RowResult, error) {
	if !r.isInitialized {
		return nil, errors.New("CsvDataReader has to be initialized before it can be used!")
	}

	var fileName string
	var deleteTempFileInOuterFunction bool

	if !r.isRemoteFile {
		fileName = r.sourceConfig.File
	} else {
		tempFile, err := os.CreateTemp(os.TempDir(), "rdf2go-temp-*.csv")

		if err != nil {
			return nil, err
		}

		deleteTempFileInOuterFunction = true

		defer func() {
			if deleteTempFileInOuterFunction {
				os.Remove(tempFile.Name())
			}
		}()

		err = download.DownloadFile(r.sourceConfig.File, tempFile)
		tempFile.Close()

		if err != nil {
			return nil, err
		}

		fileName = tempFile.Name()
	}

	f, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(f)

	record, err := csvReader.Read()
	headers := record

	if err == io.EOF {
		return nil, errors.New("CSV file contains no data!")
	} else if err != nil {
		return nil, err
	}

	// TODO f.Close() is never called if we return early due to an error

	deleteTempFileInOuterFunction = false

	channel := make(chan datareader.RowResult)

	go func() {
		defer close(channel)
		defer f.Close()

		defer func(isTempFile bool, file string) {
			if isTempFile {
				os.Remove(file)
			}
		}(r.isRemoteFile, fileName)

		for {
			record, err := csvReader.Read()

			if err == io.EOF {
				return
			}

			if err != nil {
				channel <- datareader.RowResult{Error: err}
				return
			}

			row := datareader.DataRow{}
			for i, header := range headers {
				if i >= len(record) {
					break
				}

				row[header] = record[i]
			}

			channel <- datareader.RowResult{Row: row}
		}
	}()

	return channel, nil
}
