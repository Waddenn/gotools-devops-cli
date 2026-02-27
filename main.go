package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"gotools/config"
	"gotools/fileops"
	"gotools/infraops"
	"gotools/procops"
	"gotools/secureops"
	"gotools/webops"
)

var (
	cfg    *config.Config
	reader *bufio.Reader
)

func main() {
	configPath := flag.String("config", "", "chemin vers config.txt ou config.json")
	flag.Parse()

	var err error
	cfg, err = loadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur config: %v\n", err)
		fmt.Println("Config par defaut chargee.")
		cfg = config.DefaultConfig()
	}

	if err := cfg.EnsureOutDir(); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur creation dossier out: %v\n", err)
		os.Exit(1)
	}

	reader = bufio.NewReader(os.Stdin)

	// boucle principale
	for {
		printMenu()
		choice := readLine("Choix : ")

		switch strings.ToUpper(choice) {
		case "A":
			menuAnalysis()
		case "B":
			menuMultiFiles()
		case "C":
			menuWiki()
		case "D":
			menuProcessOps()
		case "E":
			menuSecureOps()
		case "F":
			menuContainerOps()
		case "G":
			menuHealthCheck()
		case "H":
			menuParallelScan()
		case "Q":
			fmt.Println("Au revoir !")
			return
		default:
			fmt.Println("Choix invalide.")
		}
		fmt.Println()
	}
}

func loadConfig(path string) (*config.Config, error) {
	if path != "" {
		return config.Load(path)
	}
	// on tente json d'abord, sinon txt
	if _, err := os.Stat("config.json"); err == nil {
		return config.Load("config.json")
	}
	if _, err := os.Stat("config.txt"); err == nil {
		return config.Load("config.txt")
	}
	return nil, fmt.Errorf("aucun fichier de config trouve")
}

func printMenu() {
	fmt.Println("\n========================================")
	fmt.Println("          GOTOOLS - Menu principal")
	fmt.Println("========================================")
	fmt.Println("  [A] FileOps  - Analyse d'un fichier")
	fmt.Println("  [B] FileOps  - Analyse d'un dossier (.txt)")
	fmt.Println("  [C] WebOps   - Wikipedia")
	fmt.Println("  [D] ProcOps  - Gestion processus")
	fmt.Println("  [E] SecureOps- Securite / permissions")
	fmt.Println("  [F] InfraOps - Docker")
	fmt.Println("  [G] InfraOps - Etat disque")
	fmt.Println("  [H] InfraOps - Scan parallele (.txt)")
	fmt.Println("  [Q] Quitter")
	fmt.Println("========================================")
}

// ---- Choix A ----

func menuAnalysis() {
	path := readLineDefault("Fichier a analyser", cfg.DefaultFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Erreur: '%s' introuvable.\n", path)
		return
	}

	fmt.Println("\n--- Infos fichier ---")
	if err := fileops.FileInfo(path); err != nil {
		fmt.Println("Erreur:", err)
		return
	}

	fmt.Println("\n--- Stats mots ---")
	if err := fileops.WordStats(path); err != nil {
		fmt.Println("Erreur:", err)
		return
	}

	keyword := readLine("Mot-cle pour filtrage (optionnel) : ")
	if keyword != "" {
		fmt.Println("\n--- Comptage ---")
		if _, err := fileops.CountKeyword(path, keyword); err != nil {
			fmt.Println("Erreur:", err)
		}

		fmt.Println("\n--- Filtrage ---")
		if err := fileops.FilterKeyword(path, keyword, cfg.OutDir); err != nil {
			fmt.Println("Erreur:", err)
		}
	}

	n := readIntMin("Nombre de lignes pour head/tail", 5, 0)
	fmt.Println("\n--- Head ---")
	if err := fileops.Head(path, n, cfg.OutDir); err != nil {
		fmt.Println("Erreur:", err)
	}
	fmt.Println("\n--- Tail ---")
	if err := fileops.Tail(path, n, cfg.OutDir); err != nil {
		fmt.Println("Erreur:", err)
	}
}

// ---- Choix B ----

func menuMultiFiles() {
	dir := readLineDefault("Dossier a analyser", cfg.BaseDir)
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		fmt.Printf("Erreur: '%s' n'est pas un dossier valide.\n", dir)
		return
	}

	runStep("Batch", func() error { return fileops.BatchAnalyze(dir) })
	runStep("Rapport", func() error { return fileops.GenerateReport(dir, cfg.OutDir) })
	runStep("Index", func() error { return fileops.GenerateIndex(dir, cfg.OutDir) })
	runStep("Fusion", func() error { return fileops.MergeFiles(dir, cfg.OutDir) })
}

// ---- Choix C ----

func menuWiki() {
	input := readLine("Article(s) Wikipedia (virgule pour separer, ex: Go_(langage)) : ")
	if input == "" {
		fmt.Println("Aucun article saisi.")
		return
	}

	// on split par virgule si plusieurs articles
	var articles []string
	for _, a := range strings.Split(input, ",") {
		a = strings.TrimSpace(a)
		if a != "" {
			articles = append(articles, a)
		}
	}

	if len(articles) == 1 {
		if err := webops.AnalyzeArticle(articles[0], cfg.WikiLang, cfg.OutDir); err != nil {
			fmt.Println("Erreur:", err)
		}
	} else {
		// plusieurs articles => on telecharge en parallele
		fmt.Printf("Telechargement de %d articles en parallele...\n", len(articles))
		webops.AnalyzeArticlesParallel(articles, cfg.WikiLang, cfg.OutDir)
	}
}

// ---- Choix D ----

func menuProcessOps() {
	for {
		fmt.Println("\n  --- ProcOps ---")
		fmt.Println("  [1] Lister les processus")
		fmt.Println("  [2] Rechercher")
		fmt.Println("  [3] Arreter un processus")
		fmt.Println("  [R] Retour")

		switch strings.ToUpper(readLine("  Choix : ")) {
		case "1":
			procs, err := procops.ListProcesses(cfg.ProcessTopN)
			if err != nil {
				fmt.Println("Erreur:", err)
				continue
			}
			procops.PrintProcesses(procs)

		case "2":
			kw := readLine("  Mot-cle : ")
			procs, err := procops.SearchProcesses(kw, cfg.ProcessTopN)
			if err != nil {
				fmt.Println("Erreur:", err)
				continue
			}
			if len(procs) == 0 {
				fmt.Println("  Aucun processus trouve.")
			} else {
				procops.PrintProcesses(procs)
			}

		case "3":
			pidStr := readLine("  PID : ")
			pid, err := strconv.Atoi(pidStr)
			if err != nil {
				fmt.Println("  PID invalide.")
				continue
			}
			if err := procops.KillProcess(pid, cfg.OutDir, reader); err != nil {
				fmt.Println("  Erreur:", err)
			}

		case "R":
			return
		default:
			fmt.Println("  Choix invalide.")
		}
	}
}

// ---- Choix E ----

func menuSecureOps() {
	for {
		fmt.Println("\n  --- SecureOps ---")
		fmt.Println("  [1] Verrouiller un fichier")
		fmt.Println("  [2] Deverrouiller")
		fmt.Println("  [3] Passer en lecture seule")
		fmt.Println("  [4] Restaurer lecture/ecriture")
		fmt.Println("  [5] Verifier permissions")
		fmt.Println("  [R] Retour")

		switch strings.ToUpper(readLine("  Choix : ")) {
		case "1":
			runSecureFileAction(func(p string) error { return secureops.LockFile(p, cfg.OutDir, reader) })
		case "2":
			runSecureFileAction(func(p string) error { return secureops.UnlockFile(p, cfg.OutDir, reader) })
		case "3":
			runSecureFileAction(func(p string) error { return secureops.SetReadOnly(p, cfg.OutDir) })
		case "4":
			runSecureFileAction(func(p string) error { return secureops.SetReadWrite(p, cfg.OutDir) })
		case "5":
			runSecureFileAction(secureops.CheckPermissions)
		case "R":
			return
		default:
			fmt.Println("  Choix invalide.")
		}
	}
}

// ---- Choix F ----

func menuContainerOps() {
	containers, err := infraops.ListContainers()
	if err != nil {
		fmt.Println("Erreur:", err)
		return
	}
	infraops.PrintContainers(containers)

	if len(containers) > 0 {
		name := readLine("Stats d'un conteneur (nom ou ID, vide pour passer) : ")
		if name != "" {
			if err := infraops.ContainerStats(name); err != nil {
				fmt.Println("Erreur:", err)
			}
		}
	}
}

// ---- Choix G ----

func menuHealthCheck() {
	fmt.Println("\n--- Etat du disque ---")
	if err := infraops.CheckDiskSpace(); err != nil {
		fmt.Println("Erreur:", err)
	}
}

// ---- Choix H - scan parallele avec goroutines ----

func menuParallelScan() {
	dir := readLineDefault("Dossier a scanner", cfg.BaseDir)

	files, err := fileops.FindTxtFiles(dir)
	if err != nil {
		fmt.Println("Erreur:", err)
		return
	}
	if len(files) == 0 {
		fmt.Println("Aucun fichier .txt trouve.")
		return
	}

	fmt.Printf("Scan parallele de %d fichiers...\n", len(files))

	type scanResult struct {
		path  string
		lines int
		words int
		err   error
	}

	var wg sync.WaitGroup
	ch := make(chan scanResult, len(files))

	for _, f := range files {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			lines, err := fileops.ReadLines(p)
			if err != nil {
				ch <- scanResult{path: p, err: err}
				return
			}
			words, err := fileops.ExtractWords(p)
			if err != nil {
				ch <- scanResult{path: p, err: err}
				return
			}
			ch <- scanResult{path: p, lines: len(lines), words: len(words)}
		}(f)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	fmt.Printf("\n  %-30s %8s %8s\n", "FICHIER", "LIGNES", "MOTS")
	fmt.Printf("  %s\n", strings.Repeat("-", 50))
	for r := range ch {
		if r.err != nil {
			fmt.Printf("  %-30s %8s %8s\n", r.path, "ERR", "ERR")
			fmt.Printf("    -> %v\n", r.err)
			continue
		}
		fmt.Printf("  %-30s %8d %8d\n", r.path, r.lines, r.words)
	}
}

// ---- saisie utilisateur ----

func runStep(title string, fn func() error) {
	fmt.Printf("\n--- %s ---\n", title)
	if err := fn(); err != nil {
		fmt.Println("Erreur:", err)
	}
}

func runSecureFileAction(fn func(string) error) {
	path := readLineDefault("  Fichier", cfg.DefaultFile)
	if err := fn(path); err != nil {
		fmt.Println("  Erreur:", err)
	}
}

func readLine(prompt string) string {
	fmt.Print(prompt)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func readLineDefault(prompt, def string) string {
	fmt.Printf("%s [%s] : ", prompt, def)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return def
	}
	return line
}

func readInt(prompt string, def int) int {
	s := readLineDefault(prompt, strconv.Itoa(def))
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

func readIntMin(prompt string, def, min int) int {
	n := readInt(prompt, def)
	if n < min {
		return def
	}
	return n
}
