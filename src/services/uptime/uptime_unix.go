//go:build linux || darwin
// +build linux darwin

package uptime

import (
	"StatSniper/models"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)


func parseSecUsec(input string) (int64, error) {
	startSec := strings.Index(input, "sec =") + len("sec =")
	endSec := strings.Index(input, ",")
	secStr := strings.TrimSpace(input[startSec:endSec])

	startUsec := strings.Index(input, "usec =") + len("usec =")
	endUsec := strings.Index(input, "}")
	usecStr := strings.TrimSpace(input[startUsec:endUsec])

	sec, err := strconv.ParseInt(secStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing sec: %v", err)
	}

	usec, err := strconv.ParseInt(usecStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing usec: %v", err)
	}

	totalTimeNano := sec*1e9 + usec*1e3
	return totalTimeNano, nil
}

func GetSystemUptime() models.Uptime {
	var uptimeSeconds int64

	if runtime.GOOS == "linux" {
		data, err := os.ReadFile("/proc/uptime")
		if err != nil {
			log.Println(fmt.Errorf("error reading /proc/uptime: %v", err))
			return models.Uptime{}
		}

		parts := strings.Fields(string(data))
		if len(parts) < 1 {
			log.Println( fmt.Errorf("error: invalid uptime data"))
			return models.Uptime{}
		}

		uptimeSeconds, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			log.Println(fmt.Errorf("error converting uptime to float: %v", err))
			return models.Uptime{}
		}
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("sysctl", "-n", "kern.boottime")
		output, err := cmd.Output()
		if err != nil {
			log.Println(fmt.Errorf("error running sysctl: %v", err))
			return models.Uptime{}
		}

		outputStr := string(output)
		uptimeSeconds, err = parseSecUsec(outputStr)

		if err != nil {
			log.Println(err)
		}

	} else {
		log.Println( fmt.Errorf("unsupported OS: %s", runtime.GOOS))
		return models.Uptime{}
	}

	currentTimeNano := time.Now().UnixNano()

	uptimeNano := currentTimeNano - uptimeSeconds
	uptimeDuration := time.Duration(uptimeNano)


	days := int(uptimeDuration.Hours()) / 24
	hours := int(uptimeDuration.Hours()) % 24
	minutes := int(uptimeDuration.Minutes()) % 60
	seconds := int(uptimeDuration.Seconds()) % 60

	return models.Uptime{
		Days:    fmt.Sprintf("%d", days),
		Hours:   fmt.Sprintf("%d", hours),
		Minutes: fmt.Sprintf("%d", minutes),
		Seconds: fmt.Sprintf("%d", seconds),
	}
}
