package tableaux

import (
	"math/rand"
	"propositional_tableaux/formula"
	"reflect"
	"testing"
	"testing/quick"
)

func TestBufferSet_Has(t *testing.T) {
	tests := []struct {
		name   string
		values BufferSet
		check  formula.Formula
		want   bool
	}{
		{
			"has first value",
			NewBufferSet(formula.NewLetter("p")),
			formula.NewLetter("p"),
			true,
		},
		{
			"has second value",
			NewBufferSet(formula.NewLetter("p"), formula.NewLetter("q")),
			formula.NewLetter("q"),
			true,
		},
		{
			"does not have value",
			NewBufferSet(formula.NewLetter("p"), formula.NewLetter("q")),
			formula.NewLetter("r"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.values.Has(tt.check); got != tt.want {
				t.Errorf("BufferSet.Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBufferSet_Add(t *testing.T) {
	tests := []struct {
		name     string
		initial  BufferSet
		toAdd    formula.Formula
		expected BufferSet
	}{
		{
			"add to empty buffer",
			NewBufferSet(),
			formula.NewLetter("p"),
			NewBufferSet(formula.NewLetter("p")),
		},
		{
			"add second value",
			NewBufferSet(formula.NewLetter("p")),
			formula.NewLetter("q"),
			NewBufferSet(formula.NewLetter("p"), formula.NewLetter("q")),
		},
		{
			"add existing first value",
			NewBufferSet(formula.NewLetter("p")),
			formula.NewLetter("p"),
			NewBufferSet(formula.NewLetter("p")),
		},
		{
			"add existing second value",
			NewBufferSet(formula.NewLetter("p"), formula.NewLetter("q")),
			formula.NewLetter("q"),
			NewBufferSet(formula.NewLetter("p"), formula.NewLetter("q")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initial.Add(tt.toAdd)
			if tt.initial[0] != tt.expected[0] || tt.initial[1] != tt.expected[1] {
				t.Errorf("BufferSet.Add() = %v, want %v", tt.initial, tt.expected)
			}
		})
	}
}

func TestBufferSet_HasComplementOf(t *testing.T) {
	tests := []struct {
		name   string
		values BufferSet
		check  formula.Formula
		want   bool
	}{
		{
			"has complement of first value",
			NewBufferSet(formula.NewLetter("p")),
			formula.NewNot(formula.NewLetter("p")),
			true,
		},
		{
			"has complement of second value",
			NewBufferSet(formula.NewLetter("p"), formula.NewNot(formula.NewLetter("q"))),
			formula.NewLetter("q"),
			true,
		},
		{
			"does not have complement",
			NewBufferSet(formula.NewLetter("p"), formula.NewLetter("q")),
			formula.NewLetter("r"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.values.HasComplementOf(tt.check); got != tt.want {
				t.Errorf("BufferSet.HasComplementOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBuildBufferTableaux_SemanticCompare checks that the assignments discovered by the buffered analytic tableaux are the same of the semantic one
func TestBuildBufferTableaux_SemanticCompare(t *testing.T) {
	f := func(f formula.Formula) bool {
		semanticTab := BuildSemanticTableaux(f)
		analyticTab := BuildBufferTableaux(f)

		sAssignments := semanticTab.Eval()
		bAssignments := analyticTab.Eval()

		res := compareAssignments(sAssignments, bAssignments)
		if !res {
			t.Errorf("fail on formula: %v", f)
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
		t.Errorf("%v", err)
	}
}

// TestBuildBufferTableaux_AnalyticCompare checks that the assignments discovered by the buffered analytic tableaux are the same of the analytic one
func TestBuildBufferTableaux_AnalyticCompare(t *testing.T) {
	f := func(f formula.Formula) bool {
		semanticTab := BuildAnalyticTableaux(f)
		analyticTab := BuildBufferTableaux(f)

		sAssignments := semanticTab.Eval()
		bAssignments := analyticTab.Eval()

		res := compareAssignments(sAssignments, bAssignments)

		if !res {
			t.Errorf("fail on formula: %v", f)
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
		t.Errorf("%v", err)
	}
}

func TestBufferTableauxMarks(t *testing.T) {
	f := func(f formula.Formula) bool {
		tab := BuildBufferTableaux(f)

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
