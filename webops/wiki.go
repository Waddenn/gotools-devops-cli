package webops

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

func FetchArticle(article, lang string) (string, error) {
	url := fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", lang, article)
	fmt.Printf("  Recuperation de %s...\n", url)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("erreur creation requete HTTP: %w", err)
	}
	req.Header.Set("User-Agent", "GoTools/1.0 (+M1 DevOps project)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", lang)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erreur HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("statut HTTP %d pour %s", resp.StatusCode, url)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erreur parsing HTML: %w", err)
	}

	// on recupere le texte de chaque <p> dans le contenu principal
	var paragraphs []string
	doc.Find("#mw-content-text p").Each(func(_ int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			paragraphs = append(paragraphs, text)
		}
	})

	return strings.Join(paragraphs, "\n\n"), nil
}

// AnalyzeArticle telecharge un article, affiche des stats et sauvegarde dans out/
func AnalyzeArticle(article, lang, outDir string) error {
	text, err := FetchArticle(article, lang)
	if err != nil {
		return err
	}
	if text == "" {
		return fmt.Errorf("aucun contenu pour '%s'", article)
	}

	// stat 1 : nb mots (sans les numeriques)
	words := extractWordsFromText(text)
	totalLen := 0
	for _, w := range words {
		totalLen += len(w)
	}
	avg := 0.0
	if len(words) > 0 {
		avg = float64(totalLen) / float64(len(words))
	}
	fmt.Printf("  Mots (hors numeriques) : %d\n", len(words))
	fmt.Printf("  Longueur moyenne       : %.1f\n", avg)

	// stat 2 : nb de paragraphes
	paras := strings.Split(text, "\n\n")
	fmt.Printf("  Paragraphes            : %d\n", len(paras))

	// sauvegarde
	outPath := filepath.Join(outDir, "wiki_"+safeFilePart(article)+".txt")
	content := fmt.Sprintf("=== %s ===\nMots: %d | Moy: %.1f | Paragraphes: %d\n\n%s\n",
		article, len(words), avg, len(paras), text)

	if err := os.WriteFile(outPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("erreur ecriture: %w", err)
	}
	fmt.Printf("  -> Sauvegarde dans %s\n", outPath)
	return nil
}

// AnalyzeArticlesParallel lance le telechargement de plusieurs articles en meme temps
func AnalyzeArticlesParallel(articles []string, lang, outDir string) {
	var wg sync.WaitGroup
	results := make(chan string, len(articles))

	for _, article := range articles {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			if err := AnalyzeArticle(a, lang, outDir); err != nil {
				results <- fmt.Sprintf("  Echec [%s] : %v", a, err)
			} else {
				results <- fmt.Sprintf("  Termine [%s]", a)
			}
		}(article)
	}

	// on ferme le channel quand toutes les goroutines sont finies
	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("\n  Resultats :")
	for r := range results {
		fmt.Println(r)
	}
}

func extractWordsFromText(text string) []string {
	var words []string
	for _, w := range strings.Fields(text) {
		if !isNumeric(w) {
			words = append(words, w)
		}
	}
	return words
}

func isNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) && r != '.' && r != ',' {
			return false
		}
	}
	return true
}

func safeFilePart(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "article"
	}

	var b strings.Builder
	for _, r := range s {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			b.WriteRune(r)
		case r == '-', r == '_':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}

	name := strings.Trim(b.String(), "_")
	if name == "" {
		return "article"
	}
	return name
}
