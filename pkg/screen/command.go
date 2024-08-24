package screen

import (
	"fvim/pkg/file"
	"reflect"
	"strings"
)

func HandleCommand() {
	if COMMAND == ":w" {
		COMMAND = ""
		CMD_MESSAGE = "saving..."
		file.WriteFile(&FILEDATA, &CMD_MESSAGE)
	} else {
		COMMAND = ""
		CMD_MESSAGE = "command not found"
	}
}

func checkBeforeExit() bool {
	var filepath string
	var filedata []string
	for i := 0; i < len(FILENAME); i++ {
		filepath += string(FILENAME[i])
	}
	filepath = strings.Replace(filepath, "Editing ", "", -1)
	file.ReadFile(&filepath, &filedata)
	if reflect.DeepEqual(filedata, FILEDATA) {
		return true
	} else {
		filedata = nil
		filepath = ""
		return false
	}
}
