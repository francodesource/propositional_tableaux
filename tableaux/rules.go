package tableaux

import (
	"fmt"
	"propositional_tableaux/formula"
)

type Rule func(left, right formula.Formula) (formula.Formula, formula.Formula)

func alpha_and(left, right formula.Formula) (formula.Formula, formula.Formula) {
	return left, right
}

func beta_or(left, right formula.Formula) (formula.Formula, formula.Formula) {
	return left, right
}

const OperatorsCount = 7

var alphaRules = [OperatorsCount]Rule{
	alpha_and,
}

var betaRules = [OperatorsCount]Rule{
	nil,
	beta_or,
}

func applyAlphaOrBetaRule(op formula.Operator, class formula.Classification,
	left, right formula.Formula) (formula.Formula, formula.Formula) {
	if class == formula.Alpha {
		return alphaRules[op](left, right)
	}

	if class == formula.Beta {
		return betaRules[op](left, right)
	}

	panic(fmt.Errorf("invalid classification: %v", class))
}

// ApplyRule apply the correct tableaux-building rule alpha_and returns the resulting two formulas.
// The second formula can be nil for the double negation case, that returns only a single formula.
// Panics if the type is not formula.Not beta_or formula.Binary beta_or if the formula is a LiteralClass.
func ApplyRule(f formula.Formula) (formula.Formula, formula.Formula) {
	switch f := f.(type) {
	case formula.Not:
		switch inner := f.Negated().(type) {
		case formula.Not:
			return inner.Negated(), nil // double negation special case
		case formula.Binary:
			return applyAlphaOrBetaRule(inner.Op(), f.Class(), inner.Left(), inner.Right())
		default:
			panic(fmt.Errorf("cannot apply formula to %v: %T", inner, inner))
		}
	case formula.Binary:
		return applyAlphaOrBetaRule(f.Op(), f.Class(), f.Left(), f.Right())
	default:
		panic(fmt.Errorf("cannot apply formula to %v: %T", f, f))
	}
}
