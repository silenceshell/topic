package pkg

import (
	"fmt"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

var (
	prevUser   int64 = 0
	prevSystem int64 = 0
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
func GetCpuUsage() (user, system, idle float64) {
	var user_, system_ int64
	user_ = getCgoupValueByPath("/sys/fs/cgroup/cpuacct/cpuacct.usage_user")
	system_ = getCgoupValueByPath("/sys/fs/cgroup/cpuacct/cpuacct.usage_sys")

	if prevUser == 0 && prevSystem == 0 {
		prevUser = user_
		prevSystem = system_
		return
	}

	//todo: should divided by count of cpus.
	user = float64(user_-prevUser) / 10000000       // / 1000,000,000 * 100 = /10,000,000
	system = float64(system_-prevSystem) / 10000000 // / 1000,000,000 * 100 = /10,000,000
	idle = 100 - user - system

	prevUser = user_
	prevSystem = system_

	return
}
