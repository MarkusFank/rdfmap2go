package csv

import (
	"encoding/csv"
	"errors"
	"io"
	"os"

	"github.com/MarkusFank/rdfmap2go/internal/datareader"
)

type CsvDataReader struct {
	csvFile       string
	isInitialized bool
	csvReader     *csv.Reader
	file          *os.File
	headers       []string
}

func (r *CsvDataReader) Init(file string) error {
	r.csvFile = file

	f, err := os.Open(file)
	r.file = f

	if err != nil {
		return err
	}

	// defer f.Close()

	r.csvReader = csv.NewReader(f)

	record, err := r.csvReader.Read()
	r.headers = record

	if err == io.EOF {
		return errors.New("CSV file contains no data!")
	} else if err != nil {
		return err
	}

	r.isInitialized = true
	return nil
}

func (r *CsvDataReader) Close() {
	if r.isInitialized && r.csvReader != nil {
		r.file.Close()
	}
}

func (r *CsvDataReader) ReadRow() (*datareader.DataRow, error) {
	if !r.isInitialized {
		return nil, errors.New("CsvDataReader has to be initialized before it can be used!")
	}

	record, err := r.csvReader.Read()

	if err == io.EOF {
		return nil, nil
	}

	row := datareader.DataRow{}
	for i, header := range r.headers {
		row[header] = record[i]
	}

	return &row, nil
}
