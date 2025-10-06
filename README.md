# DFA Evaluator

A desktop application for evaluating strings against Deterministic Finite Automata (DFA) with a user-friendly graphical interface.
## Demo

https://github.com/user-attachments/assets/a5590a40-b8fa-4caa-8c1f-eb2fb5154f4d

## Prerequisites

- Follow platform-specific Wails installation: https://wails.io/docs/gettingstarted/installation
- Check that you have Go, npm and other OS specific files installed.
- Verify the installations with `wails doctor`

## Quick Start

**Run Built App**: Use the pre-built executable called ```dfa-evaluator.app``` in `build/bin/` directory


**Development Mode**: Run `wails dev` for live development with hot reload

**Building**: Run `wails build --platform` for production package depending on your OS

## Key Features

- Load DFA tables (`.dfa`) and input strings (`.in`)
- Validate strings against DFA rules (VALID/INVALID)
- Auto-generate output files (`.out`)
- Interactive GUI with real-time feedback
- Comprehensive error handling

## How It Works

1. Load DFA transition table and input strings
2. Validate DFA format and completeness
3. Process each string through the DFA
4. Generate results in UI and output file

## Input Files

### DFA Files (`.dfa`)
CSV format with:
- **Line 1**: Two alphabet symbols (e.g., `0,1`)
- **Other lines**: State transitions with 4 elements:
  - Type: `-` (start), `+` (final), or blank
  - Source state (A-Z)
  - Destination for symbol 1
  - Destination for symbol 2

### Input Strings (`.in`)
- One string per line
- Uses alphabet symbols from DFA

## Example

**transitions.dfa**
```
0,1
-,A,B,A
,B,B,C
+,C,B,A
```

**strings.in**
```
110011
0111110001
```

**strings.out** (generated)
```
INVALID
VALID
```

## Usage

1. Launch application
2. Load DFA file (`.dfa`)
3. Load input strings (`.in`)
4. Click "Process" to evaluate
5. View results and check generated `.out` file

## Project Structure

- `dfa.go` - DFA logic and validation
- `evaluator.go` - String evaluation engine  
- `app.go` - Wails application bridge
- `frontend/` - Svelte UI components
