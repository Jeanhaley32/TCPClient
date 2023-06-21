package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	// TCPAddr is the address of the TCP Server
	TCPAddr = "localhost:6000"
	// TCPMessage Channel
	ServerString string
	SendMessage  = make(chan string)
	closech      = make(chan interface{})
	ReadRefresh  = 2 * time.Second // Time Between Server Reads.
)

// model struct, pretty simple model as is.
type model struct {
	message   string
	UserMsg   string
	TxtCursor int
}

func (m model) Init() tea.Cmd {
	return nil
}

// Updates the model for Bubbletea.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			closech <- "" // Close the TCP Connection
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
		case "enter":
			SendMessage <- m.UserMsg
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
	}
	return m, nil
}

func initialModel() model {
	return model{
		message:   "",
		UserMsg:   "",
		TxtCursor: 0,
	}
}

// Constructs the View for the Bubble Tea program.
func (m model) View() string {
	// Construct Message Prompt for user.
	view := ServerString + "\nEnter Message: "
	for i, r := range m.UserMsg {
		if i == m.TxtCursor {
			view += "_"
		}
		view += string(r)
	}

	return view
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() error {
		TeaStart()
		wg.Done()
		return nil
	}()
	go func() error {
		err := TCPDialer()
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
		return nil
	}()
	wg.Wait()
	fmt.Println("Program Finished")
}

// Starts the Tea Program
func TeaStart() error {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

// Dials into the TCP Server, and starts the Connection Reader.
func TCPDialer() error {
	c, err := net.Dial("tcp", "localhost:6000")
	if err != nil {
		return err
	}
	defer c.Close()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		err := ConnReader(c)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()
	go func() {
		err := ConnWriter(c)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("TCP Dialer Finished")
	return nil
}

// Reads from TCP Connection, and sends to TCPMessage Channel
// for the Tea Program to read from, and update the model.
func ConnReader(c net.Conn) error {
	defer func() {
		fmt.Println("connection Reader Closed")
	}()
	buffer := make([]byte, 1024)
	for {
		// Read from the connection
		_, err := c.Read(buffer)
		if err != nil {
			return err
		}
		ServerString = string(buffer)
	}
}

// Connection writer for the TCP connection.
// Reads from msg channel, and writes that message to the TCP connection.
func ConnWriter(c net.Conn) error {
	defer func() {
		fmt.Println("connection Writer Closed")
	}()
	for {
		select {
		case msg := <-SendMessage:
			_, err := c.Write([]byte(msg))
			if err != nil {
				return err
			}
		}
	}
}
