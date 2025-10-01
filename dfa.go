package main

import (
	"fmt"
	"strings"
)

type DFA struct {
	States       [26]State `json:"states"`
	CurrentState uint      `json:"currentState"`
	StartState   uint      `json:"startState"`
}

type State struct {
	IsValid bool `json:"isValid"`
	IsStart bool `json:"isStart"`
	IsEnd   bool `json:"IsEnd"`
	NextA   uint `json:"nextA"`
	NextB   uint `json:"nextB"`
}

func (dfa *DFA) GetCurrentState() State {
	return dfa.States[dfa.CurrentState]
}

// Move advances the DFA based on the `read` input.
//
// This function expects `read` to be either 0 or 1.
func (dfa *DFA) Move(read uint) error {
	if read == 0 {
		dfa.CurrentState = dfa.GetCurrentState().NextA
	} else {
		dfa.CurrentState = dfa.GetCurrentState().NextB
	}

	if !dfa.GetCurrentState().IsValid {
		return fmt.Errorf("went to an invalid state.")
	}

	return nil
}

// Constants for readability in line parsing
const (
	colStateType = iota
	colStateName
	colNextA
	colNextB
)

// parseLetterState converts uppercase letters to their corresponding 0-index.
func parseLetterState(ch byte) (uint, error) {
	if ch < 'A' || ch > 'Z' {
		return 0, fmt.Errorf("invalid state name `%c`", ch)
	}

	return uint(ch - 'A'), nil
}

func NewDFA(in string) (*DFA, error) {
	dfa := &DFA{}
	isStartSet := false

	lines := strings.Split(in, "\n")
	for line_no, line := range lines {
		cols := strings.Split(line, ",")
		if len(cols) != 4 {
			return nil, fmt.Errorf("%d:%v does not have 4 columns", line_no+1, cols)
		}

		// identify which state is being set ======================================
		idx, err := parseLetterState(cols[colStateName][0])
		if err != nil {
			return nil, err
		}

		state := &dfa.States[idx]

		// identify the state type (start, end) ===================================
		switch cols[colStateType] {
		case "-":
			if isStartSet {
				return nil, fmt.Errorf("multiple start states parsed")
			}
			state.IsStart = true
			dfa.StartState = idx
			isStartSet = true
		case "+":
			state.IsEnd = true
		}

		// set next states ========================================================
		nextA, err := parseLetterState(cols[colNextA][0])
		if err != nil {
			return nil, err
		}
		state.NextA = nextA

		nextB, err := parseLetterState(cols[colNextB][0])
		if err != nil {
			return nil, err
		}
		state.NextB = nextB

		state.IsValid = true
	}

	return dfa, nil
}
