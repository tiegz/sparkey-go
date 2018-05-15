package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"

type LogIterator struct {
  LogReaderInstance *C.sparkey_logreader
  Instance *C.sparkey_logiter
}

// TODO
// func New(lri *C.sparkey_logreader) *LogIterator {
//   li := LogIterator{
//     LogReaderInstance: lri,
//     Instance: nil,
//   }
//   li.Create()
//   return &li
// }

// TODO
//  // * sparkey_returncode returncode = sparkey_logiter_create(&myiter, myreader);
// func (li *LogIterator) Create() {
//   C.sparkey_logiter_create(
//     &li.Instance,
//     C.sparkey_logreader(&li.LogReaderInstance))
// }
