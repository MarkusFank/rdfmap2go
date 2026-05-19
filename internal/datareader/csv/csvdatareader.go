package csv

import (
	"encoding/csv"
	"errors"
	"io"
	"os"

	"github.com/MarkusFank/rdfmap2go/internal/datareader"
	"github.com/MarkusFank/rdfmap2go/internal/mapping"
)

type CsvDataReader struct {
	isInitialized bool
	sourceConfig  *mapping.CsvSourceConfig
}

func (r *CsvDataReader) Init(sourceConfig mapping.SourceConfig) error {

	csvSourceConfig, isCsvSourceConfig := sourceConfig.(mapping.CsvSourceConfig)

	if !isCsvSourceConfig {
		return errors.New("Specified source cofing is not valid")
	}

	_, err := os.Stat(csvSourceConfig.File) // check if file exists

	if err != nil {
		return err
	}

	r.sourceConfig = &csvSourceConfig

	r.isInitialized = true
	return nil
}

func (r *CsvDataReader) Read() (<-chan datareader.RowResult, error) {
	if !r.isInitialized {
		return nil, errors.New("CsvDataReader has to be initialized before it can be used!")
	}

	f, err := os.Open(r.sourceConfig.File)

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

	channel := make(chan datareader.RowResult)

	go func() {
		defer f.Close()

		for {
			record, err := csvReader.Read()

			if err == io.EOF {
				close(channel)
				return
			}

			if err != nil {
				close(channel)
				channel <- datareader.RowResult{Error: err}
			}

			row := datareader.DataRow{}
			for i, header := range headers {
				row[header] = record[i]
			}

			channel <- datareader.RowResult{Row: row}
		}
	}()

	return channel, nil
}

func (r *CsvDataReader) Close() {

}
