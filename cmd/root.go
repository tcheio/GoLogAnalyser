package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var (
    verbose bool
)

var rootCmd = &cobra.Command{
    Use:   "loganalyzer",
    Short: "Analyse concurrente et robuste de fichiers de logs",
    Long: `GoLog Analyzer (loganalyzer) est un outil CLI pour analyser des fichiers de logs en parallèle.`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func init() {
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "affiche plus de détails")
}
