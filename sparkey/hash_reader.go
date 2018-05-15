package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"

func NewHashReader(basename string) *C.sparkey_hashreader {
  var hr *C.sparkey_hashreader
  hash_filename := basename + ".spi"
  log_filename  := basename + ".spl"
  C.sparkey_hash_open(&hr, C.CString(hash_filename), C.CString(log_filename))
  return hr
}

