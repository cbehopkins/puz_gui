package main

import (
	"log"
	"os"

	"github.com/cbehopkins/puz_gui"
	"github.com/icza/gowut/gwu"
)

func main() {

	// Create and build a window
	boggleWin := puz_gui.BoggleWindow()
	sudokuWin := puz_gui.SudokuWindow()
	// Create and start a GUI server (omitting error check)
	server := gwu.NewServer("puzzle", "localhost:8081")
	server.SetLogger(log.New(os.Stdout, "", log.Lshortfile))
	server.SetText("Puzzles")
	server.AddWin(boggleWin)
	server.AddWin(sudokuWin)
	// TBD add countdown

	server.Start("") // Also opens windows list in browser
}
