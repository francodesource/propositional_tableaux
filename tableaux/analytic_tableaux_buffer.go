package tableaux

import (
	"fmt"
	"iter"
	"maps"
	"propositional_tableaux/formula"
	"slices"
	"strings"
)

// BufferSet is a buffer of two different formulas
type BufferSet []formula.Formula

// NewBufferSet creates a new BufferSet containing the given formulas.
// Panics if more than two different formulas are given.
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

// Add adds a formula to the BufferSet.
// Panics if more than two different formulas are added.
// If the formula is already contained it will not be added.
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

// Has returns true if the formula f is contained in the BufferSet.
func (b BufferSet) Has(f formula.Formula) bool {
	return b[0] == f || b[1] == f
}

// HasComplementOf returns true if the BufferSet contains the complement of at least one of the given formulas.
func (b BufferSet) HasComplementOf(fs ...formula.Formula) bool {
	for _, f := range fs {
		if b.Has(formula.Complement(f)) {
			return true
		}
	}

	return false
}

func (b BufferSet) HasOnlyLiterals() bool {
	return b[0] == nil || formula.IsLiteral(b[0]) && b[1] == nil || formula.IsLiteral(b[1])
}

// BufferNode is a node in a buffered analytic tableaux.
type BufferNode struct {
	formulas            BufferSet
	left, right, father *BufferNode
	mark                Mark
}

// Father returns the father node of the current node. Returns nil if the node is root.
func (b *BufferNode) Father() *BufferNode {
	if b.father == nil {
		return nil
	}
	return b.father
}

// Left returns the left child node of the current node. Returns nil if no left child exists.
func (b *BufferNode) Left() Node {
	if b.left == nil {
		return nil
	}
	return b.left
}

// Right returns the right child node of the current node. Returns nil if no right child exists.
func (b *BufferNode) Right() Node {
	if b.right == nil {
		return nil
	}
	return b.right
}

// Formulas returns an iterator over all formulas contained in the current node, filtered of nil values.
func (b *BufferNode) Formulas() iter.Seq[formula.Formula] {
	res := make([]formula.Formula, 0, 2)

	for _, f := range b.formulas {
		if f != nil {
			res = append(res, f)
		}
	}

	return slices.Values(res)
}

// BranchHasComplementPairOf checks if the current node branch has a complement pair of at least one of the given formulas.
func (b *BufferNode) BranchHasComplementPairOf(fs ...formula.Formula) bool {
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
func (b *BufferNode) ChooseAlphaFormula(visited map[formula.Formula]bool) (res formula.Formula) {
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
func (b *BufferNode) ChooseBetaFormula(visited map[formula.Formula]bool) (res formula.Formula) {
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

// MarkAsClosed marks the current node as closed.
func (b *BufferNode) MarkAsClosed() {
	b.mark = Closed
}

// MarkAsOpen marks the current node as open.
func (b *BufferNode) MarkAsOpen() {
	b.mark = Open
}

// IsLeaf returns true if the current node is a leaf (i.e., has no children).
func (b *BufferNode) IsLeaf() bool {
	return b.left == nil && b.right == nil
}

// IsClosed returns true if the current node is a closed leaf.
func (b *BufferNode) IsClosed() bool {
	return b.IsLeaf() && b.mark == Closed
}

// IsOpen returns true if the current node is an open leaf.
func (b *BufferNode) IsOpen() bool {
	return b.IsLeaf() && b.mark == Open
}

func (b *BufferNode) String() string {
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

func (b *BufferNode) collectLiterals(literals map[formula.Literal]bool) {
	for _, f := range b.formulas {
		if f != nil && formula.IsLiteral(f) {
			literals[formula.AsLiteral(f)] = true
		}
	}

	if b.father != nil {
		b.father.collectLiterals(literals)
	}
}

func (b *BufferNode) eval() []Assignment {

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

// Eval evaluates the buffered tableaux and returns all satisfying assignments cleaned of the redundant assignments.
func (b *BufferNode) Eval() []Assignment {
	return CleanAssignments(b.eval())
}

func buildBufferedTableaux(a *BufferNode, visited map[formula.Formula]bool) {
	alpha := a.ChooseAlphaFormula(visited)

	if alpha != nil {
		visited[alpha] = true

		left, right := ApplyRule(alpha)
		newBuff := NewBufferSet(left, right)

		a.left = &BufferNode{
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

		a.left = &BufferNode{
			formulas: leftSet,
			father:   a,
		}

		a.right = &BufferNode{
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

	// If it reaches this point it means that it is not possible to apply any rule so the algorithm has finished,
	// Except for the special case where an alpha formula contains two equals subformulas: the second one will already mark
	// as visited, and it will create a leaf that does not contain just literals. So I check for this and enforce the
	// application of the rule.

	if !a.formulas.HasOnlyLiterals() {
		newVisited := maps.Clone(visited)
		for f := range a.Formulas() {
			newVisited[f] = false
		}
		buildBufferedTableaux(a, newVisited)
	}

	a.MarkAsOpen()

}

// BuildBufferTableaux builds a buffered analytic tableaux for the given formula.
func BuildBufferTableaux(f formula.Formula) *BufferNode {
	set := NewBufferSet(f)

	res := &BufferNode{formulas: set}
	buildBufferedTableaux(res, make(map[formula.Formula]bool))

	return res
}
