package main

import (
	"log"
	"net"
	"os"

	"github.com/cbehopkins/puz_gui"
	"github.com/icza/gowut/gwu"
)

func buildWin(s gwu.Session) {
	// Create and build a window
	boggleWin := puzGui.BoggleWindow()
	sudokuWin := puzGui.SudokuWindow()
	countdownWin := puzGui.CountdownWindow()
	anaWin := puzGui.AnaWindow()
	s.AddWin(boggleWin)
	s.AddWin(sudokuWin)
	s.AddWin(countdownWin)
	s.AddWin(anaWin)
}

// SessHandler is our session handler to build the showcases window.
type sessHandler struct{}

func (h sessHandler) Created(s gwu.Session) {
	buildWin(s)
	win := gwu.NewWindow("show", "Available Puzzles")
	cntLnk := gwu.NewLink("Countdown", "cnt")
	sodLnk := gwu.NewLink("Sudoku", "sod")
	bogLnk := gwu.NewLink("Boggle", "bog")
	anaLnk := gwu.NewLink("Anagram", "ana")

	cntLnk.SetTarget("")
	sodLnk.SetTarget("")
	bogLnk.SetTarget("")
	anaLnk.SetTarget("")

	win.Add(cntLnk)
	win.Add(sodLnk)
	win.Add(bogLnk)
	win.Add(anaLnk)
	s.AddWin(win)
}

func (h sessHandler) Removed(s gwu.Session) {}

// Get preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr, _, _ := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		localAddr = ""
	}

	return localAddr
}

func serverString(local bool) string {
	host_string := "localhost"
	port := "8081"
	if !local {
		ip := GetOutboundIP()
		if ip != "" {
			host_string = ip
		}
	}
	return host_string + ":" + port
}

func main() {
	localServer := false

	// Create and start a GUI server (omitting error check)
	server := gwu.NewServer("", serverString(localServer))
	server.SetLogger(log.New(os.Stdout, "", log.Lshortfile))
	server.SetText("Puzzles")

	if localServer {
		// Create and build a window
		boggleWin := puzGui.BoggleWindow()
		sudokuWin := puzGui.SudokuWindow()
		countdownWin := puzGui.CountdownWindow()
		anaWin := puzGui.AnaWindow()
		//
		server.AddWin(boggleWin)
		server.AddWin(sudokuWin)
		server.AddWin(countdownWin)
		server.AddWin(anaWin)

		server.Start("") // Also opens windows list in browser
	} else {

		server.AddSessCreatorName("show", "Puzzle Creator")
		server.AddSHandler(sessHandler{})
		//autoOpen := false
		// Start GUI server
		//		var openWins []string
		//		if autoOpen {
		//			openWins = []string{"show"}
		//		}
		if err := server.Start("show"); err != nil {
			log.Println("Error: Cound not start GUI server:", err)
			return
		}
	}
}
