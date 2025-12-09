package tableaux

import (
	"encoding/json"
	"math/rand"
	"propositional_tableaux/formula"
	"reflect"
	"slices"
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

// TestBuildAnalyticTableaux checks that the assignments discovered by the buffered analytic tableaux are the same of the semantic one
func TestBuildBufferedTableaux(t *testing.T) {
	f := func(f formula.Formula) bool {
		semanticTab := BuildSemanticTableaux(f)
		analyticTab := BuildBufferedTableaux(f)

		sAssignments := semanticTab.Eval()
		bAssignments := analyticTab.Eval()

		if len(sAssignments) != len(bAssignments) {
			t.Errorf("fail on formula: %v", f)
			return false
		}

		slice1 := make([]string, len(sAssignments))
		slice2 := make([]string, len(bAssignments))

		for i := range sAssignments {
			bytes1, _ := json.Marshal(sAssignments[i])
			bytes2, _ := json.Marshal(bAssignments[i])

			slice1[i] = string(bytes1)
			slice2[i] = string(bytes2)
		}
		slices.Sort(slice1)
		slices.Sort(slice2)
		res := slices.Compare(slice1, slice2) == 0
		if !res {
			t.Errorf("fail on formula: %v", f)
		}
		return res
	}
	maxSize := 50
	config := &quick.Config{
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = reflect.ValueOf(formula.GenerateRandom(r, r.Intn(maxSize)))
		},
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("%v", err)
	}
}

// TestBuildAnalyticTableaux checks that the assignments discovered by the buffered analytic tableaux are the same of the analytic one
func TestBuildBufferedTableaux2(t *testing.T) {
	f := func(f formula.Formula) bool {
		semanticTab := BuildAnalyticTableaux(f)
		analyticTab := BuildBufferedTableaux(f)

		sAssignments := semanticTab.Eval()
		bAssignments := analyticTab.Eval()

		if len(sAssignments) != len(bAssignments) {
			t.Errorf("fail on formula: %v", f)
			return false
		}

		slice1 := make([]string, len(sAssignments))
		slice2 := make([]string, len(bAssignments))

		for i := range sAssignments {
			bytes1, _ := json.Marshal(sAssignments[i])
			bytes2, _ := json.Marshal(bAssignments[i])

			slice1[i] = string(bytes1)
			slice2[i] = string(bytes2)
		}
		slices.Sort(slice1)
		slices.Sort(slice2)
		res := slices.Compare(slice1, slice2) == 0
		if !res {
			t.Errorf("fail on formula: %v", f)
		}
		return res
	}
	maxSize := 50
	config := &quick.Config{
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = reflect.ValueOf(formula.GenerateRandom(r, r.Intn(maxSize)))
		},
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("%v", err)
	}
}
