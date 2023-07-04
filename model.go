package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	readdeadline   = 500 * time.Millisecond
	writedeadline  = 500 * time.Millisecond
	characterLimit = 124
)

// model struct represents the state of the application.
type model struct {
	textinput textinput.Model
	conn      net.Conn // Connection Object
	message   string   // Message taken from Server
	viewport  viewport.Model
	ready     bool
}

// Runs initial cmd.
func (m model) Init() tea.Cmd {
	return m.getServerMessage
}

// func (m model) headerView() string {
// 	title := titleStyle.Render("TCP Client")
// 	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
// 	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
// }

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

// Set model's initial state.
func initialModel() model {
	// Dial into server.
	c, err := net.Dial("tcp", TCPAddr)
	if err != nil {
		log.Fatal(err)
	}
	ti := textinput.New()
	ti.CharLimit = characterLimit
	ti.Width = characterLimit + 1
	ti.Cursor.Blink = true
	ti.Placeholder = "Enter Message"
	ti.Focus()
	model := model{
		conn:      c,
		textinput: ti,
		message:   "Retrieving Server Message",
	}
	// Return model with initial state.
	return model
}

// Constructs the View for the Bubble Tea program.
func (m model) View() string {
	return fmt.Sprintf("%v\n%v\n%v", m.viewport.View(), m.footerView(), m.textinput.View())
}
