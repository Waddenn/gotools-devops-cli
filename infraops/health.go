package infraops

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const (
	red   = "\033[31m"
	green = "\033[32m"
	reset = "\033[0m"
)

// CheckDiskSpace verifie l'espace disque et affiche une alerte si < 10% libre
func CheckDiskSpace() error {
	var used float64
	var err error

	if runtime.GOOS == "windows" {
		used, err = diskSpaceWindows()
	} else {
		used, err = diskSpaceUnix()
	}
	if err != nil {
		return fmt.Errorf("impossible de verifier le disque: %w", err)
	}

	free := 100 - used
	fmt.Printf("  Espace utilise : %.1f%%\n", used)
	fmt.Printf("  Espace libre   : %.1f%%\n", free)

	if free < 10 {
		fmt.Printf("\n  %sALERTE: Espace disque critique (%.1f%% libre) !%s\n", red, free, reset)
	} else {
		fmt.Printf("\n  %sEspace disque: etat normal%s\n", green, reset)
	}
	return nil
}

func diskSpaceUnix() (float64, error) {
	output, err := exec.Command("df", "/").Output()
	if err != nil {
		return 0, err
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("sortie df inattendue")
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 5 {
		return 0, fmt.Errorf("sortie df inattendue")
	}
	pct := strings.TrimSuffix(fields[4], "%")
	return strconv.ParseFloat(pct, 64)
}

func diskSpaceWindows() (float64, error) {
	output, err := exec.Command("wmic", "logicaldisk", "where", "DeviceID='C:'",
		"get", "FreeSpace,Size", "/format:csv").Output()
	if err != nil {
		return 0, err
	}
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		parts := strings.Split(strings.TrimSpace(line), ",")
		if len(parts) >= 3 {
			free, e1 := strconv.ParseFloat(parts[1], 64)
			total, e2 := strconv.ParseFloat(parts[2], 64)
			if e1 == nil && e2 == nil && total > 0 {
				return (1 - free/total) * 100, nil
			}
		}
	}
	return 0, fmt.Errorf("impossible de parser wmic")
}
