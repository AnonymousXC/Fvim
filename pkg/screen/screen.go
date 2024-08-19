package screen

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type Screen struct {
	Screen     tcell.Screen
	Dimensions map[string]int
	Cursor     map[string]int
}

const (
	X           string = "x"
	Y           string = "y"
	WIDTH       string = "width"
	HEIGHT      string = "height"
	START       string = "start"
	END         string = "end"
	LEFT        string = "left"
	BOTTOM      string = "bottom"
	TOP         string = "top"
	MODE_NORMAL string = "NORMAL"
	MODE_INSERT string = "INSERT"
	MODE_CMD    string = "COMMAND_LINE"
)

var (
	FILENAME    []rune
	FILEDATA    []string
	MODE        = MODE_NORMAL
	COMMAND     = ""
	CMD_MESSAGE = ""
	PADDING     = map[string]int{
		LEFT:   3,
		BOTTOM: 2,
		TOP:    2,
	}
	Cursor = map[string]int{
		X: PADDING[LEFT],
		Y: 1,
	}
	Dimensions = map[string]int{
		WIDTH:  0,
		HEIGHT: 0,
	}
	ViewLine = map[string]int{
		START: 0,
		END:   Dimensions[HEIGHT] - 2,
	}
)

func CreateScreen() (*Screen, error) {

	var screen, err = tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	var initErr = screen.Init()
	if initErr != nil {
		return nil, err
	}

	var s = Screen{Screen: screen, Dimensions: Dimensions, Cursor: Cursor}
	return &s, nil
}

func (s *Screen) Close() {
	s.Screen.Fini()
}

func (s *Screen) RenderCursor() {
	if MODE == MODE_INSERT {
		s.Screen.ShowCursor(Cursor[X], Cursor[Y])
	}
}

func (s *Screen) RenderFileName() {
	for i := 0; i < len(FILENAME); i++ {
		s.Screen.SetContent(i, 0, FILENAME[i], nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}
}

func (s *Screen) RenderFileData() {
	for y, fileLine := 0, ViewLine[START]; y < Dimensions[HEIGHT]-PADDING[BOTTOM]; y, fileLine = y+1, fileLine+1 {
		if fileLine < len(FILEDATA) {
			var lineRune = []rune(fmt.Sprintf("%02d", fileLine+1) + " " + FILEDATA[fileLine])
			for x := 0; x < Dimensions[WIDTH]; x++ {
				if x < len(lineRune) {
					s.Screen.SetContent(x, y+1, lineRune[x], nil, tcell.StyleDefault)
				}
			}
		} else {
			s.Screen.SetContent(0, y+1, '>', nil, tcell.StyleDefault)
		}
	}
}

func (s *Screen) SetFileName(filename string) {
	filename = "Editing " + filename
	FILENAME = []rune(filename)
}

func (s *Screen) SetFileData(filedata *[]string) {
	FILEDATA = *filedata
}

func (s *Screen) Size() {
	Dimensions[WIDTH], Dimensions[HEIGHT] = s.Screen.Size()
}

func (s *Screen) RenderCommandBox() {
	for x := 0; x < Dimensions[WIDTH]; x++ {
		s.Screen.SetContent(x, Dimensions[HEIGHT]-PADDING[BOTTOM], tcell.RuneHLine, nil, tcell.StyleDefault.Foreground(tcell.ColorTeal))
	}
}

func (s *Screen) SetMode(mode string) {
	MODE = mode
	s.Screen.Clear()
}

func (s *Screen) RenderCommandText(cmd string, style tcell.Style) {
	var cmdRund = []rune(cmd)
	for x := 0; x < len(cmdRund); x++ {
		s.Screen.SetContent(x, Dimensions[HEIGHT]-PADDING[BOTTOM]+1, cmdRund[x], nil, style)
	}
}

func (s *Screen) RenderCommand() {
	s.RenderCommandBox()
	if MODE == MODE_INSERT {
		s.RenderCommandText("press esc to exit insert mode", tcell.StyleDefault.Foreground(tcell.ColorTeal))
	}
	if MODE == MODE_NORMAL {
		if COMMAND == "" && CMD_MESSAGE == "" {
			s.RenderCommandText("press i to enter insert mode", tcell.StyleDefault.Foreground(tcell.ColorTeal))
		} else if COMMAND != "" {
			s.RenderCommandText(COMMAND, tcell.StyleDefault.Foreground(tcell.ColorTeal))
		} else {
			s.RenderCommandText(CMD_MESSAGE, tcell.StyleDefault.Foreground(tcell.ColorRed))
		}
	}
}
