package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"
import "fmt"

func OpenLogReader(log_filename string) *C.sparkey_logreader {
  var lr *C.sparkey_logreader
  return_code := C.sparkey_logreader_open(
    &lr,
    C.CString(log_filename))

  fmt.Printf("OpenLogReader:\t\t%s, Return Code: %d\n", log_filename, return_code)

  return lr
}
