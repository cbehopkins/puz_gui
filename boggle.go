package puzgui

import (
	"github.com/cbehopkins/boggle"
	"github.com/icza/gowut/gwu"
)

type boggleProcessHandler struct {
	size  int
	table gwu.Table
	lab   gwu.TextBox
	dic   *boggle.DictMap
}

func (h *boggleProcessHandler) HandleEvent(e gwu.Event) {
	if _, isButton := e.Src().(gwu.Button); isButton {
		txt := ""
		success := true
		var ra [][]rune
		ra = make([][]rune, h.size)
		for x := 0; x < h.size; x++ {
			ra[x] = make([]rune, h.size)
		}
		// Go through the input table and extract data from it
		extractFunc := func(x, y int) {
			fr, err := gtRune(x, y, h.table, alphaRune)
			if err != nil {
				success = false
			} else {
				ra[x][y] = fr
			}
		}
		tableVals(h.size, h.table, e, extractFunc)

		if !success {
			txt = ""
		} else {
			wrdsFound := make(map[string]struct{})
			wrkFunc := func(wrd string) {
				wrdsFound[wrd] = struct{}{}
			}

			pz := h.dic.NewPuzzle(h.size, ra)
			pz.StartWorker(wrkFunc)
			pz.RunWalk()
			pz.Shutdown()
			wrdCnt := len(wrdsFound)
			h.lab.SetRows(wrdCnt)
			for wrd := range wrdsFound {
				txt += wrd + "\n"
			}
		}
		h.lab.SetText(txt)
		e.MarkDirty(h.lab)
	}
}

type boggleClearHandler struct {
	size  int
	table gwu.Table
	lab   gwu.TextBox
}

func (h *boggleClearHandler) HandleEvent(e gwu.Event) {
	if _, isButton := e.Src().(gwu.Button); isButton {
		clearFunc := func(x, y int) string {
			return ""
		}
		stTableVals(h.size, h.table, e, clearFunc)
		h.lab.SetText("")
		e.MarkDirty(h.lab)
	}
}

// BoggleWindow Return a window for the boggle gui
func BoggleWindow() gwu.Window {
	size := 4

	dic := boggle.NewDictMap([]string{})
	dic.PopulateFile("../boggle/wordlist.txt")

	win := gwu.NewWindow("bog", "Boggle")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HACenter)
	win.SetCellPadding(2)

	// A panel for each major thing
	panelTable := gwu.NewHorizontalPanel()
	panelButtons := gwu.NewHorizontalPanel()

	table := newInputTable(size, size)
	panelTable.Add(table)

	resultTxt := gwu.NewTextBox("")

	buttonProcess := gwu.NewButton("Process")
	buttonProcess.AddEHandler(&boggleProcessHandler{size: size, table: table, lab: resultTxt, dic: dic}, gwu.ETypeClick)
	buttonClear := gwu.NewButton("Clear")
	buttonClear.AddEHandler(&boggleClearHandler{size: size, table: table, lab: resultTxt}, gwu.ETypeClick)
	panelButtons.Add(buttonProcess)
	panelButtons.Add(buttonClear)

	win.Add(panelButtons)
	win.Add(panelTable)
	win.Add(resultTxt)
	return win
}
