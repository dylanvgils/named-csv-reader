# Named csv reader
Package named csv reader provides methods to easly read csv files and parse columsn to basic types.

This package is not a CSV parser, it uses the [encodig/csv reader](https://pkg.go.dev/encoding/csv#NewReader) to prase the csv file.

## Installation
```
go get github.com/dylanvgils/namedcsvreader
```

## Examples

### CSV file with headers
Read csv file using the headers in the csv file.

```go
file, err := os.Open("testdata/valid.csv")
if err != nil {
    panic(err)
}

for record := range namedcsvreader.NewReader(file).Read() {
    fmt.Printf("Name: %s; Age: %d\n", record.GetString("name"), record.GetInt("age"))
}

// Output:
// Name: Colorado Leon; Age: 26
// Name: Rajah Fletcher; Age: 47
// Name: Tobias Snow; Age: 16
// Name: Jared Finch; Age: 19
```

### CSV file without headers
Read a csv file that does not contain headers, by specifying your own.

```go
file, err := os.Open("testdata/valid_without_headers.csv")
if err != nil {
    panic(err)
}

reader := namedcsvreader.NewReader(file).
    WithHeaders("name", "age")

for record := range reader.Read() {
    fmt.Printf("Name: %s; Age: %d\n", record.GetString("name"), record.GetInt("age"))
}

// Output:
// Name: Colorado Leon; Age: 26
// Name: Rajah Fletcher; Age: 47
```
