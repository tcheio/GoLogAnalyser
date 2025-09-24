package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tcheio/GOLogAnalyser/internal/analyzer"
)

func ExportJSON(results []analyzer.Result, outPath string, withDatePrefix bool) (string, error) {
	if len(results) == 0 {
		return "", fmt.Errorf("aucun résultat à exporter")
	}

	dir := filepath.Dir(outPath)
	base := filepath.Base(outPath)

	if withDatePrefix {
		date := time.Now().Format("060102")
		base = fmt.Sprintf("%s_%s", date, base)
	}

	finalPath := filepath.Join(dir, base)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	tmpPath := finalPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return "", fmt.Errorf("write: %w", err)
	}
	if err := os.Rename(tmpPath, finalPath); err != nil {
		return "", fmt.Errorf("rename: %w", err)
	}

	return finalPath, nil
}

func FilterByStatus(results []analyzer.Result, status string) []analyzer.Result {
	status = strings.ToUpper(status)
	out := make([]analyzer.Result, 0, len(results))
	for _, r := range results {
		if strings.EqualFold(r.Status, status) {
			out = append(out, r)
		}
	}
	return out
}
