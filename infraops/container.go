package infraops

import (
	"fmt"
	"os/exec"
	"strings"
)

type ContainerInfo struct {
	ID     string
	Name   string
	Image  string
	Status string
}

func ListContainers() ([]ContainerInfo, error) {
	cmd := exec.Command("docker", "ps", "--format", "{{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("erreur docker ps (docker lance ?): %w", err)
	}

	var containers []ContainerInfo
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 4)
		if len(parts) < 4 {
			continue
		}
		containers = append(containers, ContainerInfo{
			ID: parts[0], Name: parts[1], Image: parts[2], Status: parts[3],
		})
	}
	return containers, nil
}

func ContainerStats(nameOrID string) error {
	cmd := exec.Command("docker", "stats", "--no-stream", "--format",
		"table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}", nameOrID)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("erreur docker stats: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

func PrintContainers(containers []ContainerInfo) {
	if len(containers) == 0 {
		fmt.Println("  Aucun conteneur actif.")
		return
	}
	fmt.Printf("  %-12s %-20s %-25s %s\n", "ID", "NOM", "IMAGE", "STATUT")
	fmt.Printf("  %s\n", strings.Repeat("-", 75))
	for _, c := range containers {
		fmt.Printf("  %-12s %-20s %-25s %s\n", c.ID, c.Name, c.Image, c.Status)
	}
}
