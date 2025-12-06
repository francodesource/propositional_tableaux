package tableaux

import (
	"maps"
	"propositional_tableaux/formula"
	tu "propositional_tableaux/internal/testutil"
	"testing"
)

var (
	simpleUnsat = formula.NewAnd(
		tu.P,
		formula.Complement(tu.P),
	)

	simpleTaut = formula.NewOr(
		tu.P,
		formula.Complement(tu.P),
	)

	unsat = formula.NewAnd(
		formula.NewOr(tu.P, tu.Q),
		formula.NewAnd(formula.Complement(tu.P), formula.Complement(tu.Q)),
	)

	sat = formula.NewAnd(
		tu.P,
		formula.NewOr(formula.NewNot(tu.Q), formula.NewNot(tu.P)),
	)
)

// TestBuildSemanticTableaux tests whether the function produces tableaux that leads to correct assignments.
// Since the construction of a tableaux is deterministic it proves just some properties
func TestBuildSemanticTableaux(t *testing.T) {
	tests := []struct {
		name string
		f    formula.Formula
		want []Assignment
	}{
		{
			"simple unsat",
			simpleUnsat,
			[]Assignment{},
		},
		{
			"simple tautology",
			simpleTaut,
			[]Assignment{
				{"P": true},
				{"P": false},
			},
		},
		{
			"unsat",
			unsat,
			[]Assignment{},
		},
		{
			"sat",
			sat,
			[]Assignment{{
				"P": true,
				"Q": false,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tab := BuildSemanticTableaux(tt.f)
			got := tab.Eval()

			if len(got) != len(tt.want) {
				t.Errorf("want %v, got %v", tt.want, got)
			}

			equalCounter := 0
			for _, assignment1 := range tt.want {
				for _, assignment2 := range got {
					if maps.Equal(assignment2, assignment1) {
						equalCounter++
					}
				}
			}

			if equalCounter != len(tt.want) {
				t.Errorf("want %v, got %v", tt.want, got)
			}
		})
	}
}

// TestBuildSemanticTableaux2 checks the improvement of checking for complementary formula.
func TestBuildSemanticTableaux2(t *testing.T) {
	inner := formula.NewNand(
		formula.NewOr(tu.P, formula.NewBiconditional(tu.Q, tu.R)),
		formula.NewBiconditional(tu.P1, formula.NewXor(tu.R, tu.S)))

	f := formula.NewAnd(inner, formula.NewNot(inner))
	tab := BuildSemanticTableaux(f)
	assignment := tab.Eval()

	if len(assignment) > 0 {
		t.Errorf("expected formula to be unasat")
	}

	if tab.Height() != 2 {
		t.Errorf("expected tableaux to be of minimal length 2")
	}
}
