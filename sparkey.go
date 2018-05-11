package main

import (
  "fmt"
  "./sparkey"
)

func main () {
  s := sparkey.New("sparkey_db.spl")

  s.Put("first", "Hello")
  s.Put("second", "Worlb")
  s.Put("third", "Goodbye")
  s.Flush()

  fmt.Printf("Hello, Worlb. Sparkey info:\n  Filename: %d\n  LogWriter.Filename: %s\n  LogWriter.CompressionType: %s\n  LogWriter.BlockSize: %s\n  %s\n",
    s.Size(),
    s.Filename,
    s.LogWriter.Filename,
    s.LogWriter.CompressionType,
    s.LogWriter.BlockSize)
}
