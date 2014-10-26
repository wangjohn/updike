package main

import (
  "github.com/wangjohn/updike/philarios"
  "log"
)

func main() {
  words, err := philarios.AlternativeWords("amazing", 10)
  if err != nil {
    log.Fatal("Fatal Error: %s", err.Error())
  }

  log.Print("Words: %s", words)
}
