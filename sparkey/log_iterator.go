package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"

func NewLogIterator(lr *C.sparkey_logreader) *C.sparkey_logiter {
	var li *C.sparkey_logiter
	C.sparkey_logiter_create(&li, lr)
	return li
}
