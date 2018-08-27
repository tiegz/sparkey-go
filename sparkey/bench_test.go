package sparkey

import (
  "testing"
  "fmt"
)

func BenchmarkPut(b *testing.B) {
  b.Run("Put", func (b *testing.B) {
    s := New("sparkey-bench-db-put", COMPRESSION_NONE, 1024)

    for i := 0; i < b.N; i++ {
      s.Put(fmt.Sprintf("key_%d", i), "foo")
    }
    s.Flush()
    s.Close()
    s.DeleteSparkey()
  })
}

func BenchmarkGet(b *testing.B) {
  s := New("sparkey-bench-db-get", COMPRESSION_NONE, 1024)

  for i := 0; i < b.N; i++ {
    s.Put(fmt.Sprintf("key_%d", i), "foo")
  }
  s.Flush()

  b.Run("Get", func (b *testing.B) {
    for i := 0; i < b.N; i++ {
      s.Get(fmt.Sprintf("key_%d", i))
    }
  })
  s.Close()
  s.DeleteSparkey()
}
