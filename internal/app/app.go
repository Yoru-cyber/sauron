package app

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	viewport viewport.Model
	content  string
	view     string
	width    int
	height   int
	ready    bool
}

type dataUpdateTestMsg struct {
	result string
	err    error
}
type refreshMsg time.Time

func fetchNetworkDataCmd() tea.Cmd {
	return func() tea.Msg {
		result := FetchNetworkInfo()
		if result.Error != nil {
			return dataUpdateTestMsg{result: "", err: result.Error}
		}
		return dataUpdateTestMsg{result: result.Result, err: result.Error}
	}
}
func fetchDefaultDataCmd() tea.Cmd {
	return func() tea.Msg {
		result := FetchDefaultInfo()
		if result.Error != nil {
			return dataUpdateTestMsg{result: "", err: result.Error}
		}
		return dataUpdateTestMsg{result: result.Result, err: result.Error}
	}
}
func fetchDiskCmd() tea.Cmd {
	return func() tea.Msg {
		result := FetchDiskInfo()
		if result.Error != nil {
			return dataUpdateTestMsg{result: "", err: result.Error}
		}
		return dataUpdateTestMsg{result: result.Result, err: result.Error}
	}
}
func Run() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	f, err := tea.LogToFile("sauron-debug.log", "debug")
	if err != nil {
		fmt.Println("couldn't open a file for logging:", err)
		os.Exit(1)
	}
	defer f.Close()
	// Initialize the model
	m := model{view: "default"}

	// Start Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	// Schedule the initial fetch command and start the periodic refresh.
	return tea.Batch(
		tick(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "s":
			log.Println("View(): Pressed 's', changing to default view.")
			m.view = "default"
			m.viewport.SetContent(m.content)
			return m, tea.Batch(cmds...)
		case "n":
			log.Println("View(): Pressed 'n', changing to network view.")
			m.view = "network"
			m.viewport.SetContent(m.content)
			return m, tea.Batch(cmds...)
		case "d":
			log.Println("View(): Pressed 'd', changing to disk view.")
			m.view = "disk"
			m.viewport.SetContent(m.content)
			return m, tea.Batch(cmds...)
		}

	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}
		m.width = msg.Width
		m.height = msg.Height
	case dataUpdateTestMsg:
		if msg.err != nil {
			m.content = fmt.Sprintf("Error: %v\nRetrying in 1 second...", msg.err)
		} else {
			m.content = msg.result
		}
		m.viewport.SetContent(m.content)

		return m, nil
	case refreshMsg:
		// When a tick message is received, trigger a new data fetch
		// based on the current view.
		switch m.view {
		case "default":
			log.Println("Refresh(): Scheduling a refresh default data fetch command.")
			cmds = append(cmds, fetchDefaultDataCmd())
		case "network":
			log.Println("Refresh(): Scheduling a refresh network data fetch command.")
			cmds = append(cmds, fetchNetworkDataCmd())
		case "disk":
			log.Println("Refresh(): Scheduling a refresh disk data fetch command.")
			cmds = append(cmds, fetchDiskCmd())
		}
		cmds = append(cmds, tick())
		return m, tea.Batch(cmds...)
	case os.Signal:
		return m, tea.Quit
	}
	log.Println("Debug(): updating viewport")
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return m.viewport.View()
}
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return refreshMsg(t)
	})
}
