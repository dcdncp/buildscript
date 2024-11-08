package main

import (
	"bscript/runtime"
	"fmt"
	"os"
)

func main() {
	args := os.Args
	var path string
	if len(args) == 0 {
		path = "build.script"
	} else {
		path = args[1]
	}
	errors := runtime.RunFile(path)
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
	}
}
