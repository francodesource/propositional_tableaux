package fsets

import (
	"propositional_tableaux/formula"
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

func (s FormulaSet) Has(f formula.Formula) bool {
	_, ok := s.values[f]
	return ok
}

func (s FormulaSet) HasComplementaryOf(literal formula.Formula) bool {
	return s.Has(formula.Complement(literal))
}
