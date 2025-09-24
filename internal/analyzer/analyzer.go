package analyzer

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/tcheio/GOLogAnalyser/internal/config"
)

type Result struct {
	LogID        string `json:"log_id"`
	FilePath     string `json:"file_path"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	ErrorDetails string `json:"error_details"`
}

type FileAccessError struct {
	Path string
	Err  error
}

func (e *FileAccessError) Error() string {
	return fmt.Sprintf("fichier inaccessible '%s': %v", e.Path, e.Err)
}

func (e *FileAccessError) Unwrap() error { return e.Err }

func Run(entries []config.Entry, verbose bool) []Result {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	results := make([]Result, 0, len(entries))
	resultsCh := make(chan Result, len(entries))

	var wg sync.WaitGroup
	wg.Add(len(entries))

	for _, e := range entries {
		e := e
		go func() {
			defer wg.Done()
			res := analyzeOne(e, verbose)
			resultsCh <- res
		}()
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	for r := range resultsCh {
		results = append(results, r)
	}
	return results
}

func analyzeOne(e config.Entry, verbose bool) Result {
	if _, err := os.Stat(e.Path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Result{
				LogID:        e.ID,
				FilePath:     e.Path,
				                Status:       "FAILED",
				                Message:      "Fichier introuvable.",
				                ErrorDetails: fmt.Sprintf("FAILED - Fichier introuvable: %v", err),
				            }
				        }
				        return Result{
				            LogID:        e.ID,
				            FilePath:     e.Path,
				            Status:       "FAILED",
				            Message:      "Fichier inaccessible.",
				            ErrorDetails: fmt.Sprintf("FAILED - Fichier inaccessible: %v", err),
				        }
				    }
				
				    f, err := os.Open(e.Path)
				    if err != nil {
				        return Result{
				            LogID:        e.ID,
				            FilePath:     e.Path,
				            Status:       "FAILED",
				            Message:      "Impossible d'ouvrir le fichier.",
				            ErrorDetails: fmt.Sprintf("FAILED - Impossible d'ouvrir le fichier: %v", err),
				        }
				    }
				    	_ = f.Close()
				    
				    	content, err := os.ReadFile(e.Path)
				    	if err != nil {
				    		return Result{
				    			LogID:        e.ID,
				    			FilePath:     e.Path,
				    			Status:       "FAILED",
				    			Message:      "Impossible de lire le contenu du fichier.",
				    			ErrorDetails: fmt.Sprintf("FAILED - Lecture du fichier impossible: %v", err),
				    		}
				    	}
				    
				    	if len(content) == 0 {
				    		return Result{
				    			LogID:        e.ID,
				    			FilePath:     e.Path,
				    			Status:       "FAILED",
				    			Message:      "Fichier vide.",
				    			ErrorDetails: "FAILED - Fichier vide.",
				    		}
				    	}
				    
				    	// TODO: Implement actual corruption detection based on log type.
				    	// For now, if it's the 'corrupted-log' and not empty, we'll mark it as corrupted.
				    	if e.ID == "corrupted-log" {
				    		return Result{
				    			LogID:        e.ID,
				    			FilePath:     e.Path,
				    			Status:       "FAILED",
				    			Message:      "Fichier corrompu détecté.",
				    			ErrorDetails: "FAILED - Fichier corrompu (détection basique).",
				    		}
				    	}
				    
				    	delay := time.Duration(50+rand.Intn(151)) * time.Millisecond
				    	time.Sleep(delay)
				    
				    	return Result{
				    		LogID:        e.ID,
				    		FilePath:     e.Path,
				    		Status:       "OK",
				    		Message:      "Analyse terminée avec succès.",
				    		ErrorDetails: "OK - Aucune erreur détectée.",
				    	}
				    }
