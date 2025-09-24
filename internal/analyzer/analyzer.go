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
				ErrorDetails: (&FileAccessError{Path: e.Path, Err: err}).Error(),
			}
		}
		return Result{
			LogID:        e.ID,
			FilePath:     e.Path,
			Status:       "FAILED",
			Message:      "Fichier inaccessible.",
			ErrorDetails: (&FileAccessError{Path: e.Path, Err: err}).Error(),
		}
	}

	f, err := os.Open(e.Path)
	if err != nil {
		return Result{
			LogID:        e.ID,
			FilePath:     e.Path,
			Status:       "FAILED",
			Message:      "Impossible d'ouvrir le fichier.",
			ErrorDetails: (&FileAccessError{Path: e.Path, Err: err}).Error(),
		}
	}
	_ = f.Close()

	delay := time.Duration(50+rand.Intn(151)) * time.Millisecond
	time.Sleep(delay)

	return Result{
		LogID:        e.ID,
		FilePath:     e.Path,
		Status:       "OK",
		Message:      "Analyse terminée avec succès.",
		ErrorDetails: "",
	}
}
