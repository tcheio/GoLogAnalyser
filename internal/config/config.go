package config

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "os"
)

type Entry struct {
    ID   string `json:"id"`
    Path string `json:"path"`
    Type string `json:"type"`
}

func Load(path string) ([]Entry, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("impossible d'ouvrir la config: %w", err)
    }
    defer f.Close()

    data, err := io.ReadAll(f)
    if err != nil {
        return nil, fmt.Errorf("lecture config échouée: %w", err)
    }

    var entries []Entry
    if err := json.Unmarshal(data, &entries); err != nil {
        var se *json.SyntaxError
        if errors.As(err, &se) {
            return nil, &ParseError{Offset: se.Offset, Err: err}
        }
        return nil, &ParseError{Offset: 0, Err: err}
    }
    return entries, nil
}

func Append(path string, e Entry) error {
    entries, err := Load(path)
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            entries = []Entry{e}
        } else {
            return err
        }
    } else {
        entries = append(entries, e)
    }

    out, err := json.MarshalIndent(entries, "", "  ")
    if err != nil {
        return fmt.Errorf("échec de sérialisation: %w", err)
    }
    if err := os.WriteFile(path, out, 0o644); err != nil {
        return fmt.Errorf("échec d'écriture config: %w", err)
    }
    return nil
}

type ParseError struct {
    Offset int64
    Err    error
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("erreur de parsing JSON (offset=%d): %v", e.Offset, e.Err)
}

func (e *ParseError) Unwrap() error { return e.Err }
