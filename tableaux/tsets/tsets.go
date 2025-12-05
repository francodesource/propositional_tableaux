package tsets

import (
	"iter"
	"propositional_tableaux/formula"
	"propositional_tableaux/tableaux/fsets"
)

// TSet represents a set specialized for tableaux formulas sets.
type TSet struct {
	literals      fsets.FormulaSet
	alphaFormulas fsets.FormulaSet
	betaFormulas  fsets.FormulaSet
}

func NewTSet() TSet {
	return TSet{
		literals:      fsets.New(),
		alphaFormulas: fsets.New(),
		betaFormulas:  fsets.New(),
	}
}

func (s TSet) Add(f formula.Formula) TSet {
	switch f.Class() {
	case formula.LiteralClass:
		s.literals = s.literals.Add(f)
	case formula.Alpha:
		s.alphaFormulas = s.alphaFormulas.Add(f)
	case formula.Beta:
		s.betaFormulas = s.betaFormulas.Add(f)
	}
	return s
}

func (s TSet) Len() int {
	return s.literals.Len() + s.alphaFormulas.Len() + s.betaFormulas.Len()
}

func (s TSet) IsEmpty() bool {
	return s.Len() == 0
}

func (s TSet) HasOnlyLiterals() bool {
	return s.literals.Len() > 0 && s.alphaFormulas.Len() == 0 && s.betaFormulas.Len() == 0
}

func (s TSet) HasComplementaryLiterals() bool {
	for literal := range s.literals.Iter() {
		if s.literals.HasComplementaryOf(literal) {
			return true
		}
	}
	return false
}

func (s TSet) HasAlpha() bool {
	return s.alphaFormulas.Len() > 0
}

func (s TSet) HasBeta() bool {
	return s.betaFormulas.Len() > 0
}

func (s TSet) IterLiterals() iter.Seq[formula.Formula] {
	return s.literals.Iter()
}

func (s TSet) IterAlpha() iter.Seq[formula.Formula] {
	return s.alphaFormulas.Iter()
}

func (s TSet) IterBeta() iter.Seq[formula.Formula] {
	return s.betaFormulas.Iter()
}

func RemoveAlpha(set TSet, f formula.Formula) TSet {
	alphaSet := fsets.Remove(set.alphaFormulas, f)
	return TSet{
		literals:      fsets.Clone(set.literals),
		alphaFormulas: alphaSet,
		betaFormulas:  fsets.Clone(set.literals),
	}
}

func RemoveBeta(set TSet, f formula.Formula) TSet {
	betaSet := fsets.Remove(set.betaFormulas, f)
	return TSet{
		literals:      fsets.Clone(set.literals),
		alphaFormulas: fsets.Clone(set.literals),
		betaFormulas:  betaSet,
	}
}
