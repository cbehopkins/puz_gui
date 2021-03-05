package puzgui

import (
	"log"
	"strconv"

	"github.com/cbehopkins/sod"
	"github.com/icza/gowut/gwu"
)

func runSudoku(input [][]int) (output [][]int, testPuzzle *sod.Puzzle) {
	size := len(input)
	output = make([][]int, size)
	for i, arr := range input {

		if len(arr) != size {
			log.Fatal("Size error in input")
		}
		output[i] = make([]int, size)
	}

	testPuzzle = sod.NewPuzzle()

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			val := input[y][x]
			if val != 0 {
				tc := sod.Coord{x, y}
				testPuzzle.SetVal(sod.Value(val), tc)
			}
		}
	}
	result := testPuzzle.SelfCheck()
	if result != nil {
		// TBD add error field we can report this to
		return
	}
	testPuzzle.SolveAll()

	result = testPuzzle.SelfCheck()
	if result != nil {
		// TBD return error
		return
	}
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
		success := true
		var intArray [][]int
		intArray = make([][]int, h.size)
		for x := 0; x < h.size; x++ {
			intArray[x] = make([]int, h.size)
		}
		for x := 0; x < h.size; x++ {
			for y := 0; y < h.size; y++ {
				tt, err := gt(x, y, h.table)

				if err != nil {
					success = false
				} else {
					if tt == "" {
						// This is okay, blank input cells are fine
					} else {
						val, ok := isCNum(tt, h.size)
						if ok {
							intArray[y][x] = val // need to invert here annoyingly
						} else {
							success = false
						}
					}
				}
			}
		}
		if !success {
			log.Println("Error! Invalid input")
		} else {
			resultInt, resultPuz := runSudoku(intArray)
			size := h.size
			funcResult := func(x, y int) string {
				var txt string
				val := resultInt[y][x]
				if (val > 0) && (val <= size) {
					txt = strconv.Itoa(val)
				}
				return txt
			}
			funcPartial := func(x, y int) string {
				vals := resultPuz.GetCel(sod.Coord{x, y}).Values()
				txt := ""
				for _, val := range vals {
					txt += strconv.Itoa(int(val))
				}
				return txt
			}

			stTableVals(h.size, h.resultTable, e, funcResult)
			stTableVals(h.size, h.partialTable, e, funcPartial)
		}
	}
}

type sudokuClearHandler struct {
	size        int
	table       gwu.Table
	resultTable gwu.Table
}

func (h *sudokuClearHandler) HandleEvent(e gwu.Event) {
	if _, isButton := e.Src().(gwu.Button); isButton {
		fc := func(x, y int) string {
			return ""
		}
		stTableVals(h.size, h.resultTable, e, fc)
	}
}

// SudokuWindow Return gui window element
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

	tableInput := newInputTable(size, size)

	panelInput.Add(tableInput)

	tableResult := newOutputTable(size, size)
	tablePartial := newOutputTable(size, size)

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
