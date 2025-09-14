package app

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Yoru-cyber/Sauron/internal/constants"
	"github.com/Yoru-cyber/Sauron/internal/utils"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

const footerHelp = "Press s for default, n for network, d for disk and q to quit"

// cache for Default info as it is OS, CPU, Kernel info so it won't change
var (
	headerInfoCache Result
	cacheMutex      sync.Mutex
)

type SystemData struct {
	HeadInfo string
	RAMInfo  string
	CPUInfo  string
	NetInfo  string
	DiskInfo string
}
type Result struct {
	Result string
	Error  error
}
type ResultChan chan Result

func FetchNetworkInfo() Result {
	nChan := make(ResultChan)
	go GetNetwork(nChan)
	nResult := <-nChan
	return nResult
}
func FetchDiskInfo() Result {
	dChan := make(ResultChan)
	go GetDiskInfo(dChan)
	dResult := <-dChan
	return dResult
}
func getOrCacheHeader() Result {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Check if the cache is already populated.
	if headerInfoCache.Result != "" {
		return headerInfoCache
	}

	// If the cache is empty, perform the expensive data fetch.
	hChan := make(ResultChan)
	go GetHeader(hChan) // No need for a goroutine, as this is a blocking call within this function.
	hResult := <-hChan

	// Cache the result for future calls.
	headerInfoCache = hResult
	return hResult
}
func FetchDefaultInfo() Result {
	var hResult Result
	if headerInfoCache.Result == "" {
		hResult = getOrCacheHeader()
	} else {
		hResult = headerInfoCache
	}
	var sb strings.Builder
	sb.Grow(1024)
	var wg sync.WaitGroup
	wg.Add(5)
	var errs []error
	rChan := make(ResultChan)
	cChan := make(ResultChan)

	go func() {
		defer wg.Done()
		GetRamUsage(rChan)
	}()
	go func() {
		defer wg.Done()
		GetCPUInfo(cChan)
	}()

	rResult := <-rChan
	cResult := <-cChan

	if rResult.Error != nil {
		errs = append(errs, fmt.Errorf("RAM fetch failed: %w", rResult.Error))
	}
	if cResult.Error != nil {
		errs = append(errs, fmt.Errorf("CPU fetch failed: %w", cResult.Error))
	}
	content := BuildDefaultView(hResult.Result, cResult.Result, rResult.Result)
	var finalResult Result
	if len(errs) > 0 {
		finalResult = Result{"", errors.Join(errs...)}
	} else {
		finalResult = Result{content, nil}
	}
	return finalResult
}
func FetchAllData() (*SystemData, error) {
	var wg sync.WaitGroup
	var errs []error
	hChan := make(ResultChan)
	rChan := make(ResultChan)
	cChan := make(ResultChan)
	nChan := make(ResultChan)
	dChan := make(ResultChan)
	wg.Add(5)
	go func() {
		defer wg.Done()
		GetHeader(hChan)
	}()
	go func() {
		defer wg.Done()
		GetRamUsage(rChan)
	}()
	go func() {
		defer wg.Done()
		GetCPUInfo(cChan)
	}()
	go func() {
		defer wg.Done()
		GetNetwork(nChan)
	}()
	go func() {
		defer wg.Done()
		GetDiskInfo(dChan)
	}()
	hResult := <-hChan
	rResult := <-rChan
	cResult := <-cChan
	nResult := <-nChan
	dResult := <-dChan
	if hResult.Error != nil {
		errs = append(errs, fmt.Errorf("header fetch failed: %w", hResult.Error))
	}
	if rResult.Error != nil {
		errs = append(errs, fmt.Errorf("RAM fetch failed: %w", rResult.Error))
	}
	if cResult.Error != nil {
		errs = append(errs, fmt.Errorf("CPU fetch failed: %w", cResult.Error))
	}
	if nResult.Error != nil {
		errs = append(errs, fmt.Errorf("network fetch failed: %w", nResult.Error))
	}
	if dResult.Error != nil {
		errs = append(errs, fmt.Errorf("disk fetch failed: %w", dResult.Error))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return &SystemData{
		HeadInfo: hResult.Result,
		RAMInfo:  rResult.Result,
		CPUInfo:  cResult.Result,
		NetInfo:  nResult.Result,
		DiskInfo: dResult.Result,
	}, nil
}
func GetRamUsage(ch ResultChan) {
	var sb strings.Builder
	vm, err := mem.VirtualMemory()
	if err != nil {
		ch <- Result{"", err}
		return
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
	color := constants.DraculaColors.Green
	if percent > 80 {
		color = constants.DraculaColors.Red // red
	} else if percent > 60 {
		color = constants.DraculaColors.Yellow // yellow
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
	sb.WriteString("RAM: ")
	sb.WriteString(fmt.Sprintf("%.1f", usedGB))
	sb.WriteString("/")
	sb.WriteString(fmt.Sprintf("%.1f", totalGB))
	sb.WriteString(" GB ")
	sb.WriteString(bar(barContent))
	sb.WriteString(fmt.Sprintf(" %.1f", percent))
	sb.WriteString("%")
	ch <- Result{sb.String(), nil}
}
func GetHeader(ch ResultChan) {
	hostInfo, err := utils.GetHostInfo()
	if err != nil {
		ch <- Result{"", err}
		return
	}
	cores, err := cpu.Counts(false)
	if err != nil {
		ch <- Result{"", err}
		return
	}
	logicalCores, err := cpu.Counts(true)
	if err != nil {
		ch <- Result{"", err}
		return
	}
	// Should be a different function
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(constants.DraculaColors.Cyan)).
		PaddingBottom(1)

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(constants.DraculaColors.Foreground))

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(constants.DraculaColors.Comment)).
		Width(20) // Align labels

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(constants.DraculaColors.Green))

	boxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(constants.DraculaColors.Purple)).
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
	ch <- Result{boxStyle.Render(content), nil}
}
func GetCPUInfo(ch ResultChan) {
	var sb strings.Builder
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		ch <- Result{"", err}
		return
	}
	pStr := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffff")).Width(16).Render("CPU: " + strconv.FormatFloat(percent[0], 'f', 2, 64))
	color := constants.DraculaColors.Green
	if percent[0] > 80 {
		color = constants.DraculaColors.Red
	} else if percent[0] > 60 {
		color = constants.DraculaColors.Yellow
	}
	barWidth := 30
	filled := int(float64(barWidth) * percent[0] / 100)
	bar := lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render
	barContent := ""
	for i := range barWidth {
		if i < filled {
			barContent += "█"
		} else {
			barContent += "░"
		}

	}
	sb.Grow(1024)
	sb.WriteString(pStr)
	sb.WriteString(bar(barContent))
	ch <- Result{sb.String(), nil}
}
func GetNetwork(ch ResultChan) {
	var sb strings.Builder
	sb.Grow(1024)
	interfaces, err := net.Interfaces()
	if err != nil {
		ch <- Result{"", err}
		return
	}
	counters, err := net.IOCounters(true) // pernic = true
	if err != nil {
		ch <- Result{"", err}
		return
	}
	for i, iface := range interfaces {
		if iface.HardwareAddr != "" {
			str := utils.FormatNetworkOutput(iface, counters[i])
			sb.WriteString(str)
		}
	}
	ch <- Result{sb.String(), nil}
}
func GetDiskInfo(ch ResultChan) {
	var sb strings.Builder
	sb.Grow(1024)
	ioCounters, err := disk.IOCounters()
	if err != nil {
		ch <- Result{"", nil}
		return
	}
	for _, ioCounter := range ioCounters {
		sb.WriteString("Device: ")
		sb.WriteString(ioCounter.Name)
		sb.WriteString("↑ Read: ")
		sb.WriteString(strconv.FormatUint(ioCounter.ReadCount, 10))
		sb.WriteString(" ")
		sb.WriteString("↓ Write: ")
		sb.WriteString(strconv.FormatUint(ioCounter.WriteCount, 10))
		sb.WriteString("\n")
	}
	ch <- Result{sb.String(), nil}
}
