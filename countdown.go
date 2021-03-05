package puzgui

import (
	"fmt"
	"log"
	"strconv"

	cntslv "github.com/cbehopkins/countdown/cnt_slv"

	"github.com/icza/gowut/gwu"
)

// RunCountdown give a target number and available numbers
// return the solution as a string
func RunCountdown(target int, sources []int) string {
	foundValues := cntslv.NewNumMap()
	//found_values.SelfTest = true
	foundValues.UseMult = true
	foundValues.PermuteMode = cntslv.LonMap
	foundValues.SeekShort = false // TBD make this controllable

	fmt.Println("Starting permute")
	returnProofs := foundValues.CountHelper(target, sources)
	for range returnProofs {
		//fmt.Println("Proof Received", v)
	}
	//fmt.Println("Permute Complete", proof_list)
	return foundValues.GetProof(target)
}

type countdownProcessHandler struct {
	size   int
	table  gwu.Table
	lab    gwu.Label
	target gwu.TextBox
}

func (h *countdownProcessHandler) HandleEvent(e gwu.Event) {
	if _, isButton := e.Src().(gwu.Button); isButton {
		txt := ""
		success := true
		var ra []int
		ra = make([]int, h.size)

		// Go through the input table and extract data from it
		extractFunc := func(x, y int) bool {
			fmt.Println("Cd access", x, y)
			str, err := gt(x, y, h.table)
			if err != nil {
				txt += "Unable to get data from cell:" + strconv.Itoa(x) + "," + strconv.Itoa(y) + "\n"
				success = false
				return false
			}
			i, err := strconv.Atoi(str)
			if err != nil {
				txt += "\"" + str + "\" is not a valid number\n"
				success = false
				return false
			}
			if y == 0 {
				ra[x] = i
			} else {
				log.Fatal("Countdown should have a 1 dimensional table")
			}
			return false
		}
		tableUVals(h.table, e, extractFunc)
		targetVal, err := strconv.Atoi(h.target.Text())
		if err != nil {
			success = false
			txt += "Target cannot be read" + fmt.Sprint(err)
		}
		if !success {
			txt += "Failure to extract data from table"
		} else {
			txt += RunCountdown(targetVal, ra)

		}
		h.lab.SetText(txt)
		e.MarkDirty(h.lab)
	}
}

type countdownClearHandler struct {
	size   int
	table  gwu.Table
	lab    gwu.Label
	target gwu.TextBox
}

func (h *countdownClearHandler) HandleEvent(e gwu.Event) {
	if _, isButton := e.Src().(gwu.Button); isButton {
		clearFunc := func(x, y int) string {
			return ""
		}
		stTableUVals(h.table, e, clearFunc)
		h.lab.SetText("")
		h.target.SetText("")
		e.MarkDirty(h.lab)
		e.MarkDirty(h.target)
	}
}

// CountdownWindow return gui window for countdown puzzle
func CountdownWindow() gwu.Window {
	size := 6

	win := gwu.NewWindow("cnt", "Countdown")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HACenter)
	win.SetCellPadding(2)

	// A panel for each major thing
	panelButtons := gwu.NewHorizontalPanel()
	panelTable := gwu.NewHorizontalPanel()
	panelTarget := gwu.NewHorizontalPanel()
	table := newInputTableFlex(size, 1)
	panelTable.Add(gwu.NewLabel("Inputs"))
	panelTable.Add(table)

	panelTarget.Add(gwu.NewLabel("Target:"))
	targetTxt := gwu.NewTextBox("")
	panelTarget.Add(targetTxt)
	resultTxt := gwu.NewLabel("")

	buttonProcess := gwu.NewButton("Process")
	buttonProcess.AddEHandler(&countdownProcessHandler{size: size, table: table, target: targetTxt, lab: resultTxt}, gwu.ETypeClick)
	buttonClear := gwu.NewButton("Clear")
	buttonClear.AddEHandler(&countdownClearHandler{size: size, table: table, target: targetTxt, lab: resultTxt}, gwu.ETypeClick)
	panelButtons.Add(buttonProcess)
	panelButtons.Add(buttonClear)

	win.Add(panelButtons)
	win.Add(panelTable)
	win.Add(panelTarget)
	win.Add(resultTxt)
	return win
}
