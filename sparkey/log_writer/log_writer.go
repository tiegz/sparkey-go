package log_writer

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"

// TODO use iota instead?
const COMPRESSION_NONE int  = 0
const COMPRESSION_SNAPPY int = 1

type LogWriter struct {
  Filename string
  CompressionType int
  BlockSize int
  Instance *C.sparkey_logwriter
}

func New(filename string) *LogWriter {
  lw := LogWriter{
    Filename: filename,
    CompressionType: COMPRESSION_NONE,
    BlockSize: 1024, // ??
    Instance: nil,
  }
  lw.Create()
  return &lw
}

func (lw *LogWriter) Create() {
  C.sparkey_logwriter_create(
    &lw.Instance,
    C.CString(lw.Filename),
    C.sparkey_compression_type(lw.CompressionType),
    C.int(lw.BlockSize))
}

