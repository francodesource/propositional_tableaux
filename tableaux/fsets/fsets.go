package fsets

import (
	"github.com/francodesource/propositional_tableaux/formula"
	"iter"
	"maps"
	"strings"
)

// FormulaSet is a set of formulas without duplicate.
type FormulaSet struct {
	values map[formula.Formula]struct{}
}

// New creates a new FormulaSet containing the given formulas.
func New(formulas ...formula.Formula) FormulaSet {
	values := make(map[formula.Formula]struct{}, len(formulas))

	for _, v := range formulas {
		values[v] = struct{}{}
	}
	return FormulaSet{values: values}
}

func (s FormulaSet) String() string {
	strs := make([]string, 0, len(s.values))

	for f := range s.Iter() {
		strs = append(strs, f.String())
	}

	return "{" + strings.Join(strs, ", ") + "}"
}

// Has returns true if the formula f is contained in the set.
func (s FormulaSet) Has(f formula.Formula) bool {
	_, ok := s.values[f]
	return ok
}

// HasComplementaryOf returns true if the complement of the given literal is contained in the set.
func (s FormulaSet) HasComplementaryOf(literal formula.Formula) bool {
	return s.Has(formula.Complement(literal))
}

// Iter returns an iterator over the formulas in the set.
func (s FormulaSet) Iter() iter.Seq[formula.Formula] {
	return maps.Keys(s.values)
}

// Add adds all the passed formulas in-place and returns itself.
// If a formula is already contained it will be overwritten.
// If a formula is nil it won't be added.
func (s FormulaSet) Add(formulas ...formula.Formula) FormulaSet {
	for _, f := range formulas {
		if f != nil {
			s.values[f] = struct{}{}
		}
	}
	return s
}

// Len returns the number of formulas in the set.
func (s FormulaSet) Len() int {
	return len(s.values)
}

// Remove returns a new FormulaSet without the specified value.
// If s does not contain the value, a copy of s will be returned.
func Remove(s FormulaSet, value formula.Formula) FormulaSet {
	res := maps.Clone(s.values)
	delete(res, value)
	return FormulaSet{res}
}

// Clone returns a copy of the given FormulaSet.
func Clone(s FormulaSet) FormulaSet {
	return FormulaSet{maps.Clone(s.values)}
}
