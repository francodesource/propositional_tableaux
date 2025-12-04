package tableaux

import (
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
)

// TestBuildSemanticTableaux tests whether the function produces tableaux that leads to correct assignments.
// Since the construction of a tableaux is deterministic it proves just some properties
func TestBuildSemanticTableaux(t *testing.T) {
	tests := []struct {
		name string
		f    formula.Formula
		want []map[string]bool
	}{
		{
			"simple unsat",
			simpleUnsat,
			[]map[string]bool{},
		},
		{
			"simple tautology",
			simpleTaut,
			[]map[string]bool{
				{"P": true},
				{"P": false},
			},
		},
		{
			"unsat",
			unsat,
			[]map[string]bool{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tab := BuildSemanticTableaux(tt.f)
			got := tab.Eval()

			// for now, I just check the length of the assignments
			if len(got) != len(tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
				return
			}
		})
	}
}
