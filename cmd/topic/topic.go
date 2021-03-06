package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	linuxproc "github.com/c9s/goprocinfo/linux"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"topic/pkg"
)

const (
	menuPrint = "  PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND"
)

var (
	termWidth   = 0
	termHeight  = 0
	loadMonitor *pkg.LoadMonitor
)

func genSummary(stat *linuxproc.Stat) []string {
	currentTime := time.Now().Local().Format("15:04:05")
	upTime := pkg.GetContainerUpTime(stat)
	tc := pkg.GetTaskCount()
	cpuCount := pkg.GetCpuCount(stat)
	userCpu, systemCpu, idleCpu := pkg.GetCpuUsage(cpuCount)
	memInfo := pkg.GetTotalMemInMiB()
	avail := memInfo.Free + memInfo.Cache
	return []string{
		fmt.Sprintf("topic - %v up %s,  %d users,  load average: %s", currentTime, upTime, pkg.GetUsers(), loadMonitor.GetLoad()),
		fmt.Sprintf("Tasks: [%3d](mod:bold) total, [%3d](mod:bold) running, [%3d](mod:bold) sleeping, [%3d](mod:bold) stopped, [%3d](mod:bold) zombie",
			tc.Total, tc.Running, tc.Sleeping, tc.Stopped, tc.Zombie),
		fmt.Sprintf("%%Cpu(%sc): [%2.1f](mod:bold) us, [%2.1f](mod:bold) sy,  [0.0](mod:bold) ni, [%2.1f](mod:bold) id,  [0.0](mod:bold) wa,  [0.0 hi,](mod:bold)  [0.0](mod:bold) si,  [0.0](mod:bold) st",
			pkg.CpuCountToString(cpuCount), userCpu, systemCpu, idleCpu),
		fmt.Sprintf("MiB Mem : [%7.1f](mod:bold) total, [%7.1f](mod:bold) free, [%7.1f](mod:bold) used, [%7.1f](mod:bold) buff/cache",
			memInfo.Total, memInfo.Free, memInfo.Used, memInfo.Cache),
		fmt.Sprintf("MiB Swap:       [0](mod:bold) total,       [0](mod:bold) free,       [0](mod:bold) used. [%7.1f](mod:bold) avail Mem", avail),
	}
}

func genMenu() string {
	spaceFmt := fmt.Sprintf("%%%vs", termWidth-len(menuPrint)-1)
	paddingSpace := fmt.Sprintf(spaceFmt, " ")
	return fmt.Sprintf("[%s%s](fg:black,bg:white)", menuPrint, paddingSpace)
}

func genProcesses(taskMonitor *pkg.TaskMonitor) string {
	result := []string{
		genMenu(),
	}
	result = append(result, taskMonitor.GetTaskInfos()...)
	return strings.Join(result, "\n")
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Fatal("stat read fail")
	}

	termWidth, termHeight = ui.TerminalDimensions()

	taskMonitor := pkg.NewTaskMonitor(stat)

	loadMonitor = pkg.NewLoadMonitor()
	go loadMonitor.Run()

	summary := widgets.NewList()
	summary.Rows = genSummary(stat)
	summary.TextStyle = ui.NewStyle(ui.ColorWhite)
	summary.WrapText = false
	summary.SetRect(-1, -1, termWidth, termHeight)
	summary.Border = false

	processes := widgets.NewParagraph()
	processes.Text = genProcesses(taskMonitor)
	processes.WrapText = false
	processes.SetRect(-1, 5, termWidth, termHeight)
	processes.Border = false

	draw := func() {
		stat, err := linuxproc.ReadStat("/proc/stat")
		if err != nil {
			log.Fatal("stat read fail")
		}

		termWidth, termHeight = ui.TerminalDimensions()
		summary.Rows = genSummary(stat)
		processes.Text = genProcesses(taskMonitor)
		ui.Render(summary, processes)
	}
	draw()

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker:
			draw()
		}
	}
}
