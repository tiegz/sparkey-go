package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"
import "unsafe"
import "fmt"

type Sparkey struct {
  Basename string
  CompressionType int
  BlockSize int
  LogWriter   *C.sparkey_logwriter
  LogReader   *C.sparkey_logreader
  // LogIterator *C.sparkey_logiter
  HashWriter  *HashWriter // there's actually no HashWriter struct in Sparkey's C API -- HashWriter is unused except for the factory method
  HashReader  *C.sparkey_hashreader
}

// TODO use iota instead?
const COMPRESSION_NONE int   = 0
const COMPRESSION_SNAPPY int = 1

// TODO use iota instead?
const ITERATOR_STATE_NEW C.sparkey_iter_state = 0
const ITERATOR_STATE_ACTIVE C.sparkey_iter_state = 1
const ITERATOR_STATE_CLOSED C.sparkey_iter_state = 2
const ITERATOR_STATE_INVALID C.sparkey_iter_state = 3

func New(basename string, compression_type int, block_size int) *Sparkey {
  s := Sparkey{
    Basename: basename,
    CompressionType: compression_type,
    BlockSize: block_size,
    LogWriter: NewLogWriter(basename, compression_type, 1024),
    LogReader: OpenLogReader(basename),
    // LogIterator: NewLogIterator(???),
    // HashIterator: NewHashIterator(???),
    HashWriter: NewHashWriter(basename),
    HashReader: OpenHashReader(basename),
  }
  return &s
}

func Open(basename string) *Sparkey {
  s := Sparkey{
    Basename: basename,
    LogWriter: OpenLogWriter(basename),
    LogReader: OpenLogReader(basename),
    // LogIterator: NewLogIterator(???),
    // HashIterator: NewHashIterator(???),
    HashWriter: NewHashWriter(basename),
    HashReader: OpenHashReader(basename),
  }

  return &s
}

func (store *Sparkey) Put(key string, value string) {
  cKey := (*C.uchar)(unsafe.Pointer(C.CString(key)))
  cValue := (*C.uchar)(unsafe.Pointer(C.CString(value)))

  defer C.free(unsafe.Pointer(cKey))
  defer C.free(unsafe.Pointer(cValue))

  return_code := C.sparkey_logwriter_put(
    store.LogWriter,
    C.ulonglong(len(key)),
    cKey,
    C.ulonglong(len(value)),
    cValue)

  fmt.Printf("Put:\t\t%s %d %s %d, Return Code: %d\n", key, len(key), value, len(value), return_code)
}

func (store *Sparkey) Delete(key string) {
  cKey := (*C.uchar)(unsafe.Pointer(C.CString(key)))

  return_code := C.sparkey_logwriter_delete(
    store.LogWriter,
    C.ulonglong((len(key))),
    cKey)

  fmt.Printf("Delete:\t\t%s %d, Return Code: %d\n", key, len(key), return_code)
}

func (store *Sparkey) Flush() {
  // TODO could we store these in store instead of building them?
  log_filename := store.Basename + ".spl"
  index_filename := store.Basename + ".spi"

  // Flush logwriter
  return_code := C.sparkey_logwriter_flush(store.LogWriter)
  fmt.Printf("Flush logwriter\t\tReturn Code %d\n", return_code)

  // Reset to flush cached headers
  return_code = C.sparkey_logreader_open(&store.LogReader, C.CString(log_filename))
  fmt.Printf("Flush logreader\t\tReturn Code %d\n", return_code)

  return_code = C.sparkey_hash_write(C.CString(index_filename), C.CString(log_filename), 0)
  fmt.Printf("Hash write\t\tReturn Code %d\n", return_code)

  // TODO do we really  need to reopen hash reader?
  return_code = C.sparkey_hash_open(&store.HashReader, C.CString(index_filename), C.CString(log_filename))
  fmt.Printf("Hash open\t\tReturn Code %d\n", return_code)
}


func (store *Sparkey) Close() {
  C.sparkey_logwriter_close(&store.LogWriter)
  C.sparkey_logreader_close(&store.LogReader)
  C.sparkey_hash_close(&store.HashReader)
}

func (store *Sparkey) Size() (size uint64) {
  return uint64(C.sparkey_hash_numentries(store.HashReader))
}

type forEachType uint
const (
  FOR_EACH_LOG forEachType = iota
  FOR_EACH_HASH
)

func (store *Sparkey) ForEachHash(fn func(k, v string)) {
  var li *C.sparkey_logiter = NewLogIterator(store.LogReader)

  store.forEach(FOR_EACH_HASH, li, fn)

  C.sparkey_logiter_close(&li)
}

func (store *Sparkey) ForEachLog(fn func(k, v string)) {
  var li *C.sparkey_logiter = NewLogIterator(store.LogReader)

  store.forEach(FOR_EACH_LOG, li, fn)

  C.sparkey_logiter_close(&li)
}


func (store *Sparkey) forEach(t forEachType, li *C.sparkey_logiter, fn func(k, v string)) {
  var lis C.sparkey_iter_state = C.sparkey_logiter_state(li)

  var wanted_keylen C.ulonglong
  var actual_keylen C.ulonglong
  var wanted_valuelen C.ulonglong
  var actual_valuelen C.ulonglong
  var keybuf *C.uchar
  var valuebuf *C.uchar
  var return_code C.sparkey_returncode
  defer C.free(unsafe.Pointer(keybuf))
  defer C.free(unsafe.Pointer(valuebuf))

  for lis == ITERATOR_STATE_NEW || lis == ITERATOR_STATE_ACTIVE {
    // TODO: check the returncode
    switch t {
    case FOR_EACH_HASH:
      C.sparkey_logiter_hashnext(li, store.HashReader);
    case FOR_EACH_LOG:
      C.sparkey_logiter_next(li, store.LogReader);
    }

    wanted_keylen = C.sparkey_logiter_keylen(li)
    keybuf = (*C.uchar)(unsafe.Pointer(C.CString("")))
    return_code = C.sparkey_logiter_fill_key(li, store.LogReader, wanted_keylen, keybuf, &actual_keylen)
    if return_code != C.SPARKEY_SUCCESS {
      fmt.Printf("Breaking, sparkey_logiter_fill_key returned %d return code.", C.int(return_code))
      break // TODO return error instead of breaking
    }
    if (wanted_keylen != actual_keylen) {
      fmt.Printf("Breaking, sparkey_logiter_fill_key returned %d length instead of %d length.", C.int(actual_keylen), C.int(wanted_keylen))
      break // TODO return error instead of breaking
    }

    wanted_valuelen = C.sparkey_logiter_valuelen(li)
    valuebuf = (*C.uchar)(unsafe.Pointer(C.CString("")))
    return_code = C.sparkey_logiter_fill_value(li, store.LogReader, wanted_valuelen, valuebuf, &actual_valuelen)
    if return_code != C.SPARKEY_SUCCESS {
      fmt.Printf("Breaking, sparkey_logiter_fill_value returned %d return code.", C.int(return_code))
      break // TODO return error instead of breaking
    }
    if (wanted_keylen != actual_keylen) {
      fmt.Printf("Breaking, sparkey_logiter_fill_value returned %d length instead of %d length.", C.int(actual_keylen), C.int(wanted_valuelen))
      break // TODO return error instead of breaking
    }

    lis = C.sparkey_logiter_state(li)
    if lis == ITERATOR_STATE_ACTIVE {
      keybuf_value   := (*C.char)(unsafe.Pointer(keybuf))
      valuebuf_value := (*C.char)(unsafe.Pointer(valuebuf))
      fn(C.GoString(keybuf_value), C.GoString(valuebuf_value))
    }

    keybuf = nil
    valuebuf = nil
  }
}

func (store *Sparkey) PrettyPrintHash() {
  fmt.Println("\n{")
  store.ForEachHash(func(k, v string) {
    fmt.Printf("  %s => %s\n", k, v)
  })
  fmt.Println("}\n")
}

func (store *Sparkey) PrettyPrintLog() {
  fmt.Println("\n{")
  store.ForEachLog(func(k, v string) {
    fmt.Printf("  %s => %s\n", k, v)
  })
  fmt.Println("}\n")
}

func (store *Sparkey) GetAll() {
  // li := log_iterator.New(store.LogReader.Instance)
  // li = log_iterator.New(store.LogWriter)
  //   iterator = hash_reader.seek(key)

  //   return unless iterator.active?

  //   iterator.get_value
  // end
}
