package main

import (
	"context"
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

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	content := strings.ReplaceAll(string(bytes), "\r\n", "\n")
	a.evaluator, err = NewEvaluator(strings.TrimSpace(content))
	if err != nil {
		return "", err
	}

	return content, nil
}
