package utils

import (
	"sync"

	"github.com/shirou/gopsutil/host"
)

var cachedHostInfo *host.InfoStat
var hostInfoOnce sync.Once
var hostInfoErr error

func GetHostInfo() (*host.InfoStat, error) {
	hostInfoOnce.Do(func() {
		cachedHostInfo, hostInfoErr = host.Info()
	})
	return cachedHostInfo, hostInfoErr
}
