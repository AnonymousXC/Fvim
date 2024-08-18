package file

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func ReadFile(filePath *string, FILE_DATA *[]string) {
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
