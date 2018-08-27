### sparkey-go - Sparkey in Go

`sparkey-go` is a [cgo](https://golang.org/cmd/cgo/) binding around the `sparkey` library. `sparkey` is a disk-based hash table.

`sparkey` writes two files to disk, for its use: a log file, and an index file. (e.g. "blargh.spl" and "blargh.spi")

### Install

First install the `sparkey` library:

#### OSX

`brew install sparkey`

#### Unix

``

For now, `sparkey-go` uses `go dep`:

`dep ensure -add github.com/tiegz/sparkey-go`

### Example

#### Creating the store

``` go
s := sparkey.New("sparkey_db", sparkey.COMPRESSION_NONE, 1024)
```

#### Setting values

``` go
s.Put("first", "Hello")
s.Put("second", "Worlb")
s.Put("third", "Goodbye")
s.Put("fourth", "EOM")
s.Flush()
```

#### Deleting values

``` go
s.Delete("third")
s.Flush()
```

### Getting values

``` go
fmt.Printf("First value is %s", s.Get("first"))

// Hello
```

#### Iterating over values

``` go
s.ForEachHash(func(k, v string) {
  fmt.Printf("%s: %s\n", k, v)
})

// first: Hello
// second: Worlb
// fourth: EOM
// }
```

#### Inspecting the store

``` go
s.PrettyPrintHash();

// {
//   first => Hello
//   second => Worlb
//   fourth => EOM
// }
```

#### Inspecting the store

``` go
fmt.Printf("Sparkey info:\n\n")
fmt.Printf("  Basename:\t\t%s\n", s.Basename)
fmt.Printf("  CompressionType:\t\t%d\n", s.CompressionType)
fmt.Printf("  Size:\t\t%d\n", s.Size())
fmt.Printf("  LogWriter.CompressionType:\t\t%d\n", s.CompressionType)
fmt.Printf("  LogWriter.BlockSize:\t\t%d\n", s.BlockSize)
fmt.Printf("  MaxKeyLen:\t\t%d\n", s.MaxKeyLen())
fmt.Printf("  MaxValueLen:\t\t%d\n", s.MaxValueLen())
s.Close()


// Sparkey info:
//
//   Basename:   sparkey_db
//   CompressionType:    0
//   Size:   3
//   LogWriter.CompressionType:    0
//   LogWriter.BlockSize:    1024
//   MaxKeyLen:    6
//   MaxValueLen:    7
```

### Running Tests

```
cd sparkey
go test -v .
```

### Running Benches

```
cd sparkey
go test -bench=. -v
```
