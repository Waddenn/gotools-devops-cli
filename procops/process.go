package procops

import (
	"bufio"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"gotools/audit"
)

type Process struct {
	PID  int
	Name string
}

// ListProcesses recupere les processus via la commande adaptee a l'OS
func ListProcesses(topN int) ([]Process, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("tasklist", "/FO", "CSV", "/NH")
	case "darwin":
		cmd = exec.Command("ps", "-Ao", "pid,comm")
	default: // linux
		cmd = exec.Command("ps", "-Ao", "pid,comm", "--no-headers")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("erreur commande: %w", err)
	}

	procs := parseProcesses(string(output))
	if topN > 0 && topN < len(procs) {
		procs = procs[:topN]
	}
	return procs, nil
}

func SearchProcesses(keyword string, topN int) ([]Process, error) {
	all, err := ListProcesses(0)
	if err != nil {
		return nil, err
	}

	kw := strings.ToLower(keyword)
	var found []Process
	for _, p := range all {
		if strings.Contains(strings.ToLower(p.Name), kw) {
			found = append(found, p)
		}
	}
	if topN > 0 && topN < len(found) {
		found = found[:topN]
	}
	return found, nil
}

// KillProcess demande confirmation avant de tuer un processus
func KillProcess(pid int, outDir string, reader *bufio.Reader) error {
	if pid <= 0 {
		return fmt.Errorf("PID invalide: %d", pid)
	}

	name := findProcessName(pid)
	fmt.Printf("  Processus : PID=%d  Nom=%s\n", pid, name)
	fmt.Print("  Confirmer l'arret ? (yes/no ou oui/non) : ")

	answer, _ := reader.ReadString('\n')
	if !isConfirmed(answer) {
		fmt.Println("  Action annulee.")
		return nil
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/T")
	default:
		cmd = exec.Command("kill", strconv.Itoa(pid))
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("impossible de tuer PID %d: %w", pid, err)
	}

	audit.Log(outDir, fmt.Sprintf("KILL PID=%d (%s)", pid, name))
	fmt.Printf("  Processus %d termine.\n", pid)
	return nil
}

func PrintProcesses(procs []Process) {
	fmt.Printf("  %-8s  %s\n", "PID", "NOM")
	fmt.Printf("  %s\n", strings.Repeat("-", 40))
	for _, p := range procs {
		fmt.Printf("  %-8d  %s\n", p.PID, p.Name)
	}
	fmt.Printf("  Total: %d\n", len(procs))
}

// --- parsing de la sortie selon l'OS ---

func parseProcesses(output string) []Process {
	var procs []Process
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var p *Process
		if runtime.GOOS == "windows" {
			p = parseWindowsLine(line)
		} else {
			p = parseUnixLine(line)
		}
		if p != nil {
			procs = append(procs, *p)
		}
	}
	return procs
}

// format csv windows: "nom.exe","PID",...
func parseWindowsLine(line string) *Process {
	parts := strings.Split(line, ",")
	if len(parts) < 2 {
		return nil
	}
	name := strings.Trim(parts[0], "\"")
	pid, err := strconv.Atoi(strings.Trim(parts[1], "\""))
	if err != nil {
		return nil
	}
	return &Process{PID: pid, Name: name}
}

// format unix: PID COMMAND
func parseUnixLine(line string) *Process {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return nil
	}
	pid, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil
	}
	return &Process{PID: pid, Name: strings.Join(fields[1:], " ")}
}

func findProcessName(pid int) string {
	procs, err := ListProcesses(0)
	if err != nil {
		return "(inconnu)"
	}
	for _, p := range procs {
		if p.PID == pid {
			return p.Name
		}
	}
	return "(inconnu)"
}

func isConfirmed(s string) bool {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "yes", "y", "oui", "o":
		return true
	default:
		return false
	}
}
