package main

import (
	_ "goinaction/sample1/matchers"
	"goinaction/sample1/search"
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	search.Run("something")
}
