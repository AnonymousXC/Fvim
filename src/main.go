package main

import (
	"fmt"
	"fvim/pkg/editor"
	"fvim/pkg/file"
	"os"
)

var FILE_PATH string
var FILE_DATA []string
var FOLDER_DATA []string

func main() {

	if len(os.Args) < 2 {
		fmt.Print("\033[31m No path found. \033[0m")
		os.Exit(1)
	}

	FILE_PATH = os.Args[1]
	file.ReadFile(&FILE_PATH, &FILE_DATA)
	editor.EditorInit(&FILE_PATH, &FILE_DATA)
}
