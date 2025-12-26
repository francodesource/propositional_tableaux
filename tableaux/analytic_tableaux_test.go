package tableaux

import (
	"fmt"
	"github.com/francodesource/propositional_tableaux/formula"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

// TestBuildAnalyticTableaux_SemanticCompare checks that the assignments discovered by the analytic tableaux are the same of the semantic one
func TestBuildAnalyticTableaux_SemanticCompare(t *testing.T) {
	f := func(f formula.Formula) bool {
		semanticTab := BuildSemanticTableaux(f)
		analyticTab := BuildAnalyticTableaux(f)

		sAssignments := semanticTab.Eval()
		aAssignments := analyticTab.Eval()

		res := compareAssignments(sAssignments, aAssignments)

		if !res {
			t.Errorf("Failed on %v\n", f)
		}
		return res
	}
	maxSize := FormulaMaxSize
	config := &quick.Config{Values: func(values []reflect.Value, r *rand.Rand) {
		values[0] = reflect.ValueOf(formula.GenerateRandom(r, r.Intn(maxSize)))
	},
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("%v", err)
	}
}

// TestBuildAnalyticTableaux_Assignments checks that the assignments obtained by an analytic tableaux, satisfy the formula.
func TestBuildAnalyticTableaux_Assignments(t *testing.T) {
	f := func(f formula.Formula) bool {
		tab := BuildAnalyticTableaux(f)
		assignments := tab.Eval()

		for _, a := range assignments {
			if !evaluate(f, a) {
				fmt.Printf("test fail for formula %v with assignment %v", f, a)
				return false
			}
		}

		return true
	}

	maxSize := FormulaMaxSize

	config := &quick.Config{
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = reflect.ValueOf(formula.GenerateRandom(r, r.Intn(maxSize)))
		},
	}

	if err := quick.Check(f, config); err != nil {
		t.Error(err)
	}
}

// TestBuildAnalyticTableaux_TruthTables compares analytic tableaux results with the one of the truth tables
func TestBuildAnalyticTableaux_TruthTables(t *testing.T) {
	f := func(f formula.Formula) bool {
		tab := BuildAnalyticTableaux(f)

		return compareTableauxWithTruthTables(t, f, tab)
	}

	maxSize := FormulaMaxSize
	config := &quick.Config{
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = reflect.ValueOf(formula.GenerateRandom(r, r.Intn(maxSize)))
		},
	}

	if err := quick.Check(f, config); err != nil {
		t.Error(err)
	}
}

func TestAnalyticTableauxMarks(t *testing.T) {
	f := func(f formula.Formula) bool {
		tab := BuildAnalyticTableaux(f)

		res := testTableauxMarks(t, tab)

		if !res {
			t.Errorf("Fail for %v\n", f)
		}
		return res
	}

	maxSize := FormulaMaxSize

	config := &quick.Config{
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = reflect.ValueOf(formula.GenerateRandom(r, r.Intn(maxSize)))
		},
	}

	if err := quick.Check(f, config); err != nil {
		t.Error(err)
	}
}
