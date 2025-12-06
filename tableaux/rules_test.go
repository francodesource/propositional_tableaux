package tableaux

import (
	"propositional_tableaux/formula"
	tu "propositional_tableaux/internal/testutil"
	"testing"
)

// A and B are two formulas used in tests.
var A = formula.NewNot(
	formula.NewAnd(
		formula.NewOr(tu.P, tu.P1), tu.P2),
)

var B = formula.NewOr(
	formula.NewNot(
		formula.NewOr(tu.Q, tu.Q1)),
	formula.NewAnd(tu.Q, tu.Q2),
)

func TestApplyRule(t *testing.T) {
	tests := []struct {
		name     string
		input    formula.Formula
		expected [2]formula.Formula
	}{
		// ALPHA RULES
		{
			name:     "double negation",
			input:    formula.NewNot(formula.NewNot(A)),
			expected: [2]formula.Formula{A, nil},
		},
		{
			name:     "and",
			input:    formula.NewAnd(A, B),
			expected: [2]formula.Formula{A, B},
		},
		{
			name: "negated or",
			input: formula.NewNot(
				formula.NewOr(A, B),
			),
			expected: [2]formula.Formula{
				formula.NewNot(A),
				formula.NewNot(B),
			},
		},
		{
			name: "negated implication",
			input: formula.NewNot(
				formula.NewImplies(A, B)),
			expected: [2]formula.Formula{
				A,
				formula.NewNot(B),
			},
		},
		{
			name: "negated nand",
			input: formula.NewNot(
				formula.NewNand(A, B)),
			expected: [2]formula.Formula{
				A,
				B,
			},
		},
		{
			name:  "nor",
			input: formula.NewNor(A, B),
			expected: [2]formula.Formula{
				formula.NewNot(A),
				formula.NewNot(B),
			},
		},
		{
			name:  "biconditional",
			input: formula.NewBiconditional(A, B),
			expected: [2]formula.Formula{
				formula.NewImplies(A, B),
				formula.NewImplies(B, A),
			},
		},
		{
			name:  "negated xor",
			input: formula.NewBiconditional(A, B),
			expected: [2]formula.Formula{
				formula.NewImplies(A, B),
				formula.NewImplies(B, A),
			},
		},
		// BETA RULES
		{
			name: "negated and",
			input: formula.NewNot(
				formula.NewAnd(A, B),
			),
			expected: [2]formula.Formula{
				formula.NewNot(A),
				formula.NewNot(B),
			},
		},
		{
			name:  "or",
			input: formula.NewOr(A, B),
			expected: [2]formula.Formula{
				A,
				B,
			},
		},
		{
			name:  "implication",
			input: formula.NewImplies(A, B),
			expected: [2]formula.Formula{
				formula.NewNot(A),
				B,
			},
		},
		{
			name:  "nand",
			input: formula.NewNand(A, B),
			expected: [2]formula.Formula{
				formula.NewNot(A),
				formula.NewNot(B),
			},
		},
		{
			name: "negated nor",
			input: formula.NewNot(
				formula.NewNor(A, B),
			),
			expected: [2]formula.Formula{
				A,
				B,
			},
		},
		{
			name: "negated biconditional",
			input: formula.NewNot(
				formula.NewBiconditional(A, B),
			),
			expected: [2]formula.Formula{
				formula.NewNot(formula.NewImplies(A, B)),
				formula.NewNot(formula.NewImplies(B, A)),
			},
		},
		{
			name:  "xor",
			input: formula.NewXor(A, B),
			expected: [2]formula.Formula{
				formula.NewNot(formula.NewImplies(A, B)),
				formula.NewNot(formula.NewImplies(B, A)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, right := ApplyRule(tt.input)
			if left != tt.expected[0] {
				t.Errorf("ApplyRule(%v) left = %v, want %v", tt.input, left, tt.expected[0])
			}

			if right != tt.expected[1] {
				t.Errorf("ApplyRule(%v) right = %v, want %v", tt.input, right, tt.expected[1])
			}
		})
	}
}
