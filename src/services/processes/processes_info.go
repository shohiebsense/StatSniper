package processes

import (
	"StatSniper/models"
	"log"
	"sort"

	"github.com/shirou/gopsutil/v4/process"
)


func GetCpuProcessInfo() []models.CpuProcess {
	processes, err := process.Processes()
	if err != nil || len(processes) == 0 {
	
		log.Println(err)
		return []models.CpuProcess{}
	}

	processCPUData := make([]struct {
		Proc *process.Process
		CPU  float64
	}, len(processes))

	for i, proc := range processes {
		cpuPercent, err := proc.CPUPercent()
		if err != nil {
			cpuPercent = 0 
		}
		processCPUData[i] = struct {
			Proc *process.Process
			CPU  float64
		}{Proc: proc, CPU: cpuPercent}
	}

	sort.Slice(processCPUData, func(i, j int) bool {
		return processCPUData[i].CPU > processCPUData[j].CPU
	})

	limit := 10
	if len(processCPUData) < limit {
		limit = len(processCPUData)
	}

	cpuProcesses := make([]models.CpuProcess, 0, limit)
	for i := 0; i < limit; i++ {
		proc := processCPUData[i].Proc
		cpu := processCPUData[i].CPU
		name, err := proc.Name()
		if err != nil {
			name = "Unknown" 
		}

		cpuProcesses = append(cpuProcesses, models.CpuProcess{
			ProcessName: name,
			CpuUsage:    cpu,
			Pid:         proc.Pid,
		})
	}

	return cpuProcesses
}