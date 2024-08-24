package file

import (
	"fmt"
	"path/filepath"
)

func ReadDir(path string, data *[]string) {
	path, err := filepath.Abs(path)
	path = filepath.Dir(path)
	if err != nil {
		fmt.Printf("\033[31m %v \033[0m \n", err)
	}
	fmt.Print(path)
}
