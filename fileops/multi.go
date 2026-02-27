package fileops

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// BatchAnalyze parcourt tous les .txt d'un dossier et affiche les infos
func BatchAnalyze(dir string) error {
	files, err := FindTxtFiles(dir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Println("  Aucun fichier .txt dans le dossier", dir)
		return nil
	}
	for _, f := range files {
		fmt.Printf("\n--- %s ---\n", f)
		if err := FileInfo(f); err != nil {
			fmt.Println("  Erreur:", err)
			continue
		}
		if err := WordStats(f); err != nil {
			fmt.Println("  Erreur:", err)
		}
	}
	return nil
}

func GenerateReport(dir, outDir string) error {
	files, err := FindTxtFiles(dir)
	if err != nil {
		return err
	}

	out, err := os.Create(filepath.Join(outDir, "report.txt"))
	if err != nil {
		return fmt.Errorf("impossible de creer report.txt: %w", err)
	}
	defer out.Close()

	fmt.Fprintf(out, "=== RAPPORT GLOBAL ===\n")
	fmt.Fprintf(out, "Dossier   : %s\n", dir)
	fmt.Fprintf(out, "Fichiers   : %d\n\n", len(files))

	totalWords, totalLines := 0, 0
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			fmt.Fprintf(out, "Fichier : %s\n  Erreur stat: %v\n\n", f, err)
			continue
		}
		lines, err := readLines(f)
		if err != nil {
			fmt.Fprintf(out, "Fichier : %s\n  Erreur lecture: %v\n\n", f, err)
			continue
		}
		words, err := extractWords(f)
		if err != nil {
			fmt.Fprintf(out, "Fichier : %s\n  Erreur analyse mots: %v\n\n", f, err)
			continue
		}
		totalLines += len(lines)
		totalWords += len(words)

		fmt.Fprintf(out, "Fichier : %s\n", f)
		fmt.Fprintf(out, "  Taille: %d | Lignes: %d | Mots: %d\n\n", info.Size(), len(lines), len(words))
	}

	fmt.Fprintf(out, "--- TOTAUX ---\n")
	fmt.Fprintf(out, "Lignes: %d\n", totalLines)
	fmt.Fprintf(out, "Mots:   %d\n", totalWords)

	fmt.Printf("  -> Rapport ecrit dans %s/report.txt\n", outDir)
	return nil
}

func GenerateIndex(dir, outDir string) error {
	files, err := FindTxtFiles(dir)
	if err != nil {
		return err
	}

	out, err := os.Create(filepath.Join(outDir, "index.txt"))
	if err != nil {
		return fmt.Errorf("impossible de creer index.txt: %w", err)
	}
	defer out.Close()

	fmt.Fprintf(out, "%-40s %10s  %s\n", "CHEMIN", "TAILLE", "DATE MODIFICATION")
	fmt.Fprintf(out, "%s\n", strings.Repeat("-", 75))
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			fmt.Fprintf(out, "%-40s %10s  %s\n", f, "ERR", err.Error())
			continue
		}
		fmt.Fprintf(out, "%-40s %10d  %s\n", f, info.Size(), info.ModTime().Format("2006-01-02 15:04"))
	}

	fmt.Printf("  -> Index ecrit dans %s/index.txt\n", outDir)
	return nil
}

func MergeFiles(dir, outDir string) error {
	files, err := FindTxtFiles(dir)
	if err != nil {
		return err
	}

	out, err := os.Create(filepath.Join(outDir, "merged.txt"))
	if err != nil {
		return fmt.Errorf("impossible de creer merged.txt: %w", err)
	}
	defer out.Close()

	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Erreur lecture %s: %v\n", f, err)
			continue
		}
		fmt.Fprintf(out, "=== %s ===\n", f)
		if _, err := out.Write(data); err != nil {
			fmt.Fprintf(os.Stderr, "  Erreur ecriture %s: %v\n", f, err)
			continue
		}
		fmt.Fprintln(out)
	}

	fmt.Printf("  -> Fusion dans %s/merged.txt\n", outDir)
	return nil
}

func FindTxtFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("impossible de lire %s: %w", dir, err)
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".txt") {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(files)
	return files, nil
}

func ReadLines(path string) ([]string, error)    { return readLines(path) }
func ExtractWords(path string) ([]string, error) { return extractWords(path) }
