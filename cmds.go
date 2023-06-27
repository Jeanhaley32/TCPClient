// Define msg, cmds, and styles
package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errorMsg error

type TickMsg time.Time

type ServerMsg string

func (t TickMsg) time() time.Time {
	return t.time()
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

	return ServerMsg(string(buffer))
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
