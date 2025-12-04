package tableaux

import (
	"fmt"
	"propositional_tableaux/formula"
)

type RulePair struct {
	Neg bool
	Op  formula.Operator
}

type Rule func(left, right formula.Formula) (formula.Formula, formula.Formula)

func and(left, right formula.Formula) (formula.Formula, formula.Formula) {
	return left, right
}
func or(left, right formula.Formula) (formula.Formula, formula.Formula) {
	return left, right
}

var rules = map[RulePair]Rule{
	{false, formula.And}: and,
	{false, formula.Or}:  or,
}

// ApplyRule apply the correct tableaux-building rule and returns the resulting two formulas.
// The second formula can be nil for the double negation case, that returns only a single formula.
// Panics if the type is not formula.Not or formula.Binary or if the formula is a Literal.
func ApplyRule(f formula.Formula) (formula.Formula, formula.Formula) {
	switch f := f.(type) {
	case formula.Not:
		switch inner := f.Negated().(type) {
		case formula.Not:
			return inner.Negated(), nil
		case formula.Binary:
			return rules[RulePair{true, inner.Op()}](inner.Left(), inner.Right())
		default:
			panic(fmt.Errorf("cannot apply formula to %v: %T", inner, inner))
		}
	case formula.Binary:
		return rules[RulePair{false, f.Op()}](f.Left(), f.Right())
	default:
		panic(fmt.Errorf("cannot apply formula to %v: %T", f, f))
	}
}
