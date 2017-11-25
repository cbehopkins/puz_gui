package puzGui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cbehopkins/ana"
	"github.com/icza/gowut/gwu"
)

type anaProcessHandler struct {
	input  gwu.TextBox
	output gwu.Label
}
type result string
type results []result

func (a results) Len() int           { return len(a) }
func (a results) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a results) Less(i, j int) bool { return len(a[i]) < len(a[j]) }
func (a results) String() string {
	ret_txt := ""
	nl := ""
	for _, v := range a {
		ret_txt += string(nl) + string(v)
		nl = "\n"
	}
	return ret_txt
}
func (h *anaProcessHandler) HandleEvent(e gwu.Event) {
	if _, isButton := e.Src().(gwu.Button); isButton {
		filename := "../ana/wordlist.txt"
		refString := h.input.Text()
		refString = strings.Replace(refString, "\n", "", -1)
		fmt.Println("Received String", refString)
		resultChan := ana.Helper(filename, refString, 4)
		results := make(results, 0)
		for res := range resultChan {
			fmt.Println("Received Result", res)
			results = append(results, result(res))
		}

		sort.Sort(sort.Reverse(results))
		h.output.SetText(results.String())
		e.MarkDirty(h.output)
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
