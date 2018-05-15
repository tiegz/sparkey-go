package sparkey

// #include <stdio.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/lib
// #cgo LDFLAGS: -L/usr/local/lib -lsparkey
// #include <./sparkey.h>
import "C"

type Sparkey struct {
  Basename string
  LogWriter *C.sparkey_logwriter
  LogReader *C.sparkey_logreader
  HashWriter *HashWriter // there is no HashWriter class in the C API for Sparkey, it's just a thing
  HashReader *C.sparkey_hashreader
}

func New(basename string) *Sparkey {
  s := Sparkey{
    Basename: basename,
    LogWriter: NewLogWriter(basename, COMPRESSION_NONE, 1024),
    LogReader: NewLogReader(basename),
    HashWriter: NewHashWriter(basename),
    HashReader: NewHashReader(basename),
  }
  return &s
}

func (store *Sparkey) Put(key string, value string) {
  // store.LogWriter.Put(key, value)
}

func (store *Sparkey) Delete(key string) {
  // store.LogWriter.Delete(key)
}



func (store *Sparkey) GetAll() {
  // li := log_iterator.New(store.LogReader.Instance)
  // li = log_iterator.New(store.LogWriter)
  //   iterator = hash_reader.seek(key)

  //   return unless iterator.active?

  //   iterator.get_value
  // end


}
 // * The logreader is not useful by itself. You also need a sparkey_logiter to iterate through the entries.
 // * This is a highly mutable struct and should not be shared between threads. It is not threadsafe.
 // *
 // * Here is a basic workflow for iterating through all entries in a logfile:
 // * - Create a logreader
 // * \code
 // * sparkey_logreader *myreader;
 // * sparkey_returncode returncode = sparkey_logreader_open(&myreader, "mylog.spl");
 // * \endcode
 // * - Create a logiter
 // * \code
 // * sparkey_logiter *myiter;
 // * sparkey_returncode returncode = sparkey_logiter_create(&myiter, myreader);
 // * \endcode
 // * - Perform the iteration:
 // * \code
 // * while (1) {
 // *   sparkey_returncode returncode = sparkey_logiter_next(myiter, myreader);
 // *   // TODO: check the returncode
 // *   if (sparkey_logiter_state(myiter) != SPARKEY_ITER_ACTIVE) {
 // *     break;
 // *   }
 // *   uint64_t wanted_keylen = sparkey_logiter_keylen(myiter);
 // *   uint8_t *keybuf = malloc(wanted_keylen);
 // *   uint64_t actual_keylen;
 // *   returncode = sparkey_logiter_fill_key(myiter, myreader, wanted_keylen, keybuf, &actual_keylen);
 // *   // TODO: check the returncode
 // *   // TODO: assert actual_keylen == wanted_keylen
 // *   uint64_t wanted_valuelen = sparkey_logiter_valuelen(myiter);
 // *   uint8_t *valuebuf = malloc(wanted_valuelen);
 // *   uint64_t actual_valuelen;
 // *   returncode = sparkey_logiter_fill_value(myiter, myreader, wanted_valuelen, valuebuf, &actual_valuelen);
 // *   // TODO: check the returncode
 // *   // TODO: assert actual_valuelen == wanted_valuelen
 // *   // Do stuff with key and value
 // *   free(keybuf);
 // *   free(valuebuf);
 // * }
 // * \endcode
 // * Note that you have to allocate memory for the key and value manually - Sparkey does not allocate memory except for when
 // * creating readers, writers and iterators.

func (store *Sparkey) Flush() {
  store.LogWriter.Flush()
}

func (store *Sparkey) Close() {
  C.sparkey_logwriter_close(&store.LogWriter)
  C.sparkey_logreader_close(&store.LogReader)
  C.sparkey_hash_close(&store.HashReader)
}

func (store *Sparkey) Size() (size int) {
  // TODO
  return 0
}
