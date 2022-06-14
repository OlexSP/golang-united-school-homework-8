package main

import (
	"io"
	"os"
)

type Arguments map[string]string

func Perform(args Arguments, writer io.Writer) error {
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
