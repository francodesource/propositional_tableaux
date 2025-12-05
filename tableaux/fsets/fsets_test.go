package fsets

import (
	"maps"
	"propositional_tableaux/formula"
	tu "propositional_tableaux/internal/testutil"
	"testing"
)

func TestFormulaSet_String(t *testing.T) {
	tests := []struct {
		name string
		set  FormulaSet
		want string
	}{
		{
			"empty",
			New(),
			"{}",
		},
		{
			"singleton",
			New(formula.NewLetter("p")),
			"{p}",
		},
		{
			"formula",
			New(
				formula.NewAnd(
					formula.NewLetter("p"),
					formula.NewLetter("q")),
				formula.NewNot(formula.NewAnd(
					formula.NewLetter("p"),
					formula.NewLetter("q")),
				),
			),

			"{(p & q), !(p & q)}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != tt.set.String() {
				t.Errorf("got %v, want %v", tt.set.String(), tt.want)
			}
		})
	}
}

func TestFormulaSet_HasComplementaryOf(t *testing.T) {
	tests := []struct {
		name string
		f    formula.Formula
		set  FormulaSet
		want bool
	}{
		{
			"literals complementary pair",
			formula.NewLetter("p"),
			New(formula.NewLetter("p"), formula.NewNot(formula.NewLetter("p"))),
			true,
		},
		{
			"single letter",
			formula.NewLetter("p"),
			New(formula.NewLetter("p")),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.set.HasComplementaryOf(tt.f)

			if res != tt.want {
				t.Errorf("got %v, want %v", res, tt.want)
			}
		})
	}
}

func TestFormulaSet_Add(t *testing.T) {
	sampleSet := func() FormulaSet {
		return New(
			formula.NewLetter("p"),
			formula.NewNot(formula.NewLetter("q")),
			formula.NewNand(formula.NewLetter("r"), formula.NewLetter("s")),
		)
	}

	tests := []struct {
		name  string
		toAdd []formula.Formula
		set   FormulaSet
		want  FormulaSet
	}{
		{
			"nothing to add",
			[]formula.Formula{},
			New(),
			New(),
		},
		{
			"add one",
			[]formula.Formula{formula.NewLetter("a")},
			sampleSet(),
			New(
				formula.NewLetter("p"),
				formula.NewNot(formula.NewLetter("q")),
				formula.NewNand(formula.NewLetter("r"), formula.NewLetter("s")),
				formula.NewLetter("a"),
			),
		},
		{
			"add one already present",
			[]formula.Formula{formula.NewLetter("p")},
			sampleSet(),
			New(
				formula.NewLetter("p"),
				formula.NewNot(formula.NewLetter("q")),
				formula.NewNand(formula.NewLetter("r"), formula.NewLetter("s")),
			),
		},
		{
			"add two, one already present",
			[]formula.Formula{formula.NewLetter("a"), formula.NewLetter("p")},
			sampleSet(),
			New(
				formula.NewLetter("p"),
				formula.NewNot(formula.NewLetter("q")),
				formula.NewNand(formula.NewLetter("r"), formula.NewLetter("s")),
				formula.NewLetter("a"),
			),
		},
		{
			"add nil",
			[]formula.Formula{nil},
			sampleSet(),
			New(
				formula.NewLetter("p"),
				formula.NewNot(formula.NewLetter("q")),
				formula.NewNand(formula.NewLetter("r"), formula.NewLetter("s")),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.set.Add(tt.toAdd...)

			if !maps.Equal(res.values, tt.want.values) {
				t.Errorf("got %v, want %v", res, tt.want)
			}
		})
	}
}

func TestFormulaSet_Len(t *testing.T) {
	tests := []struct {
		name string
		set  FormulaSet
		want int
	}{
		{
			"empty",
			New(),
			0,
		},
		{
			"one",
			New(formula.NewLetter("p")),
			1,
		},
		{
			"three",
			New(
				tu.P,
				formula.NewNot(tu.Q),
				formula.NewNand(tu.R, tu.S),
			),
			3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.set.Len()

			if res != tt.want {
				t.Errorf("got %v, want %v", res, tt.want)
			}
		})
	}
}
