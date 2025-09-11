package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Yoru-cyber/Sauron/internal/app"
	"github.com/Yoru-cyber/Sauron/internal/utils"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/mem"
)

type model struct {
	viewport viewport.Model
	memInfo  *mem.VirtualMemoryStat
	content  string
	data     app.SystemData
	width    int
	height   int
	ready    bool
}
type dataUpdateMsg app.SystemData

func fetchDataCmd() tea.Cmd {
	return func() tea.Msg {
		return dataUpdateMsg(app.FetchAllData())
	}
}

type tickMsg time.Time

func main() {
	if err := utils.InitLogger(); err != nil {
		fmt.Printf("Failed to setup logging: %v\n", err)
		os.Exit(1)
	}
	defer utils.CleanupLogger()
	// Initialize the model
	m := model{}

	// Start Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tick(),
		fetchDataCmd(),
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

	case dataUpdateMsg:
		m.data = app.SystemData(msg)
		m.content = fmt.Sprintf("%s\n%s\n%s\n%s",
			m.data.HeadInfo,
			m.data.RAMInfo,
			m.data.CPUInfo,
			m.data.NetInfo)
		m.viewport.SetContent(m.content)
		cmds = append(cmds, tick())
	case tickMsg:
		cmds = append(cmds, fetchDataCmd())
	case os.Signal:
		return m, tea.Quit
	}
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
		return tickMsg(t)
	})
}
