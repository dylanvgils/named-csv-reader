package namedcsvreader

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const dateFormat = "2006-01-02 15:04:05"

type Fixture struct {
	data    string
	headers []string
	want    []map[string]interface{}
}

func TestReaderWithHeadersInCsv(t *testing.T) {
	executeTest(t, Fixture{
		data: "testdata/valid.csv",
		want: []map[string]interface{}{
			{
				"name":      string("Colorado Leon"),
				"age":       int(26),
				"active":    bool(false),
				"weight":    float32(08.09),
				"reference": parseDate(t, "1996-09-22 23:40:01"),
			}, {
				"name":      string("Rajah Fletcher"),
				"age":       int32(47),
				"active":    bool(false),
				"weight":    float64(25.81),
				"reference": parseDate(t, "1994-08-22 18:59:33"),
			}, {
				"name":      string("Tobias Snow"),
				"age":       int64(16),
				"active":    bool(false),
				"weight":    float32(34.64),
				"reference": parseDate(t, "2005-10-29 17:14:49"),
			}, {
				"name":      string("Jared Finch"),
				"age":       int64(19),
				"active":    bool(true),
				"weight":    float64(64.33),
				"reference": parseDate(t, "1980-01-20 11:39:40"),
			},
		},
	})
}

func TestReaderWithoutHeadersInCsv(t *testing.T) {
	executeTest(t, Fixture{
		data:    "testdata/valid_without_headers.csv",
		headers: []string{"name", "age"},
		want: []map[string]interface{}{
			{

				"name": "Colorado Leon",
				"age":  int(26),
			},
			{
				"name": "Rajah Fletcher",
				"age":  int64(47),
			},
		},
	})
}

func executeTest(t *testing.T, fixture Fixture) {
	// Arrange
	file, err := os.Open(fixture.data)
	if err != nil {
		panic(err)
	}

	reader := NewReader(file)

	if fixture.headers != nil {
		reader.WithHeaders(fixture.headers...)
	}

	// Act & Assert
	for record := range reader.Read() {
		expected := fixture.want[record.RowNum]

		for key := range expected {
			expectedValue := expected[key]
			switch expected[key].(type) {
			case string:
				assert.Equal(t, expectedValue, record.GetString(key))
			case int:
				assert.Equal(t, expectedValue, record.GetInt(key))
			case int32:
				assert.Equal(t, expectedValue, record.GetInt32(key))
			case int64:
				assert.Equal(t, expectedValue, record.GetInt64(key))
			case float32:
				assert.Equal(t, expectedValue, record.GetFloat32(key))
			case float64:
				assert.Equal(t, expectedValue, record.GetFloat64(key))
			case bool:
				assert.Equal(t, expectedValue, record.GetBoolean(key))
			case time.Time:
				assert.Equal(t, expectedValue, record.GetTime(dateFormat, key))
			}
		}
	}
}

func parseDate(t *testing.T, value string) time.Time {
	v, err := time.Parse(dateFormat, value)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	return v
}
