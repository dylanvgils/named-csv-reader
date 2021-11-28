# Named csv reader
Package named csv reader provides methods to easly read csv files and parse columsn to basic types.

This package is not a CSV parser, it used the [encodig/csv reader](https://pkg.go.dev/encoding/csv#NewReader) to prase the csv file.

## Installation
```
go get github.com/dylanvgils/namedcsvreader
```

## Examples
Read csv file using the headers in the csv file.

```go
file, err := os.Open("testdata/valid.csv")
if err != nil {
    panic(err)
}

for record := range namedcsvreader.NewReader(file).Read() {
    fmt.Printf("Name: %s; Age: %dd\n", record.GetString("name"), record.GetInt("age")))
}

// Output:
// Name: Colorado Leon; Age: 26
// Name: Rajah Fletcher; Age: 47
// Name: Tobias Snow; Age: 16
// Name: Jared Finch; Age: 19
```
