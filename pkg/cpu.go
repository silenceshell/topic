package pkg

import (
	"fmt"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

var (
	prevUsageUser   int64 = 0
	prevUsageSystem int64 = 0
)

func GetCpuCount(stat *linuxproc.Stat) (count float64) {
	cfsQuota := getCgoupValueByPath("/sys/fs/cgroup/cpu/cpu.cfs_quota_us")
	cfsPeriod := getCgoupValueByPath("/sys/fs/cgroup/cpu/cpu.cfs_period_us")

	if cfsQuota == -1 {
		return float64(len(stat.CPUStats))
	}

	return float64(cfsQuota / cfsPeriod)
}

func CpuCountToString(c float64) string {
	if c == float64(int64(c)) {
		return fmt.Sprintf("%v", c)
	}
	return fmt.Sprintf("%0.1f", c)
}

// GetCpuUsage should be called every 1 seconds. not quite precise.
func GetCpuUsage(cpus float64) (user, system, idle float64) {
	var currentUsageUser, currentUsageSystem int64
	currentUsageUser = getCgoupValueByPath("/sys/fs/cgroup/cpuacct/cpuacct.usage_user")
	currentUsageSystem = getCgoupValueByPath("/sys/fs/cgroup/cpuacct/cpuacct.usage_sys")

	if prevUsageUser == 0 && prevUsageSystem == 0 {
		prevUsageUser = currentUsageUser
		prevUsageSystem = currentUsageSystem
		return
	}

	user = float64(currentUsageUser-prevUsageUser) / 10000000 / cpus       // / 1000,000,000 * 100 = /10,000,000
	system = float64(currentUsageSystem-prevUsageSystem) / 10000000 / cpus // / 1000,000,000 * 100 = /10,000,000
	idle = 100 - user - system
	if idle < 0 {
		idle = 0
	}

	prevUsageUser = currentUsageUser
	prevUsageSystem = currentUsageSystem

	return
}
