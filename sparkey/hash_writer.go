package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"
import "fmt"

type HashWriter struct {}

func NewHashWriter(basename string) *HashWriter {
  hw := HashWriter{}
  hash_filename := basename + ".spi"
  log_filename  := basename + ".spl"
  return_code   := C.sparkey_hash_write(
    C.CString(hash_filename),
    C.CString(log_filename),
    0)
  fmt.Printf("NewHashWriter: %s, %s, Return Code: %d\n", hash_filename, log_filename, return_code)

  return &hw
}
