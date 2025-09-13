package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yoru-cyber/Sauron/internal/constants"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/net"
)

func FormatBytes(bytes uint64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	size := float64(bytes)
	unitIndex := 0

	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}

	return fmt.Sprintf("%.1f%s", size, units[unitIndex])
}
func FormatUptime(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

var (
	interfaceStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(constants.DraculaColors.Purple)). // #bd93f9
			Bold(true)

	macStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(constants.DraculaColors.Comment)). // #6272a4
			Italic(true)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(constants.DraculaColors.Cyan)).
			Width(16)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(constants.DraculaColors.Foreground))

	sentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(constants.DraculaColors.Green)).
			Bold(true)

	receivedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(constants.DraculaColors.Pink)).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(constants.DraculaColors.Red)).
			Bold(true)

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(constants.DraculaColors.Comment)).
			SetString(" • ")
)

func FormatNetworkOutput(iface net.InterfaceStat, counters net.IOCountersStat) string {
	var output strings.Builder

	// Interface name with style
	output.WriteString(interfaceStyle.Render("🖧 " + iface.Name))
	output.WriteString("\n")
	// IP address
	output.WriteString(labelStyle.Render("IP:"))
	output.WriteString(" ")
	output.WriteString(macStyle.Render(iface.Addrs[0].Addr))
	output.WriteString("\n")
	// MAC address
	output.WriteString(labelStyle.Render("MAC:"))
	output.WriteString(" ")
	output.WriteString(macStyle.Render(iface.HardwareAddr))
	output.WriteString("\n")

	// Packets with colors and icons
	output.WriteString(labelStyle.Render("Packets:"))
	output.WriteString(" ")
	output.WriteString(sentStyle.Render("↑" + strconv.FormatUint(counters.PacketsSent, 10)))
	output.WriteString(separatorStyle.String())
	output.WriteString(receivedStyle.Render("↓" + strconv.FormatUint(counters.PacketsRecv, 10)))
	output.WriteString("\n")

	// Bytes with human-readable formatting
	output.WriteString(labelStyle.Render("Data:"))
	output.WriteString(" ")
	output.WriteString(sentStyle.Render("↑" + FormatBytes(counters.BytesSent)))
	output.WriteString(separatorStyle.String())
	output.WriteString(receivedStyle.Render("↓" + FormatBytes(counters.BytesRecv)))
	output.WriteString("\n")

	// Errors (if any) - only show if there are errors
	if counters.Errin > 0 || counters.Errout > 0 {
		output.WriteString(labelStyle.Render("Errors:"))
		output.WriteString(" ")
		output.WriteString(errorStyle.Render(
			"IN:" + strconv.FormatUint(counters.Errin, 10) +
				" OUT:" + strconv.FormatUint(counters.Errout, 10),
		))
		output.WriteString("\n")
	}

	// MTU
	output.WriteString(labelStyle.Render("MTU:"))
	output.WriteString(" ")
	output.WriteString(valueStyle.Render(strconv.Itoa(iface.MTU)))
	output.WriteString("\n")

	return output.String()
}
