package audit

import (
	"fmt"
	"os"
	"time"
)

// Log ajoute une ligne horodatee dans out/audit.log
func Log(outDir, action string) {
	path := outDir + "/audit.log"
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur audit log: %v\n", err)
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), action)
}
