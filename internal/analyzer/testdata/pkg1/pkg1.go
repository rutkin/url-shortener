package main

import (
	"os"
)

func main() {
	os.Exit(1) // want "expression os.exit deprecated in main function"
}
