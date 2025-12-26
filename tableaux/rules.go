package tableaux

import (
	"fmt"
	"github.com/francodesource/propositional_tableaux/formula"
)

// Rule represent a rule of the tableaux calculus that takes two formulas and returns two formulas.
type Rule func(left, right formula.Formula) (formula.Formula, formula.Formula)

func and(left, right formula.Formula) (formula.Formula, formula.Formula) {
	return left, right
}

func or(left, right formula.Formula) (formula.Formula, formula.Formula) {
	return left, right
}

func not_and(left, right formula.Formula) (formula.Formula, formula.Formula) {
	return formula.NewNot(left), formula.NewNot(right)
}

func not_or(left, right formula.Formula) (formula.Formula, formula.Formula) {
	return formula.NewNot(left), formula.NewNot(right)
}
func implies(left, right formula.Formula) (formula.Formula, formula.Formula) {
	return formula.NewNot(left), right
}

func not_implies(left, right formula.Formula) (formula.Formula, formula.Formula) {
	return left, formula.NewNot(right)
}

func biconditional(left, right formula.Formula) (formula.Formula, formula.Formula) {
	return formula.NewImplies(left, right), formula.NewImplies(right, left)
}

func not_biconditional(left, right formula.Formula) (formula.Formula, formula.Formula) {
	return formula.NewNot(formula.NewImplies(left, right)), formula.NewNot(formula.NewImplies(right, left))
}

const OperatorsCount = 7

// alphaRules is the set of rules for alpha formulas. It assumes negation and consider inner operator.
// Double negation is considered to be a special case, so it is not contained in this array.
var alphaRules = [OperatorsCount]Rule{
	and,
	not_or,
	not_implies,
	and,    // negated nand is the same as and
	not_or, // not or is the same as nor
	biconditional,
	biconditional, // biconditional is '=', xor is '!=' so negated xor is the same as biconditional
}

var betaRules = [OperatorsCount]Rule{
	not_and,
	or,
	implies,
	not_and,
	or,
	not_biconditional,
	not_biconditional,
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
