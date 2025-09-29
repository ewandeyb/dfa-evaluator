package main

import (
	"fmt"
	"testing"
)

func TestNewDFA_ValidInput(t *testing.T) {
	input := `+,A,B,C
-,B,A,C
,C,A,B`

	dfa, err := NewDFA(input)
	if err != nil {
		t.Fatalf("NewDFA() unexpected error: %v", err)
	}

	// Check state A (index 0)
	if !dfa.States[0].IsValid {
		t.Error("State A should be valid")
	}
	if !dfa.States[0].IsStart {
		t.Error("State A should be a start state")
	}
	if dfa.States[0].IsEnd {
		t.Error("State A should not be an end state")
	}
	if dfa.States[0].NextA != 1 {
		t.Errorf("State A NextA = %d, want 1", dfa.States[0].NextA)
	}
	if dfa.States[0].NextB != 2 {
		t.Errorf("State A NextB = %d, want 2", dfa.States[0].NextB)
	}

	// Check state B (index 1)
	if !dfa.States[1].IsValid {
		t.Error("State B should be valid")
	}
	if dfa.States[1].IsStart {
		t.Error("State B should not be a start state")
	}
	if !dfa.States[1].IsEnd {
		t.Error("State B should be an end state")
	}

	// Check state C (index 2)
	if !dfa.States[2].IsValid {
		t.Error("State C should be valid")
	}
}

func TestNewDFA_Errors(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		errMsg string
	}{
		{
			name:   "Multiple start states",
			input:  "+,A,B,C\n+,B,A,C",
			errMsg: "multiple start states parsed",
		},
		{
			name:   "Wrong number of columns",
			input:  "+,A,B",
			errMsg: "1:[+ A B] does not have 4 columns",
		},
		{
			name:   "Invalid state name",
			input:  "+,a,B,C",
			errMsg: "invalid state name `a`",
		},
		{
			name:   "Invalid NextA state",
			input:  "+,A,5,C",
			errMsg: "invalid state name `5`",
		},
		{
			name:   "Invalid NextB state",
			input:  "+,A,B,@",
			errMsg: "invalid state name `@`",
		},
		{
			name:   "Too many columns",
			input:  "+,A,B,C,D",
			errMsg: "1:[+ A B C D] does not have 4 columns",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDFA(tt.input)
			if err == nil {
				t.Errorf("NewDFA() error = %v", err)
			}

			if err.Error() != tt.errMsg {
				t.Errorf("Expected message '%s' != '%s'", tt.errMsg, err.Error())
			}
		})
	}
}

func TestDFA_Move(t *testing.T) {
	input := `+,A,B,C
-,B,A,C
,C,A,B`

	dfa, err := NewDFA(input)
	if err != nil {
		t.Fatalf("NewDFA() unexpected error: %v", err)
	}

	// Start at state A (index 0)
	dfa.CurrentState = 0

	// Move with input 0 (should go to NextA = B)
	err = dfa.Move(0)
	if err != nil {
		t.Errorf("Move(0) unexpected error: %v", err)
	}
	if dfa.CurrentState != 1 {
		t.Errorf("After Move(0), CurrentState = %d, want 1", dfa.CurrentState)
	}

	// Move with input 1 (should go to NextB = C)
	err = dfa.Move(1)
	if err != nil {
		t.Errorf("Move(1) unexpected error: %v", err)
	}
	if dfa.CurrentState != 2 {
		t.Errorf("After Move(1), CurrentState = %d, want 2", dfa.CurrentState)
	}
}

func TestDFA_MoveToInvalidState(t *testing.T) {
	input := `+,A,B,D`

	dfa, err := NewDFA(input)
	if err != nil {
		t.Fatalf("NewDFA() unexpected error: %v", err)
	}

	dfa.CurrentState = 0

	// Try to move to state D (index 3), which is not defined (invalid)
	err = dfa.Move(1)
	if err == nil {
		t.Error("Move() to invalid state should return error")
	}
}

func TestDFA_MoveSequence(t *testing.T) {
	// Create a simple DFA that accepts binary strings ending in 01
	input := `+,A,A,B
,B,C,B
-,C,C,B`

	dfa, err := NewDFA(input)
	if err != nil {
		t.Fatalf("NewDFA() unexpected error: %v", err)
	}

	tests := []struct {
		sequence []uint
		endState uint
	}{
		{
			sequence: []uint{1, 0},
			endState: 2, // State C
		},
		{
			sequence: []uint{0, 1, 1},
			endState: 1,
		},
		{
			sequence: []uint{1},
			endState: 1, // State B
		},
		{
			sequence: []uint{0},
			endState: 0, // State A
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.sequence)
		t.Run(testname, func(t *testing.T) {
			dfa.CurrentState = 0 // Reset to start state

			var err error
			for _, input := range tt.sequence {
				err = dfa.Move(input)
				if err != nil {
					break
				}
			}

			if dfa.CurrentState != tt.endState {
				t.Errorf("After sequence, CurrentState = %d, want %d", dfa.CurrentState, tt.endState)
			}
		})
	}
}

func TestDFA_EdgeCases(t *testing.T) {
	t.Run("Single state DFA", func(t *testing.T) {
		input := `+,A,A,A`
		dfa, err := NewDFA(input)
		if err != nil {
			t.Fatalf("NewDFA() unexpected error: %v", err)
		}

		// Should loop back to itself
		dfa.CurrentState = 0
		err = dfa.Move(0)
		if err != nil {
			t.Errorf("Move(0) unexpected error: %v", err)
		}
		if dfa.CurrentState != 0 {
			t.Errorf("CurrentState = %d, want 0", dfa.CurrentState)
		}
	})

	t.Run("State Z (last valid state)", func(t *testing.T) {
		input := `+,Z,A,B`
		dfa, err := NewDFA(input)
		if err != nil {
			t.Fatalf("NewDFA() unexpected error: %v", err)
		}

		if !dfa.States[25].IsValid {
			t.Error("State Z should be valid")
		}
	})

	t.Run("Empty state type", func(t *testing.T) {
		input := `,A,B,C`
		dfa, err := NewDFA(input)
		if err != nil {
			t.Fatalf("NewDFA() unexpected error: %v", err)
		}

		if dfa.States[0].IsStart || dfa.States[0].IsEnd {
			t.Error("State with empty type should be neither start nor end")
		}
	})
}
