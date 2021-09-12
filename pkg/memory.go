package pkg

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/containerd/cgroups"
	v1 "github.com/containerd/cgroups/stats/v1"
)

func GetMemInfoInByte() (total, free, used, cache float64) {
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
	if stats.Memory.Usage.Limit/1024 > hostMemInfo.MemTotal {
		total = float64(hostMemInfo.MemTotal) * 1024
	} else {
		total = float64(stats.Memory.Usage.Limit)
	}
	used = float64(stats.Memory.Usage.Usage)
	cache = float64(stats.Memory.Cache)
	free = total - used
	return
}

func GetTotalMemInMiB() (total, free, used, cache float64) {
	total, free, used, cache = GetMemInfoInByte()
	return total/1048576, free/1048576, used/1048576, cache/1048576
}
func GetTotalMemInKiB() (total, free, used, cache float64) {
	total, free, used, cache = GetMemInfoInByte()
	return total/1024, free/1024, used/1024, cache/1024
}