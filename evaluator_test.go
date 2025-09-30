package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestEvaluatorNormalize(t *testing.T) {
	tests := []struct {
		eval      Evaluator
		line      string
		want      []uint
		wantError bool
		errMsg    string
	}{
		{
			eval: Evaluator{AState: 'a', BState: 'b'},
			line: "ababb",
			want: []uint{0, 1, 0, 1, 1},
		},
		{
			eval: Evaluator{AState: '0', BState: '1'},
			line: "000111",
			want: []uint{0, 0, 0, 1, 1, 1},
		},
		{
			eval:      Evaluator{AState: 'a', BState: 'b'},
			line:      "000",
			want:      nil,
			wantError: true,
			errMsg:    "0 is an invalid character in the alphabet [a b]",
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("(%c,%c) : %v", tt.eval.AState, tt.eval.BState, tt.line)
		t.Run(testname, func(t *testing.T) {
			result, err := tt.eval.Normalize(tt.line)
			if err != nil {
				if !tt.wantError {
					t.Fatalf("Evaluator.Normalize() unexpected error: %v", err)
				}

				if tt.errMsg != err.Error() {
					t.Errorf("Expected error '%v', got '%v'", tt.errMsg, err)
				}
			}

			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("%v != %v", tt.want, result)
			}
		})
	}
}

func TestNewEvaluator(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantA     byte
		wantB     byte
		wantError bool
		errMsg    string
	}{
		{
			name:      "valid input with a,b alphabet",
			input:     "a,b\n+,A,B,A\n-,B,B,B",
			wantA:     'a',
			wantB:     'b',
			wantError: false,
		},
		{
			name:      "valid input with 0,1 alphabet",
			input:     "0,1\n+,A,B,A\n-,B,B,B",
			wantA:     '0',
			wantB:     '1',
			wantError: false,
		},
		{
			name:      "no newline separator",
			input:     "a,b",
			wantError: true,
			errMsg:    "no new line was found",
		},
		{
			name:      "empty input",
			input:     "",
			wantError: true,
			errMsg:    "no new line was found",
		},
		{
			name:      "only newline",
			input:     "\n",
			wantError: true,
			errMsg:    "[] does not have two alphabet variants",
		},
		{
			name:      "single alphabet character",
			input:     "a\n+,A,A,A",
			wantError: true,
			errMsg:    "[a] does not have two alphabet variants",
		},
		{
			name:      "three alphabet characters",
			input:     "a,b,c\n+,A,A,A",
			wantError: true,
			errMsg:    "[a b c] does not have two alphabet variants",
		},
		{
			name:      "identical alphabet characters",
			input:     "a,a\n+,A,A,A",
			wantError: true,
			errMsg:    "alphabet should contain differing characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewEvaluator(tt.input)

			if tt.wantError {
				if err == nil {
					t.Fatalf("NewEvaluator() expected error but got none")
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing '%s', got '%v'", tt.errMsg, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("NewEvaluator() unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("NewEvaluator() returned nil result without error")
			}

			if result.AState != tt.wantA {
				t.Errorf("AState = %c, want %c", result.AState, tt.wantA)
			}

			if result.BState != tt.wantB {
				t.Errorf("BState = %c, want %c", result.BState, tt.wantB)
			}

			if result.Dfa == nil {
				t.Error("Dfa is nil")
			}
		})
	}
}
