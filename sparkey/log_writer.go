package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"
import "unsafe"

// TODO use iota instead?
const COMPRESSION_NONE int   = 0
const COMPRESSION_SNAPPY int = 1

func NewLogWriter(basename string, compression_type int, block_size int) *C.sparkey_logwriter {
  filename := basename + ".spl"
  var lw *C.sparkey_logwriter
  C.sparkey_logwriter_create(
    &lw,
    C.CString(filename),
    C.sparkey_compression_type(compression_type),
    C.int(block_size))

  return lw
}

func (lw *C.sparkey_logwriter) Put(key string, value string) {
  cKey := (*C.uchar)(unsafe.Pointer(C.CString(key)))
  cValue := (*C.uchar)(unsafe.Pointer(C.CString(value)))

  defer C.free(unsafe.Pointer(cKey))
  defer C.free(unsafe.Pointer(cValue))

  C.sparkey_logwriter_put(
    lw,
    C.ulonglong(len(key)),
    cKey,
    C.ulonglong(len(value)),
    cValue)
}

func (lw *C.sparkey_logwriter) Delete(key string) {
  cKey := (*C.uchar)(unsafe.Pointer(C.CString(key)))

  C.sparkey_logwriter_delete(
    lw,
    C.ulonglong((len(key))),
    cKey)
}

func (lw *C.sparkey_logwriter) Flush() {
  C.sparkey_logwriter_flush(lw)
}
