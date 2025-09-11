package app

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yoru-cyber/Sauron/internal/utils"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

const footerHelp = "Press s for default, n for network, d for disk and q to quit"

type SystemData struct {
	HeadInfo string
	RAMInfo  string
	CPUInfo  string
	NetInfo  string
	DiskInfo string
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
	DiskInfo, err := GetDiskInfo()
	if err != nil {
		return nil, err
	}
	return &SystemData{
		HeadInfo: HeadInfo,
		RAMInfo:  RAMInfo,
		CPUInfo:  CPUInfo,
		NetInfo:  NetInfo,
		DiskInfo: DiskInfo,
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

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(DraculaColors.Cyan)).
		PaddingBottom(1)

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(DraculaColors.Foreground))

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(DraculaColors.Comment)).
		Width(20) // Align labels

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(DraculaColors.Green))

	boxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(DraculaColors.Purple)).
		Padding(1, 2)

	content := lipgloss.JoinVertical(lipgloss.Left,
		headerStyle.Render("💻 System Info"),

		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("🖥️  Platform:"),
			infoStyle.Render(hostInfo.Platform),
		),

		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("🔧 OS:"),
			infoStyle.Render(hostInfo.OS),
		),

		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("⚙️  Kernel:"),
			infoStyle.Render(hostInfo.KernelArch),
		),

		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("🏷️  Hostname:"),
			valueStyle.Render(hostInfo.Hostname),
		),

		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("⏱️  Uptime:"),
			valueStyle.Render(utils.FormatUptime(time.Duration(hostInfo.Uptime)*time.Second)),
		),

		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("🔢 CPU Cores:"),
			valueStyle.Render(fmt.Sprintf("%d", cores)),
		),

		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("🧠 Logical Cores:"),
			valueStyle.Render(fmt.Sprintf("%d", logicalCores)),
		),
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
	var output strings.Builder
	output.Grow(1024)
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	counters, err := net.IOCounters(true) // pernic = true
	if err != nil {
		return "", err
	}
	for i, iface := range interfaces {
		if iface.HardwareAddr != "" {
			str := FormatNetworkOutput(iface, counters[i])
			output.WriteString(str)
		}
	}
	return output.String(), nil
}
func GetDiskInfo() (string, error) {
	var output strings.Builder
	output.Grow(1024)
	ioCounters, err := disk.IOCounters()
	if err != nil {
		return "", nil
	}
	for _, ioCounter := range ioCounters {
		output.WriteString("Device: ")
		output.WriteString(ioCounter.Name)
		output.WriteString("↑ Read: ")
		output.WriteString(strconv.FormatUint(ioCounter.ReadCount, 10))
		output.WriteString(" ")
		output.WriteString("↓ Write: ")
		output.WriteString(strconv.FormatUint(ioCounter.WriteCount, 10))
		output.WriteString("\n")
	}
	return output.String(), nil
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
	return sb.String()
}
func DefaultView(data SystemData) string {
	var sb strings.Builder
	sb.Grow(1024)
	sb.WriteString(data.HeadInfo)
	sb.WriteString("\n")
	sb.WriteString(data.RAMInfo)
	sb.WriteString("\n")
	sb.WriteString(data.CPUInfo)
	sb.WriteString("\n")

	sb.WriteString(footerHelp)
	return sb.String()
}
func NetworkView(data SystemData) string {
	var sb strings.Builder
	sb.Grow(1024)
	sb.WriteString(data.NetInfo)
	sb.WriteString("\n")
	sb.WriteString(footerHelp)
	return sb.String()
}
func DiskView(data SystemData) string {
	var sb strings.Builder
	sb.Grow(1024)
	sb.WriteString(data.DiskInfo)
	sb.WriteString("\n")
	return sb.String()
}
