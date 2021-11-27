package namedcsvreader

import (
	"encoding/csv"
	"io"
	"math/bits"
	"os"
	"strconv"
	"time"
)

// Named csv reader represents the reader instance.
type NamedCsvReader struct {
	headers     []string       // Headers available in the csv file
	headerMap   map[string]int // Headers mapped to their index
	index       int            // Index of the row the reader is currently add
	innerReader *csv.Reader    // The reader containing the csv data
}

// Record represents a single record of the named csv reader.
type Record struct {
	Error   error    // Error that occurend when reading the record, if any.
	RowNum  int      // The row number of the record
	current []string // The record that is currently beaing read.
	*NamedCsvReader
}

// NewReader returns a new named csv reader for the given file.
func NewReader(file *os.File) *NamedCsvReader {
	reader := csv.NewReader(file)
	return &NamedCsvReader{nil, make(map[string]int), 0, reader}
}

// WithHeaders configures the headers used to read the csv data.
func (reader *NamedCsvReader) WithHeaders(headers ...string) *NamedCsvReader {
	reader.headers = headers
	reader.headerMap = make(map[string]int)
	return reader
}

// Read advances the reader to the next record.
func (reader *NamedCsvReader) Read() chan *Record {
	if reader.headers == nil {
		record, err := reader.innerReader.Read()
		if err != nil {
			panic(err)
		}
		reader.headers = record
	}
	return reader.execReader()
}

// GetInt returns the value of the specified column as int.
// When no value is found it returns 0.
func (record *Record) GetInt(key string) int {
	value, err := strconv.ParseInt(record.GetString(key), 10, bits.UintSize)
	if err != nil {
		return 0
	}
	return int(value)
}

// GetInt32 returns the value of the specified column as int32.
// When no value is found it returns 0.
func (record *Record) GetInt32(key string) int32 {
	value, err := strconv.ParseInt(record.GetString(key), 10, 32)
	if err != nil {
		return 0
	}
	return int32(value)
}

// GetInt64 returns the value of the specified column as int64.
// When no value is found it returns 0.
func (record *Record) GetInt64(key string) int64 {
	value, err := strconv.ParseInt(record.GetString(key), 10, 64)
	if err != nil {
		return 0
	}
	return value
}

// GetFloat32 returns the value of the specified column as float32.
// When no value is found it returns 0.
func (record *Record) GetFloat32(key string) float32 {
	value, err := strconv.ParseFloat(record.GetString(key), 32)
	if err != nil {
		return 0
	}
	return float32(value)
}

// GetFloat64 returns the value of the specified column as float64.
// When no value is found it returns 0.
func (record *Record) GetFloat64(key string) float64 {
	value, err := strconv.ParseFloat(record.GetString(key), 64)
	if err != nil {
		return 0
	}
	return value
}

// GetBoolean returns the value of the specified column as boolean.
// 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False as value.
// Any other value returns false.
func (record *Record) GetBoolean(key string) bool {
	value, err := strconv.ParseBool(record.GetString(key))
	if err != nil {
		return false
	}
	return value
}

// GetTime returns the value of the specified column as date time.
// When no value is found it return the default date time.
// See the time documentation for the available layouts.
func (record *Record) GetTime(layout string, key string) time.Time {
	value, err := time.Parse(layout, record.GetString(key))
	if err != nil {
		return time.Time{}
	}
	return value
}

// GetString returns the value of the specified column as string.
// When no value is found it returns an empty string.
func (record *Record) GetString(key string) string {
	if index, ok := record.headerMap[key]; ok {
		return record.current[index]
	}

	for index, value := range record.headers {
		if value == key {
			record.headerMap[key] = index
			return record.current[index]
		}
	}

	return ""
}

func (reader *NamedCsvReader) execReader() chan *Record {
	ch := make(chan *Record)

	go func() {
		for {
			record, err := reader.innerReader.Read()
			if err != nil {
				if err != io.EOF {
					ch <- &Record{err, 0, nil, reader}
				}
				break
			}
			ch <- &Record{nil, reader.index, record, reader}
			reader.index = reader.index + 1
		}
		close(ch)
	}()

	return ch
}
