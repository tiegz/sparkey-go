package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"
// import "unsafe"
import "fmt"

func NewLogWriter(basename string, compression_type int, block_size int) *C.sparkey_logwriter {
  var lw *C.sparkey_logwriter
  filename    := basename + ".spl"
  return_code := C.sparkey_logwriter_create(
    &lw,
    C.CString(filename),
    C.sparkey_compression_type(compression_type),
    C.int(block_size))

  fmt.Printf("NewLogWriter:\t\t%s, ReturnCode: %d\n", filename, return_code)

  return lw
}


func OpenLogWriter(basename string) *C.sparkey_logwriter {
  var lw *C.sparkey_logwriter
  filename    := basename + ".spl"
  return_code := C.sparkey_logwriter_append(
    &lw,
    C.CString(filename))

  fmt.Printf("OpenLogWriter:\t\t%s, ReturnCode: %d\n", filename, return_code)

  return lw
}

