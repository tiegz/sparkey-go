package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"
import "unsafe"
import "fmt"


type compressionType int
const (
  COMPRESSION_NONE compressionType = iota
  COMPRESSION_SNAPPY
)

const (
  ITERATOR_STATE_NEW     = 0
  ITERATOR_STATE_ACTIVE  = 1
  ITERATOR_STATE_CLOSED  = 2
  ITERATOR_STATE_INVALID = 3
)

type forEachType uint
const (
  FOR_EACH_LOG forEachType = iota
  FOR_EACH_HASH
)

type Sparkey struct {
  Basename string
  LogFilename string
  IndexFilename string
  CompressionType compressionType
  BlockSize int
  LogWriter   *C.sparkey_logwriter
  LogReader   *C.sparkey_logreader
  HashWriter  *HashWriter // there's actually no HashWriter struct in Sparkey's C API -- HashWriter is unused except for the factory method
  HashReader  *C.sparkey_hashreader
}

func New(basename string, compression_type compressionType, block_size int) *Sparkey {
  log_filename := basename + ".spl"
  index_filename := basename + ".spi"

  s := Sparkey{
    Basename: basename,
    LogFilename: log_filename,
    IndexFilename: index_filename,
    CompressionType: compression_type,
    BlockSize: block_size,
    LogWriter: NewLogWriter(log_filename, compression_type, 1024),
    LogReader: OpenLogReader(log_filename),
    HashWriter: NewHashWriter(log_filename, index_filename),
    HashReader: OpenHashReader(log_filename, index_filename),
  }

  return &s
}

func Open(basename string) *Sparkey {
  log_filename := basename + ".spl"
  index_filename := basename + ".spi"

  s := Sparkey{
    Basename: basename,
    LogFilename: log_filename,
    IndexFilename: index_filename,
    LogWriter: OpenLogWriter(log_filename),
    LogReader: OpenLogReader(log_filename),
    HashWriter: NewHashWriter(log_filename, index_filename),
    HashReader: OpenHashReader(log_filename, index_filename),
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

  fmt.Printf("Put:\t\t\t%s %d %s %d, Return Code: %d\n", key, len(key), value, len(value), return_code)
}

func (store *Sparkey) Delete(key string) {
  cKey := (*C.uchar)(unsafe.Pointer(C.CString(key)))

  return_code := C.sparkey_logwriter_delete(
    store.LogWriter,
    C.ulonglong((len(key))),
    cKey)

  fmt.Printf("Delete:\t\t\t%s %d, Return Code: %d\n", key, len(key), return_code)
}

func (store *Sparkey) Flush() {
  // Flush logwriter
  return_code := C.sparkey_logwriter_flush(store.LogWriter)
  fmt.Printf("Flush LW:\t\tReturn Code %d\n", return_code)

  // Reset to flush cached headers
  return_code = C.sparkey_logreader_open(&store.LogReader, C.CString(store.LogFilename))
  fmt.Printf("Flush LR:\t\tReturn Code %d\n", return_code)

  return_code = C.sparkey_hash_write(C.CString(store.IndexFilename), C.CString(store.LogFilename), 0)
  fmt.Printf("Hash write:\t\tReturn Code %d\n", return_code)

  // TODO do we really  need to reopen hash reader?
  return_code = C.sparkey_hash_open(&store.HashReader, C.CString(store.IndexFilename), C.CString(store.LogFilename))
  fmt.Printf("Hash open:\t\tReturn Code %d\n", return_code)
}


func (store *Sparkey) Close() {
  C.sparkey_logwriter_close(&store.LogWriter)
  C.sparkey_logreader_close(&store.LogReader)
  C.sparkey_hash_close(&store.HashReader)
}

func (store *Sparkey) Size() (size uint64) {
  return uint64(C.sparkey_hash_numentries(store.HashReader))
}

func (store *Sparkey) Get(k string) (v string) {
  li := NewLogIterator(store.LogReader)
  defer C.sparkey_logiter_close(&li)
  cKey := (*C.uchar)(unsafe.Pointer(C.CString(k)))
  defer C.free(unsafe.Pointer(cKey))
  var wanted_valuelen C.ulonglong
  var actual_valuelen C.ulonglong
  var valuebuf *C.uchar
  defer C.free(unsafe.Pointer(valuebuf))
  var return_code C.sparkey_returncode

  // TODO handle return code
  return_code = C.sparkey_hash_get(
    store.HashReader,
    cKey,
    C.ulonglong(len(k)),
    li)

  if C.sparkey_logiter_state(li) != ITERATOR_STATE_ACTIVE {
    fmt.Printf("Get:\t\t\t'%s' not found! \n", k)
    // TODO Entry not found
  } else {
    wanted_valuelen = C.sparkey_logiter_valuelen(li)
    valuebuf = (*C.uchar)(unsafe.Pointer(C.CString("")))
    return_code = C.sparkey_logiter_fill_value(
      li,
      store.LogReader,
      wanted_valuelen,
      valuebuf,
      &actual_valuelen)

    if return_code != C.SPARKEY_SUCCESS {
      fmt.Printf("Breaking, sparkey_logiter_fill_value returned %d return code.", C.int(return_code))
      return "TODO return error here"
    }
    if (wanted_valuelen != actual_valuelen) {
      fmt.Printf("Breaking, sparkey_logiter_fill_value returned %d length instead of %d length.", C.int(actual_valuelen), C.int(wanted_valuelen))
      return "TODO return error here"
    }

    valuebuf_value := (*C.char)(unsafe.Pointer(valuebuf))
    v = C.GoString(valuebuf_value)
  }

  fmt.Printf("Get:\t\t%s %s, Return Code: %d\n", k, v, return_code)

  return v
}

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
    return_code = C.sparkey_logiter_fill_value(
      li,
      store.LogReader,
      wanted_valuelen,
      valuebuf,
      &actual_valuelen)
    if return_code != C.SPARKEY_SUCCESS {
      fmt.Printf("Breaking, sparkey_logiter_fill_value returned %d return code.", C.int(return_code))
      break // TODO return error instead of breaking
    }
    if (wanted_valuelen != actual_valuelen) {
      fmt.Printf("Breaking, sparkey_logiter_fill_value returned %d length instead of %d length.", C.int(actual_valuelen), C.int(wanted_valuelen))
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

func (store *Sparkey) MaxKeyLen() (l uint64) {
  l = uint64(C.sparkey_logreader_maxkeylen(store.LogReader))
  return l
}

func (store *Sparkey) MaxValueLen() (l uint64) {
  l = uint64(C.sparkey_logreader_maxvaluelen(store.LogReader))
  return l
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
