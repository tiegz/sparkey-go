package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"
import "fmt"

func OpenHashReader(log_filename string, index_filename string) *C.sparkey_hashreader {
  var hr *C.sparkey_hashreader
  return_code   := C.sparkey_hash_open(
    &hr,
    C.CString(index_filename),
    C.CString(log_filename))

  fmt.Printf("OpenHashReader\t\t%s, %s, Return Code: %d\n", index_filename, log_filename, return_code)

  return hr
}

