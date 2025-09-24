package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tcheio/GOLogAnalyser/internal/config"
)

var (
	addID   string
	addPath string
	addType string
	addFile string
)

var addLogCmd = &cobra.Command{
	Use:   "add-log",
	Short: "Ajoute une entrée de log au fichier de configuration JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		if addID == "" || addPath == "" || addType == "" || addFile == "" {
			return fmt.Errorf("--id, --path, --type et --file sont requis")
		}
		entry := config.Entry{ID: addID, Path: addPath, Type: addType}
		if err := config.Append(addFile, entry); err != nil {
			return err
		}
		fmt.Printf("Entrée ajoutée à %s: %s (%s)\n", addFile, addID, addPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addLogCmd)
	addLogCmd.Flags().StringVar(&addID, "id", "", "identifiant unique du log")
	addLogCmd.Flags().StringVar(&addPath, "path", "", "chemin vers le fichier de log")
	addLogCmd.Flags().StringVar(&addType, "type", "", "type de log (ex: nginx-access, custom-app)")
	addLogCmd.Flags().StringVar(&addFile, "file", "", "chemin du fichier config.json à modifier")
}
