package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"

type HashWriter struct {
  Basename string
}

func NewHashWriter(basename string) *HashWriter {
  hw := HashWriter{
    Basename: basename,
  }
  hash_filename := hw.Basename + ".spi"
  log_filename  := hw.Basename + ".spl"
  C.sparkey_hash_write(
    C.CString(hash_filename),
    C.CString(log_filename),
    0)
  return &hw
}
