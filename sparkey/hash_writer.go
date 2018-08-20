package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"
import "fmt"

type HashWriter struct {}

func NewHashWriter(log_filename string, index_filename string) *HashWriter {
  hw := HashWriter{}
  return_code   := C.sparkey_hash_write(
    C.CString(index_filename),
    C.CString(log_filename),
    0)
  fmt.Printf("NewHashWriter:\t\t%s, %s, Return Code: %d\n", log_filename, index_filename, return_code)

  return &hw
}

