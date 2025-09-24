package cmd

import (
	"errors"
	"fmt"
	"sync"

	"github.com/BanggEddy/golangWatcher/internal/analyzer"
	"github.com/BanggEddy/golangWatcher/internal/config"
	"github.com/BanggEddy/golangWatcher/internal/reporter"
	"github.com/spf13/cobra"
)

var (
	inputFilePath  string
	outputFilePath string
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyse une liste de fichiers de logs.",
	Long:  `La commande 'analyze' lit une liste de logs depuis un fichier JSON, les analyse en parallèle et exporte un rapport.`,
	Run: func(cmd *cobra.Command, args []string) { // La fonction `Run` est le cœur de la sous-commande.
		// Elle est exécutée lorsque l'utilisateur tape `gowatcher check`.
		// `cmd` représente la commande elle-même, `args` sont les arguments positionnels passés après la commande.

		if inputFilePath == "" {
			fmt.Println("Erreur: le chemin du fichier d'entrée (--input) est obligatoire.")
			return
		}

		// Charger les cibles depuis le fichier JSON d'entrée
		targets, err := config.LoadTargetsFromFile(inputFilePath)
		if err != nil {
			fmt.Printf("Erreur lors du chargement des logs: %v\n", err)
			return
		}

		if len(targets) == 0 {
			fmt.Println("0 log trouvé dans le fichier d'entrée.")
			return
		}

		var wg sync.WaitGroup
		resultsChan := make(chan analyzer.LogResult, len(targets)) // Canal pour collecter les résultats

		wg.Add(len(targets))
		for _, target := range targets {
			go func(t config.LogTarget) {
				defer wg.Done()
				result := analyzer.AnalyzeLog(t)
				resultsChan <- result // Envoyer le resultat au canal
			}(target)
		}

		wg.Wait()          // Attendre que toutes les goroutines aient fini
		close(resultsChan) // Fermer le canal après que tous les résultats ont été envoyés

		var finalReport []analyzer.LogResult
		for res := range resultsChan { // Récupérer tous les résultats du canal
			if res.Status == "FAILED" {
				var notfound *analyzer.FileNotFoundError
				var parseErr *analyzer.ParseLogError
				if errors.As(errors.New(res.ErrorDetail), &notfound) {
					fmt.Printf("🚫 %s (%s) : %s\n", res.LogID, res.FilePath, res.Message)
				} else if errors.As(errors.New(res.ErrorDetail), &parseErr) {
					fmt.Printf("⚠️ %s (%s) : %s\n", res.LogID, res.FilePath, res.Message)
				} else {
					fmt.Printf("❌ %s (%s) : %s - %s\n", res.LogID, res.FilePath, res.Message, res.ErrorDetail)
				}
			} else {
				fmt.Printf("✅ %s (%s) : %s\n", res.LogID, res.FilePath, res.Message)
			}
			finalReport = append(finalReport, res)
		}

		// Exporter les résultats si outputFilePath est spécifié
		if outputFilePath != "" {
			err := reporter.ExportResultsToJsonFile(outputFilePath, finalReport)
			if err != nil {
				fmt.Printf("Erreur lors de l'exportation des résultats: %v\n", err)
			} else {
				fmt.Printf("✅ Résultats exportés vers %s\n", outputFilePath)
			}
		}
	},
}

// init() est une fonction spéciale de Go, exécutée lors de l'initialisation du package.
func init() {
	// Cette ligne est cruciale : elle "ajoute" la sous-commande `checkCmd` à la commande racine `rootCmd`.
	// C'est ainsi que Cobra sait que 'check' est une commande valide sous 'gowatcher'.
	rootCmd.AddCommand(analyzeCmd)
	
	// Ici, vous pouvez ajouter des drapeaux (flags) spécifiques à la commande 'check'.
	// Ces drapeaux ne seront disponibles que lorsque la commande 'check' est utilisée.
	// Exemple (commenté) : checkCmd.Flags().StringVarP(&sourceFile, "source", "s", "", "Fichier contenant les URLs à vérifier")

	// Ajout des drapeaux spécifiques à la commande 'check'
	analyzeCmd.Flags().StringVarP(&inputFilePath, "input", "i", "", "Chemin vers le fichier JSON d'entrée contenant les logs")
	analyzeCmd.Flags().StringVarP(&outputFilePath, "output", "o", "", "Chemin vers le fichier JSON de sortie pour les résultats (optionnel)")
	
	// Marquer le drapeau "input" comme obligatoire
	analyzeCmd.MarkFlagRequired("input")
}