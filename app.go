package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	evaluator   *Evaluator
	inputPath   string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type DfaLoadResult struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

func (a *App) LoadDotDfa() (*DfaLoadResult, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// load file path
	filename, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Select `.dfa` File",
		DefaultDirectory: currentDir,
		Filters: []runtime.FileFilter{
			{DisplayName: "Deterministic Finite Automata templates (*.dfa)", Pattern: "*.dfa"},
		},
	})
	if err != nil {
		return nil, err
	}

	// load file into memory (should be fine since there can be at most 27 lines)
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// trim whitespace
	content := strings.ReplaceAll(string(bytes), "\r\n", "\n")
	a.evaluator, err = NewEvaluator(strings.TrimSpace(content))
	if err != nil {
		return nil, err
	}

	return &DfaLoadResult{
		Filename: filepath.Base(filename),
		Content:  content,
	}, nil
}

type InLoadResult struct {
	Filename   string   `json:"filename"`
	InputLines []string `json:"inputLines"`
}

func (a *App) LoadDotIn() (*InLoadResult, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// load file path
	filename, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Select `.in` File",
		DefaultDirectory: currentDir,
		Filters: []runtime.FileFilter{
			{DisplayName: "Input File (*.in)", Pattern: "*.in"},
		},
	})
	if err != nil {
		return nil, err
	}

	// load file into memory
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// trim whitespace
	content := strings.ReplaceAll(string(bytes), "\r\n", "\n")
	inputLines := strings.Split(strings.TrimSpace(content), "\n")

	// Store the full path for later saving
	a.inputPath = filename

	return &InLoadResult{
		Filename:   filepath.Base(filename),
		InputLines: inputLines,
	}, nil
}

func (a *App) EvaluateInput(inputLines []string) ([]bool, error) {
	if a.evaluator == nil {
		return []bool{}, fmt.Errorf("no DFA loaded")
	}

	results := make([]bool, len(inputLines))
	for i, line := range inputLines {
		accepted, err := a.evaluator.Evaluate(line)
		if err != nil {
			return []bool{}, err
		}

		results[i] = accepted
	}

	return results, nil
}

func (a *App) SaveOutput(results []bool) error {
	if a.inputPath == "" {
		return fmt.Errorf("no input file loaded")
	}

	// Generate output path by replacing .in extension with .out
	outputPath := strings.TrimSuffix(a.inputPath, filepath.Ext(a.inputPath)) + ".out"

	// Build output content
	var output strings.Builder
	for _, result := range results {
		if result {
			output.WriteString("VALID\n")
		} else {
			output.WriteString("INVALID\n")
		}
	}

	// Write to file
	err := os.WriteFile(outputPath, []byte(output.String()), 0644)
	if err != nil {
		return err
	}

	return nil
}