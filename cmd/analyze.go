package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/tcheio/GOLogAnalyser/internal/analyzer"
	"github.com/tcheio/GOLogAnalyser/internal/config"
	"github.com/tcheio/GOLogAnalyser/internal/reporter"
)

var (
	configPath string
	outputPath string
	statusOnly string
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyse les logs listés dans un fichier de configuration JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		if configPath == "" {
			return fmt.Errorf("veuillez fournir un fichier de configuration via --config")
		}

		entries, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("échec de lecture de la configuration: %w", err)
		}

		results := analyzer.Run(entries, verbose)

		if statusOnly != "" {
			want := strings.ToUpper(statusOnly)
			filtered := make([]analyzer.Result, 0, len(results))
			for _, r := range results {
				if strings.EqualFold(r.Status, want) {
					filtered = append(filtered, r)
				}
			}
			results = filtered
		}

		printSummary(results)

		if outputPath != "" {
			finalPath, wErr := reporter.ExportJSON(results, outputPath, true)
			if wErr != nil {
				return fmt.Errorf("échec d'export JSON: %w", wErr)
			}
			fmt.Printf("\nRapport écrit: %s\n", finalPath)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().StringVarP(&configPath, "config", "c", "", "chemin du fichier de configuration JSON")
	analyzeCmd.Flags().StringVarP(&outputPath, "output", "o", "", "chemin du rapport JSON (répertoires créés au besoin)")
	analyzeCmd.Flags().StringVar(&statusOnly, "status", "", "filtre les résultats par statut (OK ou FAILED)")
}

func printSummary(results []analyzer.Result) {
	if len(results) == 0 {
		fmt.Println("Aucun résultat à afficher.")
		return
	}

	fmt.Println("\n--- Résumé de l'analyse ---")
	for _, r := range results {
		msg := fmt.Sprintf("[%s] %s — %s — %s", r.LogID, r.FilePath, r.Status, r.Message)
		fmt.Println(msg)
		if r.ErrorDetails != "" {
			fmt.Println("  ↳ erreur:", r.ErrorDetails)
		}
	}

	var target *json.SyntaxError
	if errors.As(nil, &target) {
	}

	_ = time.Now()
}
