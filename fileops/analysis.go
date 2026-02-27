package fileops

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func FileInfo(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("fichier introuvable: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("%s est un dossier, pas un fichier", path)
	}

	lines, err := countLines(path)
	if err != nil {
		return err
	}

	fmt.Printf("  Fichier    : %s\n", path)
	fmt.Printf("  Taille     : %d octets\n", info.Size())
	fmt.Printf("  Modifie    : %s\n", info.ModTime().Format("2006-01-02 15:04:05"))
	fmt.Printf("  Nb lignes  : %d\n", lines)
	return nil
}

// WordStats compte les mots en ignorant ceux qui sont purement numeriques
func WordStats(path string) error {
	words, err := extractWords(path)
	if err != nil {
		return err
	}

	totalLen := 0
	for _, w := range words {
		totalLen += len(w)
	}

	avg := 0.0
	if len(words) > 0 {
		avg = float64(totalLen) / float64(len(words))
	}

	fmt.Printf("  Mots (hors numeriques) : %d\n", len(words))
	fmt.Printf("  Longueur moyenne       : %.1f caracteres\n", avg)
	return nil
}

func CountKeyword(path, keyword string) (int, error) {
	lines, err := readLines(path)
	if err != nil {
		return 0, err
	}

	count := 0
	kw := strings.ToLower(keyword)
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), kw) {
			count++
		}
	}

	fmt.Printf("  Lignes contenant \"%s\" : %d\n", keyword, count)
	return count, nil
}

// FilterKeyword separe les lignes qui contiennent le mot-clé et celles qui ne le contiennent pas
func FilterKeyword(path, keyword, outDir string) error {
	lines, err := readLines(path)
	if err != nil {
		return err
	}

	kw := strings.ToLower(keyword)
	var with, without []string
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), kw) {
			with = append(with, line)
		} else {
			without = append(without, line)
		}
	}

	if err := writeLines(filepath.Join(outDir, "filtered.txt"), with); err != nil {
		return err
	}
	fmt.Printf("  -> %d lignes dans %s/filtered.txt\n", len(with), outDir)

	if err := writeLines(filepath.Join(outDir, "filtered_not.txt"), without); err != nil {
		return err
	}
	fmt.Printf("  -> %d lignes dans %s/filtered_not.txt\n", len(without), outDir)
	return nil
}

func Head(path string, n int, outDir string) error {
	if n < 0 {
		n = 0
	}
	lines, err := readLines(path)
	if err != nil {
		return err
	}
	if n > len(lines) {
		n = len(lines)
	}
	dest := filepath.Join(outDir, "head.txt")
	if err := writeLines(dest, lines[:n]); err != nil {
		return err
	}
	fmt.Printf("  -> %d premieres lignes ecrites dans %s\n", n, dest)
	return nil
}

func Tail(path string, n int, outDir string) error {
	if n < 0 {
		n = 0
	}
	lines, err := readLines(path)
	if err != nil {
		return err
	}
	start := len(lines) - n
	if start < 0 {
		start = 0
	}
	written := len(lines[start:])
	dest := filepath.Join(outDir, "tail.txt")
	if err := writeLines(dest, lines[start:]); err != nil {
		return err
	}
	fmt.Printf("  -> %d dernieres lignes ecrites dans %s\n", written, dest)
	return nil
}

// --- helpers ---

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("impossible d'ouvrir %s: %w", path, err)
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines, sc.Err()
}

func countLines(path string) (int, error) {
	lines, err := readLines(path)
	return len(lines), err
}

func extractWords(path string) ([]string, error) {
	lines, err := readLines(path)
	if err != nil {
		return nil, err
	}
	var words []string
	for _, line := range lines {
		for _, w := range strings.Fields(line) {
			if !isNumeric(w) {
				words = append(words, w)
			}
		}
	}
	return words, nil
}

// isNumeric renvoie true si le mot ne contient que des chiffres/ponctuation
func isNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) && r != '.' && r != ',' {
			return false
		}
	}
	return true
}

func writeLines(path string, lines []string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("impossible de creer %s: %w", path, err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
