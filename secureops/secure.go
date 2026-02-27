package secureops

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gotools/audit"
)

// LockFile cree un fichier .lock pour simuler le verrouillage
func LockFile(filename, outDir string, reader *bufio.Reader) error {
	info, err := os.Stat(filename)
	if err != nil {
		return fmt.Errorf("fichier introuvable: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("%s est un dossier, pas un fichier", filename)
	}

	lockPath := filepath.Join(outDir, filepath.Base(filename)+".lock")

	if _, err := os.Stat(lockPath); err == nil {
		fmt.Printf("  '%s' est deja verrouille.\n", filename)
		return nil
	}

	fmt.Printf("  Verrouiller '%s' ? (yes/no ou oui/non) : ", filename)
	answer, _ := reader.ReadString('\n')
	if !isConfirmed(answer) {
		fmt.Println("  Action annulee.")
		return nil
	}

	f, err := os.Create(lockPath)
	if err != nil {
		return fmt.Errorf("impossible de creer le lock: %w", err)
	}
	f.Close()

	audit.Log(outDir, fmt.Sprintf("LOCK %s", filename))
	fmt.Printf("  '%s' verrouille.\n", filename)
	return nil
}

func UnlockFile(filename, outDir string, reader *bufio.Reader) error {
	lockPath := filepath.Join(outDir, filepath.Base(filename)+".lock")

	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		fmt.Printf("  '%s' n'est pas verrouille.\n", filename)
		return nil
	}

	fmt.Printf("  Deverrouiller '%s' ? (yes/no ou oui/non) : ", filename)
	answer, _ := reader.ReadString('\n')
	if !isConfirmed(answer) {
		fmt.Println("  Action annulee.")
		return nil
	}

	if err := os.Remove(lockPath); err != nil {
		return fmt.Errorf("impossible de supprimer le lock: %w", err)
	}

	audit.Log(outDir, fmt.Sprintf("UNLOCK %s", filename))
	fmt.Printf("  '%s' deverrouille.\n", filename)
	return nil
}

func IsLocked(filename, outDir string) bool {
	lockPath := filepath.Join(outDir, filepath.Base(filename)+".lock")
	_, err := os.Stat(lockPath)
	return err == nil
}

func SetReadOnly(path, outDir string) error {
	if err := os.Chmod(path, 0444); err != nil {
		return fmt.Errorf("chmod impossible: %w", err)
	}
	audit.Log(outDir, fmt.Sprintf("CHMOD read-only %s", path))
	fmt.Printf("  '%s' passe en lecture seule.\n", path)
	return nil
}

func SetReadWrite(path, outDir string) error {
	if err := os.Chmod(path, 0644); err != nil {
		return fmt.Errorf("chmod impossible: %w", err)
	}
	audit.Log(outDir, fmt.Sprintf("CHMOD read-write %s", path))
	fmt.Printf("  '%s' passe en lecture/ecriture.\n", path)
	return nil
}

func CheckPermissions(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("fichier introuvable: %w", err)
	}

	perm := info.Mode().Perm()
	fmt.Printf("  Fichier     : %s\n", path)
	fmt.Printf("  Permissions : %s\n", perm)

	if perm&0200 == 0 {
		fmt.Println("  Attention: fichier en lecture seule")
	}
	if perm&0077 != 0 {
		fmt.Println("  Attention: accessible par d'autres utilisateurs")
	}
	return nil
}

func isConfirmed(s string) bool {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "yes", "y", "oui", "o":
		return true
	default:
		return false
	}
}
