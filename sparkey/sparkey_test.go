package sparkey

import (
	"fmt"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	s1 := New("my-first-sparkey", COMPRESSION_NONE, 1024)
	s2 := New("my-second-sparkey", COMPRESSION_SNAPPY, 1024)

	assertSparkeyExists(t, "my-first-sparkey")
	assertSparkeyExists(t, "my-second-sparkey")

	s1.Close()
	s2.Close()

	s1.DeleteSparkey()
	s2.DeleteSparkey()
}

func TestOpen(t *testing.T) {
	s := New("my-first-sparkey", COMPRESSION_NONE, 1024)
	s.Put("first key", "first value")
	s.Flush()
	s.Close()

	s = Open("my-first-sparkey")
	assertSparkeyExists(t, "my-first-sparkey")
	assertSparkeyKeyValue(t, s, "first key", "first value")

	s.DeleteSparkey()
}

func TestPut(t *testing.T) {
	s := New("my-first-sparkey", COMPRESSION_NONE, 1024)

	s.Put("first key", "first value")
	s.Put("second key", "second value")
	s.Put("third key", "third value")
	if size := s.Size(); size > 1 {
		t.Errorf("Expected sparkey size to be 0, but it has a size of %d", size)
	}

	s.Flush()
	if size := s.Size(); size != 3 {
		t.Errorf("Expected sparkey size to be 3, but it has a size of %d", size)
	}
	assertSparkeyKeyValue(t, s, "first key", "first value")
	assertSparkeyKeyValue(t, s, "second key", "second value")
	assertSparkeyKeyValue(t, s, "third key", "third value")

	s.Close()
	s.DeleteSparkey()
}

func TestGet(t *testing.T) {
	s := New("my-first-sparkey", COMPRESSION_NONE, 1024)

	s.Put("first key", "first value")
	s.Put("second key", "second value")
	s.Flush()

	if val, err := s.Get("first key"); err != nil {
		t.Errorf("Unexpected sparkey value: %s", val)
	}

	s.Close()
	s.DeleteSparkey()
}

func TestDelete(t *testing.T) {
	s := New("my-first-sparkey", COMPRESSION_NONE, 1024)

	s.Put("first key", "first value")
	s.Put("second key", "second value")
	s.Delete("first key")
	s.Flush()
	assertSparkeyNilValue(t, s, "first key")

	s.Close()
	s.DeleteSparkey()
}

func assertSparkeyNilValue(t *testing.T, s *Sparkey, k string) {
	if v, err := s.Get(k); err == nil {
		fmt.Printf("Expected key %s to be nil, but was: %s\n", k, v)
	}
}

func assertSparkeyKeyValue(t *testing.T, s *Sparkey, k string, v string) {
	if val, err := s.Get(k); err != nil {
		t.Errorf(err.Error())
	} else if val != v {
		t.Errorf("Unexpected sparkey value for key %s: %s", k, val)
	}
}

func assertSparkeyExists(t *testing.T, basename string) {
	log_filename := basename + ".spl"
	index_filename := basename + ".spi"

	if _, err := os.Stat(log_filename); os.IsNotExist(err) {
		t.Errorf("Expected %s log file to exist, but it does not exist", log_filename)
	}

	if _, err := os.Stat(index_filename); os.IsNotExist(err) {
		t.Errorf("Expected %s index file to exist, but it does not exist", index_filename)
	}
}
