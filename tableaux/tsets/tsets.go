package tsets

import (
	"fmt"
	"github.com/francodesource/propositional_tableaux/formula"
	"github.com/francodesource/propositional_tableaux/tableaux/fsets"
	"iter"
)

// TSet represents a set specialized for tableaux formulas sets.
type TSet struct {
	literals      fsets.FormulaSet
	alphaFormulas fsets.FormulaSet
	betaFormulas  fsets.FormulaSet
}

// NewTSet creates and returns a new empty TSet.
func NewTSet() TSet {
	return TSet{
		literals:      fsets.New(),
		alphaFormulas: fsets.New(),
		betaFormulas:  fsets.New(),
	}
}

func asByte(b bool) byte {
	if b {
		return 1
	}

	return 0
}

// HasComplementOf returns true if the set contains the complement of the given formula f.
func (s TSet) HasComplementOf(f formula.Formula) bool {
	switch f.Class() {
	case formula.LiteralClass:
		return s.literals.HasComplementaryOf(f)
	case formula.Alpha:
		return s.betaFormulas.HasComplementaryOf(f)
	case formula.Beta:
		return s.alphaFormulas.HasComplementaryOf(f)
	}
	return false
}

// addOne add an element to the set. If the element is already contained, it will be overwritten.
// If f is nil, it will not be added.
// Returns 1 if this set contains the complement of f, 0 otherwise.
func (s TSet) addOne(f formula.Formula) byte {
	if f != nil {
		switch f.Class() {
		case formula.LiteralClass:

			s.literals.Add(f)
			return asByte(s.literals.HasComplementaryOf(f))
		case formula.Alpha:
			s.alphaFormulas.Add(f)
			// if a formula is an alpha formula, its complement is a beta formula and vice versa.
			return asByte(s.betaFormulas.HasComplementaryOf(f))
		case formula.Beta:
			s.betaFormulas.Add(f)
			return asByte(s.alphaFormulas.HasComplementaryOf(f))
		}
	}

	return 0
}

// Add adds all the passed elements to the set s and returns true if s contains the complement of at least one element,
// false otherwise.
func (s TSet) Add(fs ...formula.Formula) bool {
	var flag byte = 0
	for _, f := range fs {
		flag = flag | s.addOne(f) // using bitwise operator to avoid lazy evaluation.
	}
	return flag == 1
}

// Len returns the total number of elements in the set.
func (s TSet) Len() int {
	return s.literals.Len() + s.alphaFormulas.Len() + s.betaFormulas.Len()
}

// IsEmpty returns true if the set is empty, false otherwise.
func (s TSet) IsEmpty() bool {
	return s.Len() == 0
}

// HasOnlyLiterals returns true if the set contains only literals, false otherwise.
func (s TSet) HasOnlyLiterals() bool {
	return s.literals.Len() > 0 && s.alphaFormulas.Len() == 0 && s.betaFormulas.Len() == 0
}

// HasComplementaryLiterals returns true if the set contains at least a pair of complementary literals, false otherwise.
func (s TSet) HasComplementaryLiterals() bool {
	for literal := range s.literals.Iter() {
		if s.literals.HasComplementaryOf(literal) {
			return true
		}
	}
	return false
}

// HasAlpha returns true if the set contains at least one alpha formula, false otherwise.
func (s TSet) HasAlpha() bool {
	return s.alphaFormulas.Len() > 0
}

// HasBeta returns true if the set contains at least one beta formula, false otherwise.
func (s TSet) HasBeta() bool {
	return s.betaFormulas.Len() > 0
}

// IterLiterals returns an iterator over the literals in the set.
func (s TSet) IterLiterals() iter.Seq[formula.Formula] {
	return s.literals.Iter()
}

// IterAlpha returns an iterator over the alpha formulas in the set.
func (s TSet) IterAlpha() iter.Seq[formula.Formula] {
	return s.alphaFormulas.Iter()
}

// IterBeta returns an iterator over the beta formulas in the set.
func (s TSet) IterBeta() iter.Seq[formula.Formula] {
	return s.betaFormulas.Iter()
}

func (s TSet) String() string {
	return fmt.Sprintf("{ literals: %s, alpha: %s, beta: %s }",
		s.literals.String(), s.alphaFormulas.String(), s.betaFormulas.String())
}

// RemoveAlpha returns a copy of the given set without the given alpha formula.
func RemoveAlpha(set TSet, f formula.Formula) TSet {
	alphaSet := fsets.Remove(set.alphaFormulas, f)
	return TSet{
		literals:      fsets.Clone(set.literals),
		alphaFormulas: alphaSet,
		betaFormulas:  fsets.Clone(set.betaFormulas),
	}
}

// RemoveBeta returns a copy of the given set without the given beta formula.
func RemoveBeta(set TSet, f formula.Formula) TSet {
	betaSet := fsets.Remove(set.betaFormulas, f)
	return TSet{
		literals:      fsets.Clone(set.literals),
		alphaFormulas: fsets.Clone(set.alphaFormulas),
		betaFormulas:  betaSet,
	}
}
