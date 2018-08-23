package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"

func OpenLogReader(log_filename string) *C.sparkey_logreader {
  var lr *C.sparkey_logreader
  C.sparkey_logreader_open(
    &lr,
    C.CString(log_filename))

  return lr
}
