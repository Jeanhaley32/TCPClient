package main

import (
	"fmt"
	"log"
	"net"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// model struct, pretty simple model as is.

type model struct {
	Conn      net.Conn      // Connection Object
	message   ServerMessage // Message taken from Server
	UserMsg   string        // Space for Client's message to server.
	TxtCursor int           // Location of Cursor in Client's message.
}

type TickMsg time.Time

func EveryTwoSeconds() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Message from
type ServerMessage string

func (m *model) ReceiveServerMessage() ServerMessage {
	m.Conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	buffer := make([]byte, 5000)
	_, err := m.Conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	return ServerMessage(buffer)
}

func (m model) Init() tea.Cmd {
	return EveryTwoSeconds()
}

// Updates the model for Bubbletea.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "left":
			if m.TxtCursor > 0 {
				m.TxtCursor--
			}
			return m, nil
		case "right":
			if m.TxtCursor < len(m.UserMsg) {
				m.TxtCursor++
			}
			return m, nil
		case "enter":
			m.Conn.Write([]byte(m.UserMsg))
			m.UserMsg = ""
			m.TxtCursor = 0
			return m, nil
		case "backspace":
			if m.TxtCursor > 0 {
				m.UserMsg = m.UserMsg[:m.TxtCursor-1] + m.UserMsg[m.TxtCursor:]
				m.TxtCursor--
			}
			return m, nil
		default:
			m.UserMsg = m.UserMsg[:m.TxtCursor] + string(msg.Runes) + m.UserMsg[m.TxtCursor:]
			m.TxtCursor++
			return m, nil
		}
	case TickMsg:
		m.ReceiveServerMessage()
		return m, EveryTwoSeconds()
	}
	return m, nil
}

func initialModel() model {
	// Dial into server.
	c, err := net.Dial("tcp", TCPAddr)
	if err != nil {
		log.Fatal(err)
	}
	return model{
		Conn:      c,
		message:   "",
		UserMsg:   "",
		TxtCursor: 0,
	}
}

// Constructs the View for the Bubble Tea program.
func (m model) View() string {
	// Construct Message Prompt for user.
	userPrompt := "\nEnter Message: "
	for i, r := range m.UserMsg {
		if i == m.TxtCursor {
			userPrompt += "_"
		}
		userPrompt += string(r)
	}
	for i, s := range string(m.message) {
		fmt.Printf("%v, %v\n", i, s)
	}
	return string(m.message) + userPrompt
}
