package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// InputTarget représente une URL à vérifier lue depuis le fichier JSON d'entrée.
type InputTarget struct {
	Name  string `json:"name"`  // Nom descriptif de l'URL
	URL   string `json:"url"`   // L'URL à vérifier
	Owner string `json:"owner"` // Propriétaire ou contact lié à l'URL
}

// LoadTargetsFromFile lit une liste d'InputTarget à partir d'un fichier JSON.
func LoadTargetsFromFile(filePath string) ([]InputTarget, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("impossible de lire le fichier %s: %w", filePath, err)
	}

	var targets []InputTarget
	if err := json.Unmarshal(data, &targets); err != nil {
		return nil, fmt.Errorf("impossible de lire le fichier %s: %w", filePath, err)
	}
	return targets, nil
}

// SaveTargetsToFile écrit une liste d'InputTarget dans un fichier JSON.
func SaveTargetsToFile(filePath string, targets []InputTarget) error {
	data, err := json.MarshalIndent(targets, "", "  ")
	if err != nil {
		return fmt.Errorf("impossible de lire le fichier %s: %w", filePath, err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("impossible de lire le fichier %s: %w", filePath, err)
	}
	return nil
}
