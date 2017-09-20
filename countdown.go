package puz_gui

import (
	"fmt"
	"log"
	"strconv"

	"github.com/cbehopkins/countdown/cnt_slv"

	"github.com/icza/gowut/gwu"
)

func RunCountdown(target int, sources []int) string {
	if false {
		found_values := cntSlv.NewNumMap()
		//found_values.SelfTest = true
		found_values.UseMult = true
		found_values.PermuteMode = cntSlv.FastMap
		found_values.SeekShort = false // TBD make this controllable

		fmt.Println("Starting permute")
		return_proofs := found_values.CountHelper(target, sources)
		for _ = range return_proofs {
			//fmt.Println("Proof Received", v)
		}
		//fmt.Println("Permute Complete", proof_list)
		return found_values.GetProof(target)
	} else {
		findShortest := false
		return cntSlv.CountFastHelper(target, sources, findShortest)
	}
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
			str, err := Gt(x, y, h.table)
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
		TableUVals(h.table, e, extractFunc)
		targetVal, err := strconv.Atoi(h.target.Text())
		if err != nil {
			success = false
			txt += "Target cannot be read" + fmt.Sprint(err)
		}
		if !success {
			txt += "Failure to extract data from table"
		} else {
			//txt = "Result of countdown:\n"
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
		StTableUVals(h.table, e, clearFunc)
		h.lab.SetText("")
		h.target.SetText("")
		e.MarkDirty(h.lab)
		e.MarkDirty(h.target)
	}
}
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
