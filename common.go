package puz_gui

import (
	"errors"
	"strconv"

	"github.com/icza/gowut/gwu"
)

func St(txt string, x, y int, table gwu.Table, e gwu.Event) {
	c := table.CompAt(y, x)
	if c == nil {
	} else {
		tbox, isTextBox := c.(gwu.TextBox)
		if isTextBox {
			tbox.SetText(txt)
			e.MarkDirty(tbox)
		}
		label, isLabel := c.(gwu.Label)
		if isLabel {
			label.SetText(txt)
			e.MarkDirty(label)
		}
	}
}

func Gt(x, y int, table gwu.Table) (string, error) {
	c := table.CompAt(y, x)
	if c == nil {
		return "", errors.New("Nil component")
	} else {
		tbox, isTextBox := c.(gwu.TextBox)
		if isTextBox {
			tt := tbox.Text()
			return tt, nil
		} else {
			return "", errors.New("Unknown box type")
		}
	}
}
func GtRune(x, y int, table gwu.Table, rt runeType) (rune, error) {
	var tr rune
	var ok bool
	tt, err := Gt(x, y, table)
	if err != nil {
		return tr, err
	} else {
		tr, ok = IsERune(tt, rt)
		if ok {
			return tr, nil
		} else {
			return tr, errors.New("Not a rune")
		}
	}
}

func linkTb(prev, current gwu.TextBox) {
	prev.AddEHandlerFunc(func(e gwu.Event) {
		e.SetFocusedComp(current) // Pass the text box component you want to focus
	}, gwu.ETypeKeyUp)
}
func linkTable(table gwu.Table) {
	var prevTb gwu.TextBox
	fun := func(x, y int) bool {
		c := table.CompAt(y, x)
		if c == nil {
		} else {
			tb, isTextBox := c.(gwu.TextBox)
			if isTextBox {
				if prevTb != nil {
					linkTb(prevTb, tb)
				}
				prevTb = tb
			}
		}
		return true // mark dirty
	}
	TableUVals(table, nil, fun)
}

func TableUVals(table gwu.Table, e gwu.Event, fc func(x, y int) bool) {
	run := true
	for y := 0; run; y++ {
		c := table.CompAt(y, 0)
		if c == nil {
			run = false
		} else {
			var runInner bool
			runInner = true
			for x := 0; runInner; x++ {
				c := table.CompAt(y, x)
				if c == nil {
					runInner = false
				} else {
					update := fc(x, y)
					if update && e != nil {
						e.MarkDirty(c)
					}
				}
			}
		}
	}

}
func StTableUVals(table gwu.Table, e gwu.Event, set_func func(int, int) string) {
	fun := func(x, y int) bool {
		txt := set_func(x, y)
		St(txt, x, y, table, e)
		return false
	}
	TableUVals(table, e, fun)
}
func TableVals(size int, table gwu.Table, e gwu.Event, fc func(x, y int)) {
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			fc(x, y)
		}
	}
}
func StTableVals(size int, table gwu.Table, e gwu.Event, set_func func(int, int) string) {
	fun := func(x, y int) {
		txt := set_func(x, y)
		St(txt, x, y, table, e)
	}
	TableVals(size, table, e, fun)
}

func IsCNum(str string, size int) (val int, success bool) {
	i, err := strconv.Atoi(str)
	if err == nil {
		if (i > 0) && (i <= size) {
			return i, true
		} else {
			return
		}
	} else {
		return
	}

}

type runeType int

const (
	otherRune = 1 << iota
	alphaRune
	numRune
)

func checkRune(r rune, rt runeType) bool {
	otherMask := (rt & otherRune) > 0
	alphaMask := (rt & alphaRune) > 0
	numMask := (rt & numRune) > 0
	if otherMask {
		// Any rune is allowed
		return true
	}
	if numMask && r > '0' && r < '9' {
		return true
	}
	if alphaMask && ((r > 'a' && r < 'z') || (r > 'A' && r < 'Z')) {
		return true
	}
	return false
}

func IsERune(s string, rt runeType) (rune, bool) {
	found_char := false
	var fr rune
	var nv rune
	for _, r := range s {
		if found_char {
			return nv, false
		}
		found_char = true
		if checkRune(r, rt) {
			fr = r
		} else {
			return nv, false
		}
	}
	return fr, true
}

func newInputCel(x, y int, tab gwu.Table) {
	tmpTextBox := gwu.NewTextBox("")
	tmpTextBox.SetMaxLength(1)
	tmpTextBox.Style().SetWidthPx(10)
	tmpTextBox.AddSyncOnETypes(gwu.ETypeKeyUp)
	tab.Add(tmpTextBox, y, x)
}
func newOutputCel(x, y int, tab gwu.Table) {
	tmpLabel := gwu.NewLabel("_")
	tmpLabel.Style().SetWidthPx(10)
	tmpLabel.AddSyncOnETypes(gwu.ETypeKeyUp)
	tab.Add(tmpLabel, y, x)
}

func newInputTable(sizeX, sizeY int) gwu.Table {
	tableInput := gwu.NewTable()
	tableInput.EnsureSize(sizeX, sizeY)
	tableInput.SetCellPadding(2)

	for i := 0; i < sizeX; i++ {
		for j := 0; j < sizeY; j++ {
			newInputCel(i, j, tableInput)
		}
	}
	for i := 0; i < sizeX; i++ {
		for j := 0; j < sizeY; j++ {
			s := tableInput.CellFmt(j, i).Style()
			s.SetWidthPx(20)
			s.SetFullSize()
			s.SetFullWidth()
		}
	}
	// Auto move on
	linkTable(tableInput)
	return tableInput
}
func newInputCelFlex(x, y int, tab gwu.Table) {
	tmpTextBox := gwu.NewTextBox("")
	//tmpTextBox.SetMaxLength(1)
	tmpTextBox.Style().SetWidthPx(40)
	//tmpTextBox.AddSyncOnETypes(gwu.ETypeKeyUp)
	tab.Add(tmpTextBox, y, x)
}
func newInputTableFlex(sizeX, sizeY int) gwu.Table {
	tableInput := gwu.NewTable()
	tableInput.EnsureSize(sizeX, sizeY)
	tableInput.SetCellPadding(2)

	for i := 0; i < sizeX; i++ {
		for j := 0; j < sizeY; j++ {
			newInputCelFlex(i, j, tableInput)
		}
	}
	for i := 0; i < sizeX; i++ {
		for j := 0; j < sizeY; j++ {
			s := tableInput.CellFmt(j, i).Style()
			s.SetWidthPx(20)
			s.SetFullSize()
			s.SetFullWidth()
		}
	}
	return tableInput
}
func newOutputTable(sizeX, sizeY int) gwu.Table {
	tableResult := gwu.NewTable()
	tableResult.EnsureSize(sizeX, sizeY)

	for i := 0; i < sizeX; i++ {
		for j := 0; j < sizeY; j++ {
			newOutputCel(j, i, tableResult)
		}
	}
	for i := 0; i < sizeX; i++ {
		for j := 0; j < sizeY; j++ {
			st := tableResult.CompAt(j, i).Style()
			st.SetBorderRight(gwu.BrdStyleSolid)
			st.SetBorderLeft(gwu.BrdStyleSolid)
			st.SetBorderTop(gwu.BrdStyleSolid)
			st.SetBorderBottom(gwu.BrdStyleSolid)
		}
	}
	return tableResult
}
