package main

import (
	"io/ioutil"
	"os"
)

func main() {
	os.Exit(run())
}

func run() int {
	b, err := ioutil.ReadFile("input.txt")
	if err != nil {
		return 1
	}
	f := CreateField(b)
	f.Tick()
	ShowField(f)
	return 0
}
