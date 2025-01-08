package models

type CpuProcess struct {
	ProcessName string `json:"processName"`
	CpuUsage    float64    `json:"cpuUsage"`
	Pid 		int32 `json:"pid"`
}