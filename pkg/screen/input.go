package screen

import (
	"fvim/pkg/file"
	"slices"

	"github.com/gdamore/tcell/v2"
)

func (s *Screen) Input() {

	for {

		s.RenderFileName()
		s.RenderFileData()
		s.RenderCommand()
		s.RenderCursor()
		s.Screen.Show()
		var event = s.Screen.PollEvent()

		switch event := event.(type) {

		case *tcell.EventKey:

			if event.Rune() == 'i' && MODE != MODE_INSERT {
				s.SetMode(MODE_INSERT)
				continue
			}

			if MODE == MODE_NORMAL {
				CMD_MESSAGE = ""
				if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyDelete || event.Key() == tcell.KeyEsc {
					if len(COMMAND) > 0 {
						COMMAND = COMMAND[0 : len(COMMAND)-1]
					}
				} else if event.Key() == tcell.KeyEnter {
					if COMMAND == ":q" {
						return
					} else if COMMAND == ":w" {
						COMMAND = ""
						CMD_MESSAGE = "saving..."
						file.WriteFile(&FILEDATA, &CMD_MESSAGE)
					} else {
						COMMAND = ""
						CMD_MESSAGE = "command not found"
					}
				} else {
					COMMAND += string(event.Rune())
				}
				s.Screen.Clear()
			}

			if event.Key() == tcell.KeyEsc && MODE == MODE_INSERT {
				s.SetMode(MODE_NORMAL)
				s.Screen.HideCursor()
			}

			if MODE != MODE_INSERT {
				continue
			}

			if event.Key() == tcell.KeyLeft {

				if Cursor[X] <= PADDING[LEFT] && Cursor[Y] > PADDING[TOP] {
					Cursor[Y] -= 1
					if len(FILEDATA[Cursor[Y]-1]) == 0 {
						Cursor[X] = PADDING[LEFT]
					} else {
						Cursor[X] = len(FILEDATA[Cursor[Y]-1]) + PADDING[LEFT] + 1
					}
				}
				if Cursor[X] > PADDING[LEFT] {
					Cursor[X] = Cursor[X] - 1
				}

			} else if event.Key() == tcell.KeyRight {

				if Cursor[X] >= len(FILEDATA[Cursor[Y]-1])+PADDING[LEFT] {
					if Cursor[Y] >= len(FILEDATA) {
						continue
					}
					Cursor[Y] += 1
					Cursor[X] = PADDING[LEFT]
					continue
				}
				Cursor[X] = Cursor[X] + 1

			} else if event.Key() == tcell.KeyUp && Cursor[Y] > 1 {

				Cursor[Y] = Cursor[Y] - 1
				if Cursor[X] > len(FILEDATA[Cursor[Y]-1]) {
					Cursor[X] = len(FILEDATA[Cursor[Y]-1]) + PADDING[LEFT]
				}

			} else if event.Key() == tcell.KeyDown {

				if Cursor[Y] >= len(FILEDATA) {
					continue
				}

				if Cursor[X] >= len(FILEDATA[Cursor[Y]])+1 {
					if len(FILEDATA[Cursor[Y]]) == 0 {
						Cursor[X] = PADDING[LEFT]
						Cursor[Y] = Cursor[Y] + 1
						continue
					}
					Cursor[X] = len(FILEDATA[Cursor[Y]]) + 1
				}
				Cursor[Y] = Cursor[Y] + 1

			} else if event.Key() == tcell.KeyEnter {
				if Cursor[X] == len(FILEDATA[Cursor[Y]-1])+PADDING[LEFT] {
					FILEDATA = slices.Insert(FILEDATA, Cursor[Y], "")
					FILEDATA[Cursor[Y]] = ""
					Cursor[Y] += 1
					Cursor[X] = PADDING[LEFT]
				} else {
					var currentLine = FILEDATA[Cursor[Y]-1]
					FILEDATA = slices.Insert(FILEDATA, Cursor[Y], currentLine[Cursor[X]-PADDING[LEFT]:])
					FILEDATA[Cursor[Y]-1] = currentLine[0 : Cursor[X]-PADDING[LEFT]]
					Cursor[Y] += 1
					Cursor[X] = PADDING[LEFT]
				}
				s.Screen.Clear()
			} else if event.Key() == tcell.KeyBackspace {
				if Cursor[X] == PADDING[LEFT] {
					if Cursor[Y] == 1 {
						continue
					}
					var lastData = FILEDATA[Cursor[Y]-1]
					FILEDATA = slices.Delete(FILEDATA, Cursor[Y]-1, Cursor[Y])
					if Cursor[Y]-2 >= 0 {
						Cursor[X] = len(FILEDATA[Cursor[Y]-2]) + PADDING[LEFT]
						FILEDATA[Cursor[Y]-2] += lastData
					}
					Cursor[Y] -= 1
				} else {
					var currentLineData = FILEDATA[Cursor[Y]-1]
					FILEDATA[Cursor[Y]-1] = currentLineData[:Cursor[X]-PADDING[LEFT]-1] + currentLineData[Cursor[X]-PADDING[LEFT]:]
					if Cursor[X] > PADDING[LEFT] {
						Cursor[X] -= 1
					}
				}
				s.Screen.Clear()
			} else {
				var currentLineData = FILEDATA[Cursor[Y]-1]
				FILEDATA[Cursor[Y]-1] = currentLineData[:Cursor[X]-PADDING[LEFT]] + string(event.Rune()) + currentLineData[Cursor[X]-PADDING[LEFT]:]
				Cursor[X] += 1
			}
		}
	}
}
