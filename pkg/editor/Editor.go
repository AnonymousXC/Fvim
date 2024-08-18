package editor

import (
	"fmt"
	"fvim/pkg/screen"
)

type Editor struct {
	MODE     string
	filePath string
	data     []string
}

const (
	MODE_NORMAL = "NORMAL"
	MODE_INSERT = "INSERT"
	MODE_CMD    = "COMMAND_LINE"
)

func EditorInit(filepath *string, data *[]string) {

	var editor = Editor{filePath: *filepath, MODE: MODE_NORMAL, data: *data}
	var currentScreen, screenError = screen.CreateScreen()
	defer currentScreen.Close()

	if screenError != nil {
		fmt.Printf("\033[31m %v \033[0m", screenError)
	}

	fmt.Print(editor)
	currentScreen.Size()
	currentScreen.SetFileName(*filepath)
	currentScreen.SetFileData(data)
	currentScreen.Input()
}
