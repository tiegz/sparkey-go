package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"

type LogReader struct {
  Basename string
  Instance *C.sparkey_logreader
}

func NewLogReader(basename string) *C.sparkey_logreader {
  var lr *C.sparkey_logreader
  filename := basename + ".spl"
  C.sparkey_logreader_open(
    &lr,
    C.CString(filename))
  return lr
}

