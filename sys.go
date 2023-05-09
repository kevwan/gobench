package gobench

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func getCpuUsage() float64 {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Println("Failed to get CPU usage")
		return 0
	}

	return percent[0]
}

func getMemoryUsage() float64 {
	memory, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Failed to get memory usage")
		return 0
	}

	return memory.UsedPercent
}
