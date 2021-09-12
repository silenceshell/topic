package pkg

import (
	"fmt"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"os"
	"strconv"
	"time"
)

func GetContainerUpTime(stat *linuxproc.Stat) string {
	btime := stat.BootTime

	process1, _ := linuxproc.ReadProcessStat("/proc/1/stat")
	startTime := process1.Starttime

	uptime := btime.Add((time.Duration(startTime) * 10) * time.Millisecond)

	timeSince := time.Since(uptime)
	if timeSince > time.Hour*24 {
		days := timeSince / time.Hour / 24
		hours := timeSince / time.Hour % 24
		minutes := timeSince / time.Minute % 60
		return fmt.Sprintf("%d days, %d:%d", days, hours, minutes)
	}
	return uptime.Local().Format("15:04")
}

func GetTaskCount() (total, running, sleeping, stopped, zombie int) {
	f, err := os.Open("/proc")
	if err != nil {
		fmt.Println(err)
		return
	}
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range files {
		fileName := v.Name()
		if fileName[0] < '0' || fileName[0] > '9' {
			continue
		}
		pid, err := strconv.Atoi(v.Name())
		if err == nil {
			p, err := linuxproc.ReadProcessStat(fmt.Sprintf("/proc/%d/stat", pid))
			if err != nil {
				return
			}
			switch p.State {
			case "R":
				running++
			case "t", "T":
				stopped++
			case "Z":
				zombie++
			default:
				sleeping++
			}
			total++
		}
	}
	return
}

const (
	taskInfoFmt = "%5d %-8s%4d%4d%8d %6d %6d %s  %4.1f  %4.1f  %8s %-s"
)

func convertDuration(t time.Duration) string {
	minute := t / time.Minute
	t = t % time.Minute
	second := t / time.Second
	t = t % time.Second
	milliseconds := t / time.Millisecond / 10
	return fmt.Sprintf("%d:%02d.%02d", minute, second, milliseconds)
}

func (t *TaskMonitor) GetTaskInfos() (infos []string) {
	f, err := os.Open("/proc")
	if err != nil {
		fmt.Println(err)
		return
	}
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return
	}

	total, _, _, _ := GetTotalMemInKiB()

	taskInfos := make([]string, 0)
	for _, v := range files {
		fileName := v.Name()
		if fileName[0] < '0' || fileName[0] > '9' {
			continue
		}
		pid, err := strconv.Atoi(v.Name())
		if err == nil {
			procInfo, err := linuxproc.ReadProcess(uint64(pid), "/proc")
			if err != nil {
				fmt.Println(err)
				return
			}

			// in KiB
			virt := procInfo.Stat.Vsize / 1024
			res := procInfo.Statm.Resident * 4
			shr := procInfo.Statm.Share * 4

			var cpuUsage uint64
			if t.taskPrevUser[pid] != 0 || t.taskPrevSystem[pid] != 0 {
				cpuUsage = procInfo.Stat.Utime+procInfo.Stat.Stime - t.taskPrevUser[pid] - t.taskPrevSystem[pid]
			}
			t.taskPrevUser[pid] = procInfo.Stat.Utime
			t.taskPrevSystem[pid] = procInfo.Stat.Stime

			memUsage := float64(res) / total * 100
			uptime := (time.Duration(procInfo.Stat.Utime+procInfo.Stat.Stime) * 10) * time.Millisecond

			taskInfo := fmt.Sprintf(taskInfoFmt, pid, "root", procInfo.Stat.Priority, procInfo.Stat.Nice,
				virt, res, shr, procInfo.Stat.State, float64(cpuUsage), memUsage,
				convertDuration(uptime), procInfo.Cmdline)
			taskInfos = append(taskInfos, taskInfo)
		}
	}
	return taskInfos
}

type TaskMonitor struct {
	stat *linuxproc.Stat
	taskPrevUser map[int]uint64
	taskPrevSystem map[int]uint64
}

func NewTaskMonitor(stat *linuxproc.Stat) *TaskMonitor {
	taskMonitor := TaskMonitor{
		stat: stat,
		taskPrevUser: make(map[int]uint64),
		taskPrevSystem: make(map[int]uint64),
	}
	return &taskMonitor
}