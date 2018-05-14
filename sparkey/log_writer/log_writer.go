package log_writer

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"
import "unsafe"

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

func (lw *LogWriter) Put(key string, value string) {
  cKey := (*C.uchar)(unsafe.Pointer(C.CString(key)))
  cValue := (*C.uchar)(unsafe.Pointer(C.CString(value)))

  defer C.free(unsafe.Pointer(cKey))
  defer C.free(unsafe.Pointer(cValue))

  C.sparkey_logwriter_put(
    lw.Instance,
    C.ulonglong(len(key)),
    cKey,
    C.ulonglong(len(value)),
    cValue)
}

func (lw *LogWriter) Delete(key string) {
  cKey := (*C.uchar)(unsafe.Pointer(C.CString(key)))

  C.sparkey_logwriter_delete(
    lw.Instance,
    C.ulonglong((len(key))),
    cKey)
}

func (lw *LogWriter) Flush() {
  C.sparkey_logwriter_flush(lw.Instance)
}
