package main

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// updates are triggered by Bubble Tea cmds, and msg's from those
// cmds are handled here.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.textinput.Value() != "" {
				m.WriteServer(m.textinput.Value())
				m.textinput.SetValue("")
			}
		}
	case ServerMsg:
		m.message = []byte(msg)
		m.viewport.SetContent(string(m.message))
		if useHighPerformanceRenderer {
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
		// move viewport to bottom of content
		cmds = append(cmds, m.getServerMessage)
	case errorMsg:
		_, e := m.conn.Write([]byte("Error: " + msg.Error()))
		if e != nil {
			log.Fatalln(e)
		}
		log.Printf("Error: %s", msg.Error())
		time.Sleep(5 * time.Second)
		if err := m.conn.Close(); err != nil {
			log.Printf("Failed to Close Connection: %s", err.Error())
		}
		return m, tea.Quit
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(string(m.message))
			m.ready = true
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
		if useHighPerformanceRenderer {
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}
	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)
	m.viewport.GotoBottom()
	return m, tea.Batch(cmds...)
}
