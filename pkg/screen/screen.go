package screen

import (
	"fmt"
	"slices"

	"github.com/gdamore/tcell/v2"
)

type Screen struct {
	Screen     tcell.Screen
	Dimensions map[string]int
	Cursor     map[string]int
}

const (
	X      string = "x"
	Y      string = "y"
	WIDTH  string = "width"
	HEIGHT string = "height"
	START  string = "start"
	END    string = "end"
)

var FILENAME []rune
var FILEDATA []string
var PADDING int = 3
var Cursor = map[string]int{
	X: PADDING,
	Y: 1,
}
var Dimensions = map[string]int{
	WIDTH:  0,
	HEIGHT: 0,
}
var ViewLine = map[string]int{
	START: 0,
	END:   0,
}

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

func (s *Screen) Input() {

	for {

		s.RenderFileName()
		s.RenderCursor()
		s.RenderFileData()
		s.Screen.Show()
		var event = s.Screen.PollEvent()

		switch event := event.(type) {

		case *tcell.EventKey:

			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
				return
			} else if event.Key() == tcell.KeyLeft {

				if Cursor[X] <= PADDING && Cursor[Y] > 1 {
					Cursor[Y] -= 1
					if len(FILEDATA[Cursor[Y]-1]) == 0 {
						Cursor[X] = PADDING
					} else {
						Cursor[X] = len(FILEDATA[Cursor[Y]-1]) + PADDING + 1
					}
				}
				if Cursor[X] > PADDING {
					Cursor[X] = Cursor[X] - 1
				}

			} else if event.Key() == tcell.KeyRight {

				if Cursor[X] >= len(FILEDATA[Cursor[Y]-1])+PADDING {
					if Cursor[Y] >= len(FILEDATA) {
						continue
					}
					Cursor[Y] += 1
					Cursor[X] = PADDING
					continue
				}
				Cursor[X] = Cursor[X] + 1

			} else if event.Key() == tcell.KeyUp && Cursor[Y] > 1 {

				Cursor[Y] = Cursor[Y] - 1
				if Cursor[X] > len(FILEDATA[Cursor[Y]-1]) {
					Cursor[X] = len(FILEDATA[Cursor[Y]-1]) + PADDING
				}

			} else if event.Key() == tcell.KeyDown && Cursor[Y] < len(FILEDATA) {

				if Cursor[X] >= len(FILEDATA[Cursor[Y]])+1 {
					if len(FILEDATA[Cursor[Y]]) == 0 {
						Cursor[X] = PADDING
						Cursor[Y] = Cursor[Y] + 1
						continue
					}
					Cursor[X] = len(FILEDATA[Cursor[Y]]) + 1
				}
				Cursor[Y] = Cursor[Y] + 1

			} else if event.Key() == tcell.KeyEnter {
				if Cursor[X] == len(FILEDATA[Cursor[Y]-1])+PADDING {
					FILEDATA = slices.Insert(FILEDATA, Cursor[Y], "")
					FILEDATA[Cursor[Y]] = ""
					Cursor[Y] += 1
					Cursor[X] = PADDING
				} else {
					var currentLine = FILEDATA[Cursor[Y]-1]
					FILEDATA = slices.Insert(FILEDATA, Cursor[Y], currentLine[Cursor[X]-PADDING:])
					FILEDATA[Cursor[Y]-1] = currentLine[0 : Cursor[X]-PADDING]
					Cursor[Y] += 1
					Cursor[X] = PADDING
				}
				s.Screen.Clear()
			} else if event.Key() == tcell.KeyBackspace {
				if Cursor[X] == PADDING {
					if Cursor[Y] == 1 {
						continue
					}
					var lastData = FILEDATA[Cursor[Y]-1]
					FILEDATA = slices.Delete(FILEDATA, Cursor[Y]-1, Cursor[Y])
					if Cursor[Y]-2 >= 0 {
						Cursor[X] = len(FILEDATA[Cursor[Y]-2]) + PADDING
						FILEDATA[Cursor[Y]-2] += lastData
					}
					Cursor[Y] -= 1
				} else {
					var currentLineData = FILEDATA[Cursor[Y]-1]
					FILEDATA[Cursor[Y]-1] = currentLineData[:Cursor[X]-PADDING-1] + currentLineData[Cursor[X]-PADDING:]
					if Cursor[X] > PADDING {
						Cursor[X] -= 1
					}
				}
				s.Screen.Clear()
			} else {
				var currentLineData = FILEDATA[Cursor[Y]-1]
				FILEDATA[Cursor[Y]-1] = currentLineData[:Cursor[X]-PADDING] + string(event.Rune()) + currentLineData[Cursor[X]-PADDING:]
				Cursor[X] += 1
			}
		}
	}
}

func (s *Screen) RenderCursor() {
	s.Screen.ShowCursor(Cursor[X], Cursor[Y])
}

func (s *Screen) RenderFileName() {
	for i := 0; i < len(FILENAME); i++ {
		s.Screen.SetContent(i, 0, FILENAME[i], nil, tcell.StyleDefault)
	}
}

func (s *Screen) RenderFileData() {
	for y := 0; y < Dimensions[HEIGHT]-2; y++ {
		if y < len(FILEDATA) {
			var lineRune = []rune(fmt.Sprintf("%02d", y+1) + " " + FILEDATA[y])
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
