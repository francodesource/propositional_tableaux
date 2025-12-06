package tsets

import (
	"maps"
	"propositional_tableaux/formula"
	tu "propositional_tableaux/internal/testutil"
	"testing"
)

func TestTSet_Add(t *testing.T) {
	set := NewTSet()
	set.Add(tu.P, formula.NewAnd(tu.Q, tu.R))

	if set.Len() != 2 {
		t.Errorf("Expected 1 literal, got %d", set.literals.Len())
	}
}

func TestTSet_HasOnlyLiterals(t *testing.T) {
	set := NewTSet()
	set.Add(tu.P, formula.Complement(tu.Q))

	if !set.HasOnlyLiterals() {
		t.Errorf("Expected HasOnlyLiterals to return true, got false")
	}

	set.Add(formula.NewAnd(tu.P, tu.Q))

	if set.HasOnlyLiterals() {
		t.Errorf("Expected HasOnlyLiterals to return false, got true")
	}
}

func TestTSet_HasComplementaryLiterals(t *testing.T) {
	set := NewTSet()
	set.Add(tu.P, formula.Complement(tu.P))

	if !set.HasComplementaryLiterals() {
		t.Errorf("Expected HasComplementaryLiterals to return true, got false")
	}

	set = NewTSet()
	set.Add(tu.P)
	set.Add(tu.Q)

	if set.HasComplementaryLiterals() {
		t.Errorf("Expected HasComplementaryLiterals to return false, got true")
	}
}

func TestTSet_HasAlpha(t *testing.T) {
	set := NewTSet()

	if set.HasAlpha() {
		t.Errorf("Expected HasAlpha to return false, got true")
	}

	set.Add(formula.NewAnd(tu.P, tu.Q))

	if !set.HasAlpha() {
		t.Errorf("Expected HasAlpha to return true, got false")
	}
}

func TestTSet_HasBeta(t *testing.T) {
	set := NewTSet()

	if set.HasBeta() {
		t.Errorf("Expected HasBeta to return false, got true")
	}

	set.Add(formula.NewOr(tu.P, tu.Q))

	if !set.HasBeta() {
		t.Errorf("Expected HasBeta to return true, got false")
	}
}

func TestTSet_IsEmpty(t *testing.T) {
	set := NewTSet()

	if !set.IsEmpty() {
		t.Errorf("Expected IsEmpty to return true, got false")
	}
	set.Add(tu.P)

	if set.IsEmpty() {
		t.Errorf("Expected IsEmpty to return false, got true")
	}
}

func TestTSet_IterLiterals(t *testing.T) {
	set := NewTSet()
	set.Add(tu.P, tu.T, tu.R)
	set.Add(formula.NewNand(tu.S, tu.Q))

	res := map[formula.Formula]bool{}
	for f := range set.IterLiterals() {
		res[f] = true
	}

	if maps.Equal(res, map[formula.Formula]bool{tu.P: true, tu.T: true, tu.R: true}) == false {
		t.Errorf("Expected IterLiterals to return P, T, R; got %v", res)
	}
}

func TestTSet_IterAlpha(t *testing.T) {
	set := NewTSet()
	set.Add(formula.NewAnd(tu.P, tu.Q))
	set.Add(formula.NewNot(formula.NewNand(tu.R, tu.S)))
	set.Add(tu.T)

	res := map[formula.Formula]bool{}
	for f := range set.IterAlpha() {
		res[f] = true
	}

	expected := map[formula.Formula]bool{
		formula.NewAnd(tu.P, tu.Q):                  true,
		formula.NewNot(formula.NewNand(tu.R, tu.S)): true,
	}

	if maps.Equal(res, expected) == false {
		t.Errorf("Expected IterAlpha to return And(P, Q) and Nand(R, S); got %v", res)
	}
}

func TestTSet_IterBeta(t *testing.T) {
	set := NewTSet()
	set.Add(formula.NewOr(tu.P, tu.Q))
	set.Add(formula.NewNot(formula.NewNor(tu.R, tu.S)))
	set.Add(tu.T)

	res := map[formula.Formula]bool{}
	for f := range set.IterBeta() {
		res[f] = true
	}

	expected := map[formula.Formula]bool{
		formula.NewOr(tu.P, tu.Q):                  true,
		formula.NewNot(formula.NewNor(tu.R, tu.S)): true,
	}

	if maps.Equal(res, expected) == false {
		t.Errorf("Expected IterBeta to return Or(P, Q) and Nor(R, S); got %v", res)
	}
}

// TestTSet_Add2 tests adding function and its return value.
func TestTSet_Add2(t *testing.T) {
	set1 := NewTSet()
	set1.Add(tu.P, tu.P1, tu.P2)

	set2 := NewTSet()
	set2.Add(tu.P, formula.NewAnd(tu.P1, tu.P2))

	set3 := NewTSet()
	set3.Add(tu.Q, formula.NewNot(tu.R), formula.NewNot(formula.NewAnd(tu.P1, tu.P2)))

	set4 := NewTSet()
	set4.Add(formula.NewNot(formula.NewOr(tu.Q, tu.P)))

	tests := []struct {
		name  string
		set   TSet
		value []formula.Formula
		want  bool
	}{
		{
			name:  "empty",
			set:   NewTSet(),
			value: []formula.Formula{tu.P},
			want:  false,
		},
		{
			name:  "no complementary pair",
			set:   set1,
			value: []formula.Formula{tu.Q},
			want:  false,
		},
		{
			name:  "complementary pair",
			set:   set1,
			value: []formula.Formula{formula.NewNot(tu.P)},
			want:  true,
		},
		{
			name:  "complementary alpha formula",
			set:   set2,
			value: []formula.Formula{formula.Complement(formula.NewAnd(tu.P1, tu.P2))},
			want:  true,
		},
		{
			name:  "complementary negated alpha formula and other",
			set:   set2,
			value: []formula.Formula{formula.NewAnd(tu.P1, tu.P2), tu.T},
			want:  true,
		},
		{
			name:  "complementary negated alpha formula and other",
			set:   set3,
			value: []formula.Formula{formula.NewAnd(tu.P1, tu.P2), tu.T},
			want:  true,
		},
		{
			name:  "complementary beta formula and other",
			set:   set4,
			value: []formula.Formula{formula.NewOr(tu.Q, tu.P), tu.T},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.set.Add(tt.value...)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
