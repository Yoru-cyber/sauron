package app

import (
	"fmt"
	"time"

	"github.com/Yoru-cyber/Sauron/internal/utils"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func GetRamUsage() string {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Sprintf("RAM: Error: %v", err)
	}

	usedGB := float64(vm.Used) / 1024 / 1024 / 1024
	totalGB := float64(vm.Total) / 1024 / 1024 / 1024
	percent := vm.UsedPercent

	barWidth := 30
	filled := int(float64(barWidth) * percent / 100)
	if filled > barWidth {
		filled = barWidth
	}

	// Choose color based on usage level
	color := DraculaColors.Green
	if percent > 80 {
		color = DraculaColors.Red // red
	} else if percent > 60 {
		color = DraculaColors.Yellow // yellow
	}

	bar := lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render
	barContent := ""
	for i := 0; i < barWidth; i++ {
		if i < filled {
			barContent += "█"
		} else {
			barContent += "░"
		}
	}

	return fmt.Sprintf("RAM: %.1f/%.1f GB %s %.1f%%",
		usedGB, totalGB, bar(barContent), percent)
}

// showCursor enables the terminal cursor using an ANSI escape code.
func ShowCursor() {
	fmt.Print("\033[?25h")
}

// hideCursor disables the terminal cursor using an ANSI escape code.
func HideCursor() {
	fmt.Print("\033[?25l")
}
func PrintHead() string {
	hostInfo, err := utils.GetHostInfo()
	if err != nil {
		panic(err)
	}
	cores, err := cpu.Counts(false)
	if err != nil {
		panic(err)
	}
	logicalCores, err := cpu.Counts(true)
	if err != nil {
		panic(err)
	}
	boxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2)

	content := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("13")).Render("System Info"),
		fmt.Sprintf("Platform:    %s", hostInfo.Platform),
		fmt.Sprintf("OS:          %s", hostInfo.OS),
		fmt.Sprintf("Kernel:      %s", hostInfo.KernelArch),
		fmt.Sprintf("Hostname:    %s", hostInfo.Hostname),
		fmt.Sprintf("Uptime:      %s", time.Duration(hostInfo.Uptime)*time.Second),
		fmt.Sprintf("CPU cores:    %d", cores),
		fmt.Sprintf("CPU logical cores:    %d", logicalCores),
	)
	return boxStyle.Render(content)
}
func GetCPUInfo() string {

	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		panic(err)
	}
	color := DraculaColors.Green
	if percent[0] > 80 {
		color = DraculaColors.Red
	} else if percent[0] > 60 {
		color = DraculaColors.Yellow
	}
	barWidth := 30
	filled := int(float64(barWidth) * percent[0] / 100)
	bar := lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render
	barContent := ""
	for i := 0; i < barWidth; i++ {
		if i < filled {
			barContent += "█"
		} else {
			barContent += "░"
		}

	}
	return fmt.Sprintf("CPU: %.1f %s",
		percent[0], bar(barContent))
}
