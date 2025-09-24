package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// LogTarget représente un fichier log à analyser lu depuis le fichier JSON d'entrée.
type LogTarget struct {
	ID   string `json:"id"`
	Path string `json:"path"`
	Type string `json:"type"`
}

// LoadTargetsFromFile lit une liste de LogTarget à partir d'un fichier JSON.
func LoadTargetsFromFile(filePath string) ([]LogTarget, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("impossible de lire le fichier %s: %w", filePath, err)
	}
	var targets []LogTarget
	if err := json.Unmarshal(data, &targets); err != nil {
		return nil, fmt.Errorf("impossible de lire le fichier %s: %w", filePath, err)
	}
	return targets, nil
}
