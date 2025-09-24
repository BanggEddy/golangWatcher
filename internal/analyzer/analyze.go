// internal/checker/check.go (ajout à la structure CheckResult ou nouvelle structure)
// Pour une meilleure clarté dans le rapport final, nous allons légèrement modifier CheckResult
// pour inclure les champs "Name" et "Owner" dès le départ.

package analyzer

import (
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/axellelanca/gowatcher_correction/internal/config"
)

// LogResult représente le résultat de l'analyse d'un log.
type LogResult struct {
	LogID       string `json:"log_id"`
	FilePath    string `json:"file_path"`
	Status      string `json:"status"`
	Message     string `json:"message"`
	ErrorDetail string `json:"error_details"`
}

func AnalyzeLog(target config.LogTarget) LogResult {
	_, err := os.Stat(target.Path)
	if err != nil {
		return LogResult{
			LogID:       target.ID,
			FilePath:    target.Path,
			Status:      "FAILED",
			Message:     "Fichier introuvable.",
			ErrorDetail: (&FileNotFoundError{Path: target.Path, Err: err}).Error(),
		}
	}

	data, err := os.ReadFile(target.Path)
	if err != nil {
		return LogResult{
			LogID:       target.ID,
			FilePath:    target.Path,
			Status:      "FAILED",
			Message:     "Impossible de lire le fichier.",
			ErrorDetail: err.Error(),
		}
	}

	if strings.Contains(string(data), "INVALID") {
		return LogResult{
			LogID:       target.ID,
			FilePath:    target.Path,
			Status:      "FAILED",
			Message:     "Erreur de parsing.",
			ErrorDetail: (&ParseLogError{Path: target.Path, Msg: "ligne corrompue"}).Error(),
		}
	}

	time.Sleep(time.Duration(rand.Intn(150)+50) * time.Millisecond)

	return LogResult{
		LogID:    target.ID,
		FilePath: target.Path,
		Status:   "OK",
		Message:  "Analyse terminée avec succès.",
	}
}
