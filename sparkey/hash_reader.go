package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"

func OpenHashReader(log_filename string, index_filename string) *C.sparkey_hashreader {
	var hr *C.sparkey_hashreader
	C.sparkey_hash_open(
		&hr,
		C.CString(index_filename),
		C.CString(log_filename))

	return hr
}
