package main

import (
  "fmt"
  "./sparkey"
)

func main () {
  s := sparkey.New("sparkey_db", sparkey.COMPRESSION_NONE, 1024)

  s.Put("first", "Hello")
  s.Put("second", "Worlb")
  s.Put("third", "Goodbye")
  s.Put("fourth", "EOM")
  s.Delete("third")
  s.Flush()

  fmt.Println("\n{")
  s.ForEach(func(k, v string) {
    fmt.Printf("  %s => %s\n", k, v)
  })
  fmt.Println("}\n")


  fmt.Printf("Sparkey info:\n\n")
  fmt.Printf("  Basename:                  %s\n", s.Basename)
  fmt.Printf("  Size:                      %d\n", s.Size())
  fmt.Printf("  LogWriter.CompressionType: %d\n", s.CompressionType)
  fmt.Printf("  LogWriter.BlockSize:       %d\n", s.BlockSize)

  s.Close()
}
