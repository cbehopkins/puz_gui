package puz_gui

import (
	"github.com/cbehopkins/boggle"
	"github.com/icza/gowut/gwu"
)

func IsERune(s string) (rune, bool) {
	found_char := false
	var fr rune
	var nv rune
	for _, r := range s {
		if found_char {
			return nv, false
		}
		found_char = true
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return nv, false
		} else {
			fr = r
		}
	}
	return fr, true
}

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
		for x := 0; x < h.size; x++ {
			for y := 0; y < h.size; y++ {
				c := h.table.CompAt(x, y)
				if c == nil {
					success = false
				} else {
					tbox, isTextBox := c.(gwu.TextBox)
					if isTextBox {
						tt := tbox.Text()
						fr, ok := IsERune(tt)
						if ok {
							ra[x][y] = fr
						} else {
							success = false
						}
					} else {
						success = false
					}
				}
			}
		}
		if !success {
			txt = ""
		} else {
			//fmt.Println("Success")
			wrds_found := make(map[string]struct{})
			wrkFunc := func(wrd string) {
				//fmt.Println("Found Word", wrd)
				wrds_found[wrd] = struct{}{}
			}

			pz := h.dic.NewPuzzle(h.size, ra)
			pz.StartWorker(wrkFunc)
			pz.RunWalk()
			pz.Shutdown()
			wrd_cnt := len(wrds_found)
			h.lab.SetRows(wrd_cnt)
			for wrd := range wrds_found {
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
		for x := 0; x < h.size; x++ {
			for y := 0; y < h.size; y++ {
				c := h.table.CompAt(x, y)
				if c == nil {
				} else {
					tbox, isTextBox := c.(gwu.TextBox)
					if isTextBox {
						tbox.SetText("")
						e.MarkDirty(tbox)
					}
				}
			}
		}
		h.lab.SetText("")
		e.MarkDirty(h.lab)
	}
}
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

	//tbArray := make([][]*gwu.TextBox, size)
	table := gwu.NewTable()
	table.EnsureSize(size, size)
	table.SetCellPadding(2)
	newCel := func(x, y int) {
		tb := gwu.NewTextBox("")
		tb.SetMaxLength(1)
		tb.Style().SetWidthPx(10)
		tb.AddSyncOnETypes(gwu.ETypeKeyUp)
		if table.Add(tb, x, y) {
		}
		//tbArray[x][y] = &tb
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			newCel(i, j)
		}
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			table.CellFmt(i, j).Style().SetWidthPx(20)
			table.CompAt(i, j).Style().SetFullSize()
			table.CompAt(i, j).Style().SetFullWidth()
			//table.RowFmt(0).Style().SetBackground(gwu.ClrRed)
		}
	}
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
