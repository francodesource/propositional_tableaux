package tableaux

import (
	"maps"
	"propositional_tableaux/formula"
	"propositional_tableaux/tableaux/tsets"
)

type AnalyticNode struct {
	formulas            tsets.TSet
	left, right, father *AnalyticNode
	mark                Mark
}

func (a *AnalyticNode) BranchHasComplementPairOf(fs ...formula.Formula) bool {
	for _, f := range fs {
		if f != nil && a.formulas.HasComplementOf(f) {
			return true
		}
	}

	if a.father != nil {
		return a.father.BranchHasComplementPairOf(fs...)
	}

	return false
}

func (a *AnalyticNode) BranchHasComplementaryLiterals() bool {
	if a.formulas.HasComplementaryLiterals() {
		return true
	}

	if a.father != nil {
		return a.father.BranchHasComplementaryLiterals()
	}

	return false
}

func (a *AnalyticNode) MarkAsClosed() {
	a.mark = Closed
}

func (a *AnalyticNode) MarkAsOpen() {
	a.mark = Open
}

func (a *AnalyticNode) String() string {
	var res string

	res = "{\n  values: " + a.formulas.String()
	if a.left != nil {
		res += "\n  left: " + indentOf(a.left.String(), 3)
	}

	if a.right != nil {
		res += "\n  right: " + indentOf(a.right.String(), 3)
	}

	if a.mark != Unmarked {
		res += "\n  mark: " + a.mark.String()
	}
	res += "\n}"

	return res
}

// ChooseAlphaFormula chooses a random alpha formula from the current node branch.
// Returns nil if no alpha formula is available
func (a *AnalyticNode) ChooseAlphaFormula(visited map[formula.Formula]bool) (res formula.Formula) {
	for alpha := range a.formulas.IterAlpha() {
		if !visited[alpha] {
			return alpha
		}
	}

	if a.father != nil {
		return a.father.ChooseAlphaFormula(visited)
	}

	return
}

// ChooseBetaFormula chooses a random beta formula from the current node branch.
// Returns nil if no beta formula is available
func (a *AnalyticNode) ChooseBetaFormula(visited map[formula.Formula]bool) (res formula.Formula) {
	for beta := range a.formulas.IterBeta() {
		if !visited[beta] {
			return beta
		}
	}

	if a.father != nil {
		return a.father.ChooseBetaFormula(visited)
	}

	return
}

func (a *AnalyticNode) IsLeaf() bool {
	return a.left == nil && a.right == nil
}

func (a *AnalyticNode) IsClosed() bool {
	return a.IsLeaf() && a.mark == Closed
}

func (a *AnalyticNode) IsOpen() bool {
	return a.IsLeaf() && a.mark == Open
}

func assignLiteral(literal formula.Literal) bool {
	if literal.Neg {
		return false
	}
	return true
}

func (a *AnalyticNode) collectLiterals(literals map[formula.Literal]bool) {
	for literal := range a.formulas.IterLiterals() {
		literals[formula.AsLiteral(literal)] = true
	}

	if a.father != nil {
		a.father.collectLiterals(literals)
	}
}

func (a *AnalyticNode) eval() []Assignment {

	if a.IsClosed() {
		return []Assignment{}
	}

	if a.IsOpen() {
		assignment := make(Assignment)
		literals := make(map[formula.Literal]bool)
		a.collectLiterals(literals)

		for lit := range literals {
			assignment[lit.Name] = assignLiteral(lit)
		}

		return []Assignment{assignment}
	}

	var res []Assignment
	if a.left != nil {
		res = a.left.eval()
	}

	if a.right != nil {
		res = append(res, a.right.eval()...)
	}

	return res
}

func (a *AnalyticNode) Eval() []Assignment {
	return CleanAssignments(a.eval())
}

func buildAnalyticTableaux(a *AnalyticNode, visited map[formula.Formula]bool) {
	alpha := a.ChooseAlphaFormula(visited)

	if alpha != nil {
		visited[alpha] = true

		left, right := ApplyRule(alpha)
		newSet := tsets.NewTSet()
		newSet.Add(left, right)

		a.left = &AnalyticNode{
			formulas: newSet,
			father:   a,
		}

		if !a.left.BranchHasComplementPairOf(left, right) {
			buildAnalyticTableaux(a.left, visited)
		} else {
			a.left.MarkAsClosed()
		}
		return
	}

	beta := a.ChooseBetaFormula(visited)

	if beta != nil {
		visited[beta] = true
		left, right := ApplyRule(beta)

		leftSet := tsets.NewTSet()
		leftSet.Add(left)

		rightSet := tsets.NewTSet()
		rightSet.Add(right)

		a.left = &AnalyticNode{
			formulas: leftSet,
			father:   a,
		}

		a.right = &AnalyticNode{
			formulas: rightSet,
			father:   a,
		}

		if !a.left.BranchHasComplementPairOf(left) {
			buildAnalyticTableaux(a.left, maps.Clone(visited))
		} else {
			a.left.MarkAsClosed()
		}

		if !a.right.BranchHasComplementPairOf(right) {
			buildAnalyticTableaux(a.right, maps.Clone(visited))
		} else {
			a.right.MarkAsClosed()
		}
		return
	}

	a.MarkAsOpen()

}

func BuildAnalyticTableaux(f formula.Formula) *AnalyticNode {
	set := tsets.NewTSet()
	set.Add(f)

	res := &AnalyticNode{formulas: set}
	buildAnalyticTableaux(res, make(map[formula.Formula]bool))

	return res
}
