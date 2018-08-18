package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"
import "fmt"

type LogReader struct {
  Basename string
  Instance *C.sparkey_logreader
}

func OpenLogReader(basename string) *C.sparkey_logreader {
  var lr *C.sparkey_logreader
  filename    := basename + ".spl"
  return_code := C.sparkey_logreader_open(
    &lr,
    C.CString(filename))

  fmt.Printf("OpenLogReader:\t\t%s, Return Code: %d\n", filename, return_code)

  return lr
}
