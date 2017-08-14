package puz_gui

import (
	"fmt"
	"log"
	"strconv"

	"github.com/cbehopkins/sod"
	"github.com/icza/gowut/gwu"
)

func St(txt string, x, y int, table gwu.Table, e gwu.Event) {
	c := table.CompAt(x, y)
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
func runSudoku(input [][]int) (output [][]int, testPuzzle *sod.Puzzle) {
	size := len(input)
	output = make([][]int, size)
	for i, arr := range input {

		if len(arr) != size {
			fmt.Println("Size error in input")
		}
		output[i] = make([]int, size)
	}

	testPuzzle = sod.NewPuzzle()

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			val := input[y][x]
			if val != 0 {
				tc := sod.Coord{x, y}
				fmt.Println("Set Values", val, tc)
				testPuzzle.SetVal(sod.Value(val), tc)
			}
		}
	}
	result := testPuzzle.SelfCheck()
	if result != nil {
		fmt.Println("Self check fail", result)
		return
	}
	log.Println("Print before solve", testPuzzle)
	testPuzzle.SolveAll()
	log.Println("Print after solve", testPuzzle)

	result = testPuzzle.SelfCheck()
	if result != nil {
		fmt.Println("Self check fail after solve", result)
		return
	} else {
		for x := 0; x < size; x++ {
			for y := 0; y < size; y++ {
				vals := testPuzzle.GetCel(sod.Coord{x, y}).Values()
				if len(vals) == 1 {
					val := int(vals[0])
					if val != 0 {
						output[y][x] = val

					}
				}
			}
		}
	}

	return
}

type sudokuProcessHandler struct {
	size         int
	table        gwu.Table
	resultTable  gwu.Table
	partialTable gwu.Table
}

func (h *sudokuProcessHandler) HandleEvent(e gwu.Event) {
	if _, isButton := e.Src().(gwu.Button); isButton {
		log.Println("Handle Sudoku Button pressed")

		success := true
		var intArray [][]int
		intArray = make([][]int, h.size)
		for x := 0; x < h.size; x++ {
			intArray[x] = make([]int, h.size)
		}
		for x := 0; x < h.size; x++ {
			for y := 0; y < h.size; y++ {
				c := h.table.CompAt(y, x)
				if c == nil {
					success = false
				} else {
					tbox, isTextBox := c.(gwu.TextBox)
					if isTextBox {
						tt := tbox.Text()
						if tt == "" {
							// This is okay, blank input cells are fine
						} else {
							val, ok := IsCNum(tt, h.size)
							if ok {
								intArray[y][x] = val // need to invert here annoyingly
							} else {
								success = false
							}
						}
					} else {
						success = false
					}
				}
			}
		}
		if !success {
			fmt.Println("Error!")
		} else {
			//log.Println("Running solver")
			resultInt, resultPuz := runSudoku(intArray)
			//log.Println("Solver complete")
			fmt.Println("Success", resultInt)
			size := len(intArray)
			for x := 0; x < size; x++ {
				for y := 0; y < size; y++ {
					val := resultInt[y][x]
					if (val > 0) && (val <= size) {
						txt := strconv.Itoa(val)
						//fmt.Printf("Set Result:%v x:%v y:%v\n", txt, x, y)
						St(txt, y, x, h.resultTable, e)
					}
				}
			}
			for x := 0; x < size; x++ {
				for y := 0; y < size; y++ {
					vals := resultPuz.GetCel(sod.Coord{x, y}).Values()
					txt := ""
					for _, val := range vals {
						txt += strconv.Itoa(int(val))
						//fmt.Printf("Set Result:%v x:%v y:%v\n", txt, x, y)
					}
					St(txt, y, x, h.partialTable, e)
				}
			}
		}
		log.Println("Handle Sudoku Button complete")

		//h.resultTable.SetText(txt)
		//e.MarkDirty(h.resultTable)
	}
}

type sudokuClearHandler struct {
	size        int
	table       gwu.Table
	resultTable gwu.Table
}

func (h *sudokuClearHandler) HandleEvent(e gwu.Event) {
	if _, isButton := e.Src().(gwu.Button); isButton {
		for x := 0; x < h.size; x++ {
			for y := 0; y < h.size; y++ {
				St("", x, y, h.resultTable, e)
			}
		}
	}
}
func newInputCel(x, y int, tab gwu.Table) {
	tmpTextBox := gwu.NewTextBox("")
	tmpTextBox.SetMaxLength(1)
	tmpTextBox.Style().SetWidthPx(10)
	tmpTextBox.AddSyncOnETypes(gwu.ETypeKeyUp)
	tab.Add(tmpTextBox, x, y)
}
func newOutputCel(x, y int, tab gwu.Table) {
	tmpLabel := gwu.NewLabel("")
	tmpLabel.Style().SetWidthPx(10)
	tmpLabel.AddSyncOnETypes(gwu.ETypeKeyUp)
	//tmpLabel.Style().SetBorderBottom(gwu.BrdStyleSolid)
	//tmpLabel.Style().SetBorderBottom2(5, gwu.BrdStyleSolid, gwu.ClrGrey)
	tab.Add(tmpLabel, x, y)
}
func newLabelTable(size int) gwu.Table {
	tableResult := gwu.NewTable()
	tableResult.EnsureSize(size, size)
	tableResult.SetCellPadding(2)

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			newOutputCel(i, j, tableResult)
		}
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			tableResult.CellFmt(i, j).Style().SetHeightPx(20)
			tableResult.CellFmt(i, j).Style().SetWidthPx(20)
			tableResult.CompAt(i, j).Style().SetFullSize()
			tableResult.CompAt(i, j).Style().SetFullWidth()
			//if j > 0 {
			//	tableResult.CompAt(i, j).Style().SetBorderLeft2(1, gwu.BrdStyleSolid, gwu.ClrGrey)
			//}
			//tableResult.CompAt(i, j).Style().SetBorderTop2(1, gwu.BrdStyleSolid, gwu.ClrGrey)
			//tableResult.CompAt(i, j).Style().SetBorder(gwu.BrdStyleSolid)

		}
	}
	tableResult.Style().SetBorderTop(gwu.BrdStyleSolid)
	tableResult.Style().SetBorderBottom(gwu.BrdStyleSolid)
	tableResult.Style().SetBorderLeft(gwu.BrdStyleSolid)
	tableResult.Style().SetBorderRight(gwu.BrdStyleSolid)
	return tableResult
}
func SudokuWindow() gwu.Window {

	size := 9
	//maj := math.Sqrt(size)

	win := gwu.NewWindow("sod", "Sudoku")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HACenter)
	win.SetCellPadding(2)

	// A panel for each major thing
	panelInput := gwu.NewHorizontalPanel()
	panelButtons := gwu.NewHorizontalPanel()
	panelResults := gwu.NewHorizontalPanel()
	panelDebug := gwu.NewHorizontalPanel()
	//tbArray := make([][]*gwu.TextBox, size)
	tableInput := gwu.NewTable()
	tableInput.EnsureSize(size, size)
	tableInput.SetCellPadding(2)

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			newInputCel(i, j, tableInput)
		}
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			tableInput.CellFmt(i, j).Style().SetWidthPx(20)
			tableInput.CompAt(i, j).Style().SetFullSize()
			tableInput.CompAt(i, j).Style().SetFullWidth()
			//table.RowFmt(0).Style().SetBackground(gwu.ClrRed)
		}
	}
	panelInput.Add(tableInput)

	tableResult := newLabelTable(size)
	tablePartial := newLabelTable(size)

	panelResults.Add(tableResult)
	panelDebug.Add(tablePartial)

	buttonProcess := gwu.NewButton("Process")
	buttonProcess.AddEHandler(&sudokuProcessHandler{size: size, table: tableInput, resultTable: tableResult, partialTable: tablePartial}, gwu.ETypeClick)
	buttonClear := gwu.NewButton("Clear")
	buttonClear.AddEHandler(&sudokuClearHandler{size: size, table: tableInput, resultTable: tableResult}, gwu.ETypeClick)
	panelButtons.Add(buttonProcess)
	panelButtons.Add(buttonClear)

	win.Add(panelButtons)
	win.Add(panelInput)
	win.Add(panelResults)
	win.Add(panelDebug)
	return win
}
