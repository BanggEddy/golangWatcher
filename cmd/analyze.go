package cmd

import (
	"errors"
	"fmt"
	"sync"

	"github.com/axellelanca/gowatcher_correction/internal/checker"
	"github.com/axellelanca/gowatcher_correction/internal/config"
	"github.com/axellelanca/gowatcher_correction/internal/reporter"
	"github.com/spf13/cobra"
)

var (
	inputFilePath  string
	outputFilePath string
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Vérifie l'accessibilité d'une liste d'URLs.",
	Long:  `La commande 'check' parcourt une liste prédéfinie d'URLs et affiche leur statut d'accessibilité en utilisant des goroutines pour la concurrence.`,
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
			fmt.Printf("Erreur lors du chargement des URLs: %v\n", err)
			return
		}

		if len(targets) == 0 {
			fmt.Println("Aucune URL à vérifier trouvée dans le fichier d'entrée.")
			return
		}

		var wg sync.WaitGroup
		resultsChan := make(chan checker.CheckResult, len(targets)) // Canal pour collecter les résultats

		wg.Add(len(targets))
		for _, target := range targets {
			go func(t config.InputTarget) {
				defer wg.Done()
				result := checker.CheckURLSync(t)
				resultsChan <- result // Envoyer le resultat au canal
			}(target)
		}

		wg.Wait()          // Attendre que toutes les goroutines aient fini
		close(resultsChan) // Fermer le canal après que tous les résultats ont été envoyés

		var finalReport []checker.ReportEntry
		for res := range resultsChan { // Récupérer tous les résultats du canal
			reportEntry := checker.ConvertToReportEntry(res)
			finalReport = append(finalReport, reportEntry)

			// Affichage immédiat comme avant
			if res.Err != nil {
				var unreachable *checker.UnreachableURLError
				if errors.As(res.Err, &unreachable) {
					fmt.Printf("🚫 %s (%s) est inaccessible : %v\n", res.InputTarget.Name, unreachable.URL, unreachable.Err)
				} else {
					fmt.Printf("❌ %s (%s) : erreur - %v\n", res.InputTarget.Name, res.InputTarget.URL, res.Err)
				}
			} else {
				fmt.Printf("✅ %s (%s) : OK - %s\n", res.InputTarget.Name, res.InputTarget.URL, res.Status)
			}
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
	rootCmd.AddCommand(checkCmd)

	// Ici, vous pouvez ajouter des drapeaux (flags) spécifiques à la commande 'check'.
	// Ces drapeaux ne seront disponibles que lorsque la commande 'check' est utilisée.
	// Exemple (commenté) : checkCmd.Flags().StringVarP(&sourceFile, "source", "s", "", "Fichier contenant les URLs à vérifier")

	// Ajout des drapeaux spécifiques à la commande 'check'
	checkCmd.Flags().StringVarP(&inputFilePath, "input", "i", "", "Chemin vers le fichier JSON d'entrée contenant les URLs")
	checkCmd.Flags().StringVarP(&outputFilePath, "output", "o", "", "Chemin vers le fichier JSON de sortie pour les résultats (optionnel)")

	// Marquer le drapeau "input" comme obligatoire
	checkCmd.MarkFlagRequired("input")
}
