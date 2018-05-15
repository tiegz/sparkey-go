package main

import (
  "fmt"
  "./sparkey"
)

func main () {
  s := sparkey.New("sparkey_db")

  // s.Put("first", "Hello")
  // s.Put("second", "Worlb")
  // s.Put("third", "Goodbye")
  // s.Delete("third")
  // s.Flush()

  // s.GetAll()

  fmt.Printf("Hello, Worlb. Sparkey info:\n  Basename: %d\n", s.Basename)
    //\n  LogWriter.Basename: %s\n  LogWriter.CompressionType: %s\n  LogWriter.BlockSize: %s\n  %s\n",
  //   s.Size(),
  //   s.Basename,
  //   s.LogWriter.Basename,
  //   s.LogWriter.CompressionType,
  //   s.LogWriter.BlockSize)

  s.Close()
}
