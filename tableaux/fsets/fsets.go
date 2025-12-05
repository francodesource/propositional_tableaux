package fsets

import (
	"iter"
	"maps"
	"propositional_tableaux/formula"
	"strings"
)

type FormulaSet struct {
	values map[formula.Formula]struct{}
}

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

func (s FormulaSet) Has(f formula.Formula) bool {
	_, ok := s.values[f]
	return ok
}

func (s FormulaSet) HasComplementaryOf(literal formula.Formula) bool {
	return s.Has(formula.Complement(literal))
}

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

func Clone(s FormulaSet) FormulaSet {
	return FormulaSet{maps.Clone(s.values)}
}
