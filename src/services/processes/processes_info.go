package processes

import (
	"fmt"
	"sort"

	"github.com/shirou/gopsutil/v4/process"
)

func GetProcess() {
	processess, _ := process.Processes()
	sort.Slice(processess, func(i, j int) bool {
		p1, _ := processess[i].CPUPercent()
		p2, _ := processess[j].CPUPercent()
		return p1 > p2
	})

	processessRow := ""
	for i := 0; i < 10; i++ {
		n, _ := processess[i].Name()
		cp, _ := processess[i].CPUPercent()
		rowColor := ""
		if i%2 == 0 {
			rowColor = "bg-gray-500"
		}
		processessRow += fmt.Sprintf(`
		<li class="flex justify-between gap-x-4 py-1 rounded-sm %s">
			<span class="mx-2 p-1">%s (PID %d)</span>
			<span class="mx-2 p-1">%.2f%% CPU</span>
		</li>
		`, rowColor, n, processess[i].Pid, cp)
	}
}
