package puzgui

import (
	"fmt"

	"github.com/cbehopkins/ana"
	"github.com/cbehopkins/wordlist"

	"github.com/icza/gowut/gwu"
)

type anaProcessHandler struct {
	input  gwu.TextBox
	output gwu.Label
}

func (h *anaProcessHandler) HandleEvent(e gwu.Event) {
	if _, isButton := e.Src().(gwu.Button); isButton {
		defer e.MarkDirty(h.output)
		refString := h.input.Text()
		data, err := wordlist.Asset("data/wordlist.txt")
		if err != nil {
			h.output.SetText(fmt.Sprintln("Asset wordlist not found:", err))
			return
		}

		h.output.SetText(ana.AnagramWord(refString, data).String())
	}
}

type anaClearHandler struct {
	input  gwu.TextBox
	output gwu.Label
}

func (h *anaClearHandler) HandleEvent(e gwu.Event) {
	if _, isButton := e.Src().(gwu.Button); isButton {
		h.input.SetText("")
		h.output.SetText("")
		e.MarkDirty(h.input)
		e.MarkDirty(h.output)
	}
}

// AnaWindow creates the window object that
// all the anagram if resides in
func AnaWindow() gwu.Window {

	win := gwu.NewWindow("ana", "Anagram")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HACenter)
	win.SetCellPadding(2)

	// A panel for each major thing
	panelTb := gwu.NewHorizontalPanel()
	panelButtons := gwu.NewHorizontalPanel()
	srcTxt := gwu.NewTextBox("")
	panelTb.Add(srcTxt)

	resultTxt := gwu.NewLabel("")

	buttonProcess := gwu.NewButton("Process")
	buttonProcess.AddEHandler(&anaProcessHandler{input: srcTxt, output: resultTxt}, gwu.ETypeClick)
	buttonClear := gwu.NewButton("Clear")
	buttonClear.AddEHandler(&anaClearHandler{input: srcTxt, output: resultTxt}, gwu.ETypeClick)
	panelButtons.Add(buttonProcess)
	panelButtons.Add(buttonClear)

	win.Add(panelButtons)
	win.Add(panelTb)
	win.Add(resultTxt)
	return win
}
