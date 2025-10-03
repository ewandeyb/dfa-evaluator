package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx       context.Context
	evaluator *Evaluator
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

func (a *App) LoadDotDfa() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
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
		return "", err
	}

	// load file into memory (should be fine since there can be at most 27 lines)
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// trim whitespace
	content := strings.ReplaceAll(string(bytes), "\r\n", "\n")
	a.evaluator, err = NewEvaluator(strings.TrimSpace(content))
	if err != nil {
		return "", err
	}

	return content, nil
}

func (a *App) LoadDotIn() ([]string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return []string{}, err
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
		return []string{}, err
	}

	// load file into memory
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return []string{}, err
	}

	// trim whitespace
	content := strings.ReplaceAll(string(bytes), "\r\n", "\n")
	inputLines := strings.Split(strings.TrimSpace(content), "\n")
	

	return inputLines, nil
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