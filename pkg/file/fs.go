package file

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

var FILE_PATH *string

func ReadFile(filePath *string, FILE_DATA *[]string) {
	FILE_PATH = filePath
	var file, openError = os.Open(*filePath)
	defer file.Close()
	if openError != nil {
		if errors.Is(openError, os.ErrNotExist) {

			var _, err = os.Create(*filePath)
			if err != nil {
				fmt.Printf("\033[31m %v \033[0m \n", err)
				os.Exit(1)
			}

			*FILE_DATA = make([]string, 0)
			return

		}
		fmt.Printf("\033[31m %v \033[0m \n", openError)
		os.Exit(1)
	}

	var scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		*FILE_DATA = append(*FILE_DATA, scanner.Text())
	}
}

func WriteFile(FILE_DATA *[]string, COMMAND_MESSAGE *string) {
	var file = *FILE_DATA
	var fileString = ""
	for i := 0; i < len(file); i++ {
		fileString += file[i] + "\n"
	}
	var err = os.WriteFile(*FILE_PATH, []byte(fileString), os.ModeAppend)
	if err != nil {
		*COMMAND_MESSAGE = "error saving file"
	} else {
		*COMMAND_MESSAGE = "file saved succesfully"
	}

}
