package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"
import "fmt"

func NewHashReader(basename string) *C.sparkey_hashreader {
  var hr *C.sparkey_hashreader
  hash_filename := basename + ".spi"
  log_filename  := basename + ".spl"
  return_code   := C.sparkey_hash_open(
    &hr,
    C.CString(hash_filename),
    C.CString(log_filename))

  fmt.Printf("NewHashReader %s, %s, Return Code: %d\n", hash_filename, log_filename, return_code)

  return hr
}
