// Define msg, cmds, and styles
package main

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ServerMsg struct {
	topBanner     string
	welcome       string
	connections   string
	help          string
	connID        string
	sessionLength string
	banner        string
	factoid       string
	chats         string
}

type errorMsg error

type TickMsg time.Time

type ServerString string

func (t TickMsg) time() time.Time {
	return t.time()
}

// Splits ServerString into ServerMsg Components.
func (s ServerString) msgSplit() ServerMsg {
	strs := strings.Split(string(s), "\n")
	topBanner := strings.Join(strs[0:6], "")
	welcome := strs[8]
	connections := strs[9]
	help := strs[10]
	connID := strs[11][:15]
	sessionLength := strs[11][19:]
	factoid := strings.Join(strs[13:15], "")
	chats := strings.Join(strs[16:], "")

	ServerMsg := ServerMsg{
		topBanner:     topBanner,
		welcome:       welcome,
		connections:   connections,
		help:          help,
		connID:        connID,
		sessionLength: sessionLength,
		factoid:       factoid,
		chats:         chats,
	}
	return ServerMsg
}

// Cmds
func oneSecondTick() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) getServerMessage() tea.Msg {
	buffer := make([]byte, 4000)
	_, err := m.conn.Read(buffer)
	if err != nil {
		return errorMsg(err)
	}
	if len(strings.Split(string(buffer), "\n")) < 17 {
		return ServerString(buffer).msgSplit()
	}
	return nil
}

func (m model) WriteServer(s string) tea.Msg {
	m.conn.SetWriteDeadline(time.Now().Add(writedeadline))
	_, err := m.conn.Write([]byte(s + "\n"))
	if err != nil {
		return err
	}
	return nil
}

// LipGlows Styles
var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)
