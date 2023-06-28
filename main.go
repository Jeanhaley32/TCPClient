package main

import (
	"flag"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

// model struct, pretty simple model as is.

const (
	useHighPerformanceRenderer = true
)

var (
	// TCPAddr is the address of the TCP Server
	TCPAddr = "localhost:6000"
	// TCPMessage Channel
	ClearScreenMarker = []byte("\033[H\033[2J") // Used to cut out the clear screen message from messages received from server.
)

func init() {
	flag.StringVar(&TCPAddr, "addr", TCPAddr, "address of TCP server. default: localhost:6000")
	flag.Parse()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
