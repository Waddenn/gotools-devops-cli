package infraops

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	red   = "\033[31m"
	green = "\033[32m"
	reset = "\033[0m"
)

// CheckDiskSpace verifie l'espace disque et affiche une alerte si < 10% libre.
func CheckDiskSpace() error {
	used, err := diskSpaceUsedPercent()
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

func diskSpaceUsedPercent() (float64, error) {
	if isWindows() {
		return diskSpaceWindows()
	}
	return diskSpaceUnix()
}

func diskSpaceUnix() (float64, error) {
	output, err := exec.Command("df", "/").Output()
	if err != nil {
		return 0, err
	}
	return parseUnixDFUsedPercent(string(output))
}

func diskSpaceWindows() (float64, error) {
	// 1) Compat legacy
	output, err := exec.Command("wmic", "logicaldisk", "where", "DeviceID='C:'", "get", "FreeSpace,Size", "/format:csv").Output()
	if err == nil {
		if used, parseErr := parseWMICUsedPercent(string(output)); parseErr == nil {
			return used, nil
		}
	}

	// 2) Fallback moderne (PowerShell)
	psCmd := "Get-CimInstance Win32_LogicalDisk -Filter \"DeviceID='C:'\" | Select-Object FreeSpace,Size"
	output, err = exec.Command("powershell", "-NoProfile", "-Command", psCmd).Output()
	if err != nil {
		return 0, err
	}
	return parsePowerShellUsedPercent(string(output))
}

func parseUnixDFUsedPercent(output string) (float64, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
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

func parseWMICUsedPercent(output string) (float64, error) {
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "Node") || strings.Contains(line, "FreeSpace") {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) < 3 {
			continue
		}
		free, e1 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		total, e2 := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
		if e1 == nil && e2 == nil && total > 0 {
			return (1 - free/total) * 100, nil
		}
	}
	return 0, fmt.Errorf("impossible de parser wmic")
}

func parsePowerShellUsedPercent(output string) (float64, error) {
	var free, size float64
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "FreeSpace") || strings.HasPrefix(line, "----") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		f, e1 := strconv.ParseFloat(fields[0], 64)
		s, e2 := strconv.ParseFloat(fields[1], 64)
		if e1 == nil && e2 == nil && s > 0 {
			free, size = f, s
			break
		}
	}
	if size <= 0 {
		return 0, fmt.Errorf("impossible de parser powershell")
	}
	return (1 - free/size) * 100, nil
}

func isWindows() bool {
	return os.PathSeparator == '\\'
}
