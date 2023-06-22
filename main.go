package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	// TCPAddr is the address of the TCP Server
	TCPAddr = "localhost:6000"
	// TCPMessage Channel
	ServerString string
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
