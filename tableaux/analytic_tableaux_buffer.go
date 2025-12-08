package tableaux

import (
	"fmt"
	"maps"
	"propositional_tableaux/formula"
	"strings"
)

// BufferSet is a buffer of two different formulas
type BufferSet []formula.Formula

func NewBufferSet(fs ...formula.Formula) BufferSet {
	res := make(BufferSet, 2)

	for _, f := range fs {
		res.Add(f)
	}
	return res
}

func (b BufferSet) String() string {
	strs := []string{}

	if b[0] != nil {
		strs = append(strs, b[0].String())
	}

	if b[1] != nil {
		strs = append(strs, b[1].String())
	}

	return "{" + strings.Join(strs, ", ") + "}"
}

func (b BufferSet) Add(f formula.Formula) {
	if b[0] == nil {
		b[0] = f
		return
	}

	if b[0] == f {
		return // if the formula is already contained it skips.
	}

	if b[1] == nil {
		b[1] = f
		return
	}

	if b[1] == f {
		return // if the formula is already contained it skips.
	}

	panic(fmt.Errorf("BufferSet can't contain more than two values: can't add %v to %v", b, f))
}

func (b BufferSet) Has(f formula.Formula) bool {
	return b[0] == f || b[1] == f
}

func (b BufferSet) HasComplementOf(fs ...formula.Formula) bool {
	for _, f := range fs {
		if b.Has(formula.Complement(f)) {
			return true
		}
	}

	return false
}

type BufferedNode struct {
	formulas            BufferSet
	left, right, father *BufferedNode
	mark                Mark
}

func (b *BufferedNode) BranchHasComplementPairOf(fs ...formula.Formula) bool {
	for _, f := range fs {
		if f != nil && b.formulas.HasComplementOf(f) {
			return true
		}
	}

	if b.father != nil {
		return b.father.BranchHasComplementPairOf(fs...)
	}

	return false
}

// ChooseAlphaFormula chooses a random alpha formula from the current node branch.
// Returns nil if no alpha formula is available
func (b *BufferedNode) ChooseAlphaFormula(visited map[formula.Formula]bool) (res formula.Formula) {
	for _, f := range b.formulas {
		if f != nil && f.Class() == formula.Alpha && !visited[f] {
			return f
		}
	}

	if b.father != nil {
		return b.father.ChooseAlphaFormula(visited)
	}

	return
}

// ChooseBetaFormula chooses a random beta formula from the current node branch.
// Returns nil if no beta formula is available
func (b *BufferedNode) ChooseBetaFormula(visited map[formula.Formula]bool) (res formula.Formula) {
	for _, f := range b.formulas {
		if f != nil && f.Class() == formula.Beta && !visited[f] {
			return f
		}
	}

	if b.father != nil {
		return b.father.ChooseBetaFormula(visited)
	}

	return
}

func (b *BufferedNode) MarkAsClosed() {
	b.mark = Closed
}

func (b *BufferedNode) MarkAsOpen() {
	b.mark = Open
}

func (b *BufferedNode) IsLeaf() bool {
	return b.left == nil && b.right == nil
}

func (b *BufferedNode) IsClosed() bool {
	return b.IsLeaf() && b.mark == Closed
}

func (b *BufferedNode) IsOpen() bool {
	return b.IsLeaf() && b.mark == Open
}

func (b *BufferedNode) String() string {
	var res string

	res = "{\n  values: " + b.formulas.String()
	if b.left != nil {
		res += "\n  left: " + indentOf(b.left.String(), 3)
	}

	if b.right != nil {
		res += "\n  right: " + indentOf(b.right.String(), 3)
	}

	if b.mark != Unmarked {
		res += "\n  mark: " + b.mark.String()
	}
	res += "\n}"

	return res
}

func (b *BufferedNode) collectLiterals(literals map[formula.Literal]bool) {
	for _, f := range b.formulas {
		if f != nil && formula.IsLiteral(f) {
			literals[formula.AsLiteral(f)] = true
		}
	}

	if b.father != nil {
		b.father.collectLiterals(literals)
	}
}

func (b *BufferedNode) eval() []Assignment {

	if b.IsClosed() {
		return []Assignment{}
	}

	if b.IsOpen() {
		assignment := make(Assignment)
		literals := make(map[formula.Literal]bool)
		b.collectLiterals(literals)

		for lit := range literals {
			assignment[lit.Name] = assignLiteral(lit)
		}

		return []Assignment{assignment}
	}

	var res []Assignment
	if b.left != nil {
		res = b.left.eval()
	}

	if b.right != nil {
		res = append(res, b.right.eval()...)
	}

	return res
}

func (b *BufferedNode) Eval() []Assignment {
	return b.eval()
}

func buildBufferedTableaux(a *BufferedNode, visited map[formula.Formula]bool) {
	alpha := a.ChooseAlphaFormula(visited)

	if alpha != nil {
		visited[alpha] = true

		left, right := ApplyRule(alpha)
		newBuff := NewBufferSet(left, right)

		a.left = &BufferedNode{
			formulas: newBuff,
			father:   a,
		}

		if !a.left.BranchHasComplementPairOf(left, right) {
			buildBufferedTableaux(a.left, visited)
		} else {
			a.left.MarkAsClosed()
		}
		return
	}

	beta := a.ChooseBetaFormula(visited)

	if beta != nil {
		visited[beta] = true
		left, right := ApplyRule(beta)

		leftSet := NewBufferSet(left)

		rightSet := NewBufferSet(right)

		a.left = &BufferedNode{
			formulas: leftSet,
			father:   a,
		}

		a.right = &BufferedNode{
			formulas: rightSet,
			father:   a,
		}

		if !a.left.BranchHasComplementPairOf(left) {
			buildBufferedTableaux(a.left, maps.Clone(visited))
		} else {
			a.left.MarkAsClosed()
		}

		if !a.right.BranchHasComplementPairOf(right) {
			buildBufferedTableaux(a.right, maps.Clone(visited))
		} else {
			a.right.MarkAsClosed()
		}
		return
	}

	a.MarkAsOpen()

}

func BuildBufferedTableaux(f formula.Formula) *BufferedNode {
	set := NewBufferSet(f)

	res := &BufferedNode{formulas: set}
	buildBufferedTableaux(res, make(map[formula.Formula]bool))

	return res
}
