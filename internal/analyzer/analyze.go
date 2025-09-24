// internal/checker/check.go (ajout à la structure CheckResult ou nouvelle structure)
// Pour une meilleure clarté dans le rapport final, nous allons légèrement modifier CheckResult
// pour inclure les champs "Name" et "Owner" dès le départ.

package checker

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/axellelanca/gowatcher_correction/internal/config"
)

// ReportEntry représente une entrée dans le rapport final JSON.
type ReportEntry struct {
	Name   string
	URL    string
	Owner  string
	Status string // "OK", "Inaccessible", "Error"
	ErrMsg string // Message d'erreur, omis si vide
}

// CheckResult (modifié ou nouvelle version pour le workflow)
// Cette structure peut être utilisée en interne pour le résultat immédiat.
// Nous la convertirons en ReportEntry pour l'export.
type CheckResult struct {
	InputTarget config.InputTarget
	Status      string
	Err         error
}

func CheckURLSync(target config.InputTarget) CheckResult {
	// Timeout court pour éviter de bloquer trop longtemps
	client := http.Client{Timeout: 3 * time.Second}

	resp, err := client.Get(target.URL)
	if err != nil {
		return CheckResult{
			InputTarget: target,
			Err:         &UnreachableURLError{URL: target.URL, Err: err},
		}
	}
	defer resp.Body.Close()

	return CheckResult{InputTarget: target, Status: resp.Status}
}

// ConvertToCheckReport convertit un CheckResult interne en ReportEntry pour l'exportation.
func ConvertToReportEntry(res CheckResult) ReportEntry {
	report := ReportEntry{
		Name:   res.InputTarget.Name,
		URL:    res.InputTarget.URL,
		Owner:  res.InputTarget.Owner,
		Status: res.Status, // Statut par défaut
	}

	if res.Err != nil {
		var unreachable *UnreachableURLError
		if errors.As(res.Err, &unreachable) {
			report.Status = "Inaccessible"
			report.ErrMsg = fmt.Sprintf("Unreachable URL: %v", unreachable.Err)
		} else {
			report.Status = "Error"
			report.ErrMsg = fmt.Sprintf("Erreur générique: %v", res.Err)
		}
	}
	
	return report
}
