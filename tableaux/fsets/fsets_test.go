package fsets

import (
	"propositional_tableaux/formula"
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
