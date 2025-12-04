package formula

import (
	"testing"
)

// letters is a struct that holds commonly used literals for testing.
var letters = struct {
	p  Letter
	p1 Letter
	p2 Letter
	p3 Letter

	q  Letter
	q1 Letter
	q2 Letter
	q3 Letter

	r  Letter
	r1 Letter
	r2 Letter
	r3 Letter

	s  Letter
	s1 Letter
	s2 Letter
	s3 Letter

	t  Letter
	t1 Letter
	t2 Letter
	t3 Letter
}{
	p:  NewLetter("p"),
	p1: NewLetter("p1"),
	p2: NewLetter("p2"),
	p3: NewLetter("p3"),

	q:  NewLetter("q"),
	q1: NewLetter("q1"),
	q2: NewLetter("q2"),
	q3: NewLetter("q3"),

	r:  NewLetter("r"),
	r1: NewLetter("r1"),
	r2: NewLetter("r2"),
	r3: NewLetter("r3"),

	s:  NewLetter("s"),
	s1: NewLetter("s1"),
	s2: NewLetter("s2"),
	s3: NewLetter("s3"),

	t:  NewLetter("t"),
	t1: NewLetter("t1"),
	t2: NewLetter("t2"),
	t3: NewLetter("t3"),
}

func TestLetter_Name(t *testing.T) {
	name := "p"

	letter := NewLetter(name)

	if letter.Name() != name {
		t.Errorf("letter.Name() = %s, wanted %s", letter.Name(), name)
	}
}

func TestFormula_String(t *testing.T) {
	tests := []struct {
		name string
		f    Formula
		want string
	}{
		{
			name: "literal",
			f:    NewLetter("p"),
			want: "p",
		},
		{
			name: "negated literal",
			f:    NewNot(NewLetter("q")),
			want: "!q",
		},
		{
			name: "double negated literal",
			f:    NewNot(NewNot(NewLetter("r"))),
			want: "!!r",
		},
		{
			name: "and",
			f:    NewBinary(letters.p, letters.q, And),
			want: "(p & q)",
		},
		{
			name: "complex formula 1",
			f: NewOr(
				NewNot(NewAnd(letters.p1, letters.p2)),
				NewImplies(letters.q1, letters.q2),
			),
			want: "(!(p1 & p2) | (q1 -> q2))",
		},

		{
			name: "complex formula 2",
			f: NewNot(NewBiconditional(
				NewNot(NewXor(letters.p1, letters.p2)),
				NewNand(letters.q1, NewNor(letters.p2, letters.q1)),
			)),
			want: "!(!(p1 ^ p2) <-> (q1 !& (p2 !| q1)))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.String(); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormula_Class(t *testing.T) {
	// A and B are just some random formulas, used for testing a more generic kind of formulas instead of only literals
	// and operators.
	A := NewNot(NewOr(letters.t, letters.r))
	B := NewBiconditional(NewImplies(letters.p, letters.q), letters.r)

	tests := []struct {
		name     string
		formula  Formula
		expected Classification
	}{
		// Literals
		{
			name:     "letter",
			formula:  letters.p,
			expected: Literal,
		},
		{
			name:     "literal",
			formula:  NewNot(letters.p),
			expected: Literal,
		},
		// Alpha-formulas
		{
			name:     "double negation",
			formula:  NewNot(NewNot(letters.p)),
			expected: Alpha,
		},
		{
			name:     "or negation",
			formula:  NewNot(NewOr(A, letters.p)),
			expected: Alpha,
		},
		{
			name:     "implication negation",
			formula:  NewNot(NewImplies(A, B)),
			expected: Alpha,
		},
		{
			name:     "nand negation",
			formula:  NewNot(NewNand(letters.p, letters.p1)),
			expected: Alpha,
		},
		{
			name:     "nor",
			formula:  NewNor(A, B),
			expected: Alpha,
		},
		{
			name:     "biconditional",
			formula:  NewBiconditional(A, B),
			expected: Alpha,
		},
		{
			name:     "xor negation",
			formula:  NewNot(NewXor(A, B)),
			expected: Alpha,
		},
		// Beta-formulas
		{
			name:     "and negation",
			formula:  NewNot(NewAnd(A, letters.p)),
			expected: Beta,
		},
		{
			name:     "or",
			formula:  NewOr(A, B),
			expected: Beta,
		},
		{
			name:     "implication",
			formula:  NewImplies(letters.p, letters.p1),
			expected: Beta,
		},
		{
			name:     "nand",
			formula:  NewNand(letters.p, letters.p1),
			expected: Beta,
		},
		{
			name:     "nor negation",
			formula:  NewNot(NewNor(A, B)),
			expected: Beta,
		},
		{
			name:     "biconditional negation",
			formula:  NewNot(NewBiconditional(A, B)),
			expected: Beta,
		},
		{
			name:     "xor",
			formula:  NewXor(A, B),
			expected: Beta,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.formula.Class()

			if res != tt.expected {
				t.Errorf("got %v want %v", res, tt.expected)
			}
		})
	}
}

func TestBinary_Getters(t *testing.T) {
	and := NewAnd(letters.p, letters.q)

	if and.Op() != And {
		t.Errorf("expected op = %v, got %v", And, and.Op())
	}

	if and.Left() != letters.p {
		t.Errorf("expected left = %v, got %v", letters.p, and.Left())
	}

	if and.Right() != letters.q {
		t.Errorf("expected right = %v, got %v", letters.q, and.Right())
	}
}

func TestNot_Negated(t *testing.T) {
	neg := NewNot(letters.p)

	if neg.Negated() != letters.p {
		t.Errorf("Negated() = %v, expected %v", neg.Negated(), letters.p)
	}
}
