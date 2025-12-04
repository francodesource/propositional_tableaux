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
