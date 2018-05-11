package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"
import "./log_writer"

type Sparkey struct {
  Filename string
  LogWriter *log_writer.LogWriter
}

func New(filename string) *Sparkey {
  s := Sparkey{
    Filename: filename,
    LogWriter: log_writer.New(filename),
  }
  return &s
}

func (store *Sparkey) Put(key string, value string) {
  // TODO
  // store.LogWriter.Put(key, value)
}

func (store *Sparkey) Flush() {
  // TODO
}

func (store *Sparkey) Size() (size int) {
  // TODO
  return 0
}
