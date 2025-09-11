package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/Yoru-cyber/Sauron/internal/utils"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type SystemData struct {
	HeadInfo string
	RAMInfo  string
	CPUInfo  string
	NetInfo  string
}

func FetchAllData() (*SystemData, error) {
	HeadInfo, err := GetHeader()
	if err != nil {
		return nil, err
	}
	RAMInfo, err := GetRamUsage()
	if err != nil {
		return nil, err
	}
	CPUInfo, err := GetCPUInfo()
	if err != nil {
		return nil, err
	}
	NetInfo, err := GetNetwork()
	if err != nil {
		return nil, err
	}
	return &SystemData{
		HeadInfo: HeadInfo,
		RAMInfo:  RAMInfo,
		CPUInfo:  CPUInfo,
		NetInfo:  NetInfo,
	}, nil
}
func GetRamUsage() (string, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return "", err
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
		usedGB, totalGB, bar(barContent), percent), nil
}
func GetHeader() (string, error) {
	hostInfo, err := utils.GetHostInfo()
	if err != nil {
		return "", nil
	}
	cores, err := cpu.Counts(false)
	if err != nil {
		return "", nil
	}
	logicalCores, err := cpu.Counts(true)
	if err != nil {
		return "", nil
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
	return boxStyle.Render(content), nil
}
func GetCPUInfo() (string, error) {

	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return "", err
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
		percent[0], bar(barContent)), nil
}
func GetNetwork() (string, error) {
	var output string
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range interfaces {
		if iface.HardwareAddr != "" {
			output += "Interface: " + iface.Name + "\n" + "MAC: " + iface.HardwareAddr + "\n"
		}
	}
	return output, nil
}
func BuildContent(data SystemData) string {
	var sb strings.Builder
	sb.Grow(1024)
	sb.WriteString(data.HeadInfo)
	sb.WriteString("\n")
	sb.WriteString(data.RAMInfo)
	sb.WriteString("\n")
	sb.WriteString(data.CPUInfo)
	sb.WriteString("\n")
	sb.WriteString(data.NetInfo)
	return sb.String()
}
