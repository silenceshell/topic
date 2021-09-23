package pkg

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/containerd/cgroups"
	v1 "github.com/containerd/cgroups/stats/v1"
)

type MemInfo struct {
	Total, Free, Used, Cache float64
}

func GetMemInfoInByte() (memInfo MemInfo) {
	stats := v1.Metrics{}
	mem := cgroups.NewMemory("/sys/fs/cgroup/", cgroups.IgnoreModules("memsw"))
	err := mem.Stat("", &stats)
	if err != nil {
		return
	}

	hostMemInfo, err := linuxproc.ReadMemInfo("/proc/meminfo")
	if err != nil {
		return
	}

	// convert Byte to MiB
	var total float64
	if stats.Memory.Usage.Limit/1024 > hostMemInfo.MemTotal {
		total = float64(hostMemInfo.MemTotal) * 1024
	} else {
		total = float64(stats.Memory.Usage.Limit)
	}
	used := float64(stats.Memory.Usage.Usage)
	cache := float64(stats.Memory.Cache)
	free := total - used
	return MemInfo{Total: total, Free: free, Used: used, Cache: cache}
}

func GetTotalMemInMiB() (memInfo MemInfo) {
	memInfo = GetMemInfoInByte()
	memInfo.Total = memInfo.Total / 1048576
	memInfo.Free = memInfo.Free / 1048576
	memInfo.Used = memInfo.Used / 1048576
	memInfo.Cache = memInfo.Cache / 1048576
	return memInfo
}
func GetTotalMemInKiB() (memInfo MemInfo) {
	memInfo = GetMemInfoInByte()
	memInfo.Total = memInfo.Total / 1024
	memInfo.Free = memInfo.Free / 1024
	memInfo.Used = memInfo.Used / 1024
	memInfo.Cache = memInfo.Cache / 1024
	return memInfo
}
