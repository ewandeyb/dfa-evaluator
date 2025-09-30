package main

import (
	"fmt"
	"strings"
)

type Evaluator struct {
	Dfa    *DFA `json:"dfa"`
	AState byte `json:"aState"`
	BState byte `json:"bState"`
}

func NewEvaluator(input string) (*Evaluator, error) {
	evaluator := &Evaluator{}

	states_line, dfa_lines, found := strings.Cut(input, "\n")
	if !found {
		return nil, fmt.Errorf("no new line was found")
	}

	// set states ===============================================================
	states := strings.Split(states_line, ",")
	if len(states) != 2 {
		return nil, fmt.Errorf("%v does not have two alphabet variants", states)
	}
	evaluator.AState, evaluator.BState = states[0][0], states[1][0]
	if evaluator.AState == evaluator.BState {
		return nil, fmt.Errorf("alphabet should contain differing characters")
	}

	// set dfa ==================================================================
	dfa, err := NewDFA(dfa_lines)
	if err != nil {
		return nil, err
	}
	evaluator.Dfa = dfa

	return evaluator, nil
}

// `Normalize` converts a string input into an array of 0 and 1 using `AState` and `BState`
func (e *Evaluator) Normalize(line string) ([]uint, error) {
	result := []uint{}

	for _, c := range []byte(line) {
		switch c {
		case e.AState:
			result = append(result, 0)
		case e.BState:
			result = append(result, 1)
		default:
			return nil, fmt.Errorf("%c is an invalid character in the alphabet [%c %c]", c, e.AState, e.BState)
		}
	}

	return result, nil
}
