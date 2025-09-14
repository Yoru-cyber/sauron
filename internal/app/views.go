package app

import "strings"

func BuildDefaultView(headInfo string, cpuInfo string, ramInfo string) string {
	var sb strings.Builder
	sb.Grow(1024)
	sb.WriteString(headInfo)
	sb.WriteString("\n")
	sb.WriteString(ramInfo)
	sb.WriteString("\n")
	sb.WriteString(cpuInfo)
	sb.WriteString("\n")

	sb.WriteString(footerHelp)
	return sb.String()
}
