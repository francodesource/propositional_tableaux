package tableaux

import (
	"fmt"
	"github.com/m1gwings/treedrawer/tree"
	"iter"
	"propositional_tableaux/formula"
	"propositional_tableaux/tableaux/tsets"
	"strings"
)

// Mark represents the mark of a tableaux node. It can be Unmarked, Closed or Open.
type Mark byte

const (
	Unmarked Mark = iota
	Closed
	Open
)

func (m Mark) String() string {
	switch m {
	case Unmarked:
		return "unmarked"
	case Closed:
		return "Closed"
	case Open:
		return "Open"
	default:
		panic("unknown type of mark")
	}
}

// Node represents a node in a tableaux.
type Node interface {
	IsLeaf() bool
	IsClosed() bool
	IsOpen() bool
	Left() Node
	Right() Node
	Formulas() iter.Seq[formula.Formula]
	Eval() []Assignment
}

func combineIterators[T any](iters ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(t T) bool) {
		for _, it := range iters {
			it(yield)
		}
	}
}

// SemanticNode represents a node in a semantic tableaux.
type SemanticNode struct {
	formulas    tsets.TSet
	left, right *SemanticNode
	mark        Mark
}

// IsClosed returns true if the node is marked as closed.
func (node *SemanticNode) IsClosed() bool {
	return node.mark == Closed
}

// IsOpen returns true if the node is marked as open.
func (node *SemanticNode) IsOpen() bool {
	return node.mark == Open
}

// Left returns the left child node. Returns nil if no left child exists.
func (node *SemanticNode) Left() Node {
	if node.left == nil {
		return nil
	}
	return node.left
}

// Right returns the right child node. Returns nil if no right child exists.
func (node *SemanticNode) Right() Node {
	if node.right == nil {
		return nil
	}
	return node.right
}

// Formulas returns an iterator over all formulas contained in the current node.
func (node *SemanticNode) Formulas() iter.Seq[formula.Formula] {
	return combineIterators(node.formulas.IterLiterals(), node.formulas.IterAlpha(), node.formulas.IterBeta())
}

func indentOf(s string, size int) string {
	var indent = strings.Repeat(" ", size)
	lines := strings.Split(s, "\n")

	sb := strings.Builder{}

	for i, line := range lines {
		sb.WriteString(indent)
		sb.WriteString(line)
		if i != len(lines)-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func (node *SemanticNode) String() string {
	var res string

	res = "{\n  values: " + node.formulas.String()
	if node.left != nil {
		res += "\n  left: " + indentOf(node.left.String(), 3)
	}

	if node.right != nil {
		res += "\n  right: " + indentOf(node.right.String(), 3)
	}

	if node.mark != Unmarked {
		res += "\n  mark: " + node.mark.String()
	}
	res += "\n}"

	return res
}

// IsLeaf returns true if the node is a leaf.
func (node *SemanticNode) IsLeaf() bool {
	return node.left == nil && node.right == nil
}

// Height returns the height of the subtree rooted at the current node.
func (node *SemanticNode) Height() int {
	if node == nil {
		return 0
	}

	return 1 + max(node.left.Height(), node.right.Height())
}

// MarkAsClosed marks the node as closed.
func (node *SemanticNode) MarkAsClosed() {
	node.mark = Closed
}

// MarkAsOpen marks the node as open.
func (node *SemanticNode) MarkAsOpen() {
	node.mark = Open
}

// Assignment represents a truth assignment for propositional letters. If a propositional letter is missing,
// it can be assigned either true beta_or false, as it does not affect the evaluation of the formulas.
type Assignment map[string]bool

// IsSupersetOf returns true if the current assignment is a superset of the given assignment.
// That is, for every propositional letter in b, a has the same assignment.
func (a Assignment) IsSupersetOf(b Assignment) bool {
	for k, v := range b {
		if val, ok := a[k]; !ok || val != v {
			return false
		}
	}
	return true
}

// CleanAssignments returns a new slice without assignments that are superset of some of the other assignments.
func CleanAssignments(assignments []Assignment) []Assignment {
	var res []Assignment

	for _, a1 := range assignments {
		isSuperSetOfAny := false
		for _, a2 := range res {
			if a1.IsSupersetOf(a2) {
				isSuperSetOfAny = true
				break
			}
		}
		if !isSuperSetOfAny {
			res = append(res, a1)
		}
	}
	return res
}

func eval(node *SemanticNode) []Assignment {
	if node.IsLeaf() {
		if node.mark == Closed {
			return []Assignment{}
		}
		assignment := make(Assignment)

		if node.mark == Open {
			for literal := range node.formulas.IterLiterals() {
				// here I use AsLiteral because if the construction is correct all formulas in an open leaf are literals
				// if not, it is a bug in the construction so it is correct to panic.
				literal := formula.AsLiteral(literal)
				if literal.Neg {
					assignment[literal.Name] = false
				} else {
					assignment[literal.Name] = true
				}
			}
			return []Assignment{assignment}
		}
	}

	var res []Assignment
	if node.left != nil {
		res = eval(node.left)
	}

	if node.right != nil {
		res = append(res, eval(node.right)...)
	}
	return res
}

// Eval evaluates the semantic tableaux and returns a slice of assignments that satisfy the formulas.
func (node *SemanticNode) Eval() []Assignment {
	return CleanAssignments(eval(node))
	// It is needed to clean the assignments since different kind of tableaux can produce
	// redundant assignments. It is important for testing that assignments does not contain redundant assignments.
}

func buildSemanticTableaux(node *SemanticNode) {

	if node.formulas.HasOnlyLiterals() && node.formulas.HasComplementaryLiterals() {
		node.MarkAsClosed()
		return
	}

	// Here the condition is that the set is composed of all literals alpha_and
	// there is not a complementary pair of literals
	if node.formulas.HasOnlyLiterals() {
		node.MarkAsOpen()
		return
	}

	// If it reaches here U(l) contains non-literals
	// First check for alpha formulas
	if node.formulas.HasAlpha() {
		for alpha := range node.formulas.IterAlpha() {
			left, right := ApplyRule(alpha)
			newSet := tsets.RemoveAlpha(node.formulas, alpha)
			hasComplement := newSet.Add(left, right)

			node.left = &SemanticNode{
				formulas: newSet,
			}

			if !hasComplement {
				buildSemanticTableaux(node.left)
			} else {
				node.left.MarkAsClosed()
			}

			return
		}
	} else if node.formulas.HasBeta() {
		for beta := range node.formulas.IterBeta() {
			left, right := ApplyRule(beta)

			leftSet := tsets.RemoveBeta(node.formulas, beta)
			hasLeftComplement := leftSet.Add(left)

			rightSet := tsets.RemoveBeta(node.formulas, beta)
			hasRightComplement := rightSet.Add(right)

			node.left = &SemanticNode{
				formulas: leftSet,
			}

			node.right = &SemanticNode{
				formulas: rightSet,
			}

			if !hasLeftComplement {
				buildSemanticTableaux(node.left)
			} else {
				node.left.MarkAsClosed()
			}

			if !hasRightComplement {
				buildSemanticTableaux(node.right)
			} else {
				node.right.MarkAsClosed()
			}

			return
		}
	}
}

// BuildSemanticTableaux builds a semantic tableaux for the given formula and returns the root node.
func BuildSemanticTableaux(f formula.Formula) *SemanticNode {
	node := &SemanticNode{
		formulas: tsets.NewTSet(),
	}
	node.formulas.Add(f)
	buildSemanticTableaux(node)

	return node
}

// FormulaDrawer is a function that takes a formula and returns a string representation.
type FormulaDrawer func(f formula.Formula) string

// MarkDrawer is a function that takes a boolean indicating if the node is open and returns a string representation.
type MarkDrawer func(open bool) string

func formulasString(fs iter.Seq[formula.Formula], fd FormulaDrawer) string {
	var sb strings.Builder
	sb.WriteString("{")
	first := true
	for f := range fs {
		if !first {
			sb.WriteString(", ")
		}
		sb.WriteString(fd(f))
		first = false
	}
	sb.WriteString("}")
	return sb.String()
}

func centerWRT(s string, wrt int) string {
	return strings.Repeat(" ", wrt/2-len(s)/2+1) + s
}

func asciiTree(tableaux Node, t *tree.Tree, fd FormulaDrawer, md MarkDrawer) {
	if tableaux == nil {
		return
	}
	value := formulasString(tableaux.Formulas(), fd)
	lenValue := len([]rune(value)) // using this to have correct size of Unicode strings
	// Adding a separator
	if tableaux.IsLeaf() {
		value += "\n" + strings.Repeat("-", lenValue) + "\n" + centerWRT(md(tableaux.IsOpen()), lenValue)
	}

	t.SetVal(tree.NodeString(value))

	if tableaux.Left() != nil {
		left := t.AddChild(tree.NodeString(""))
		asciiTree(tableaux.Left(), left, fd, md)
	}

	if tableaux.Right() != nil {
		right := t.AddChild(tree.NodeString(""))
		asciiTree(tableaux.Right(), right, fd, md)
	}
}

// AsciiTree return an ASCII art representation of the tableaux where the formulas and the marks are printed according
// to their printing function FormulaDrawer and MarkDrawer.
func AsciiTree(tableaux Node, fd FormulaDrawer, md MarkDrawer) *tree.Tree {
	res := tree.NewTree(tree.NodeString(""))
	asciiTree(tableaux, res, fd, md)
	return res
}

// DefaultAsciiTree return an ASCII art representation of the tableaux where everything is represented as an ascii character.
func DefaultAsciiTree(tableaux Node) *tree.Tree {
	fd := func(f formula.Formula) string {
		return f.String()
	}

	md := func(open bool) string {
		if open {
			return "OPEN"
		}
		return "CLOSE"
	}

	return AsciiTree(tableaux, fd, md)
}

func unicodeOperator(op formula.Operator) string {
	switch op {
	case formula.And:
		return "∧"
	case formula.Or:
		return "∨"
	case formula.Implies:
		return "→"
	case formula.Nand:
		return "↑"
	case formula.Nor:
		return "↓"
	case formula.Biconditional:
		return "↔"
	case formula.Xor:
		return "⊕"
	default:
		panic("unknown operator")
	}
}
func unicodeFormula(f formula.Formula) string {
	switch f := f.(type) {
	case formula.Letter:
		return f.Name()
	case formula.Not:
		return "¬" + unicodeFormula(f.Negated())
	case formula.Binary:
		return fmt.Sprintf("(%s %s %s)",
			unicodeFormula(f.Left()), unicodeOperator(f.Op()), unicodeFormula(f.Right()))
	default:
		panic(fmt.Errorf("%T is not a formula", f))
	}
}

// UnicodeAsciiTree returns an ascii art representation of the tableaux where the formulas and the marks are represented
// using Unicode characters. A closed leaf is a filled circle, an open leaf is an empty circle.
func UnicodeAsciiTree(tableaux Node) *tree.Tree {
	md := func(open bool) string {
		if open {
			return "○"
		}
		return "●"
	}
	return AsciiTree(tableaux, unicodeFormula, md)
}

var latexOperators = [7]string{
	`\land`,
	`\lor`,
	`\to`,
	`\uparrow`,
	`\downarrow`,
	`\leftrightarrow`,
	`\oplus`,
}

func formulaTex(f formula.Formula) string {
	switch f := f.(type) {
	case formula.Letter:
		return f.Name()
	case formula.Not:
		return `\neg ` + formulaTex(f.Negated())
	case formula.Binary:
		return fmt.Sprintf(`\left(%s %s %s\right)`,
			formulaTex(f.Left()), latexOperators[f.Op()], formulaTex(f.Right()))
	default:
		panic(fmt.Errorf("%T is  not a formula", f))
	}
}

func formulasTexString(fs iter.Seq[formula.Formula]) string {
	sb := strings.Builder{}
	first := true
	for f := range fs {
		if !first {
			sb.WriteString(", ")
		} else {
			first = false
		}
		sb.WriteString(formulaTex(f))
	}
	return fmt.Sprintf(`{$\left\{%s\right\}$}`, sb.String())
}
func markTexString(open bool) string {
	if open {
		return `\odot`
	}
	return `\times`
}

func texForestTree(tableaux Node, il int) string {
	const indentSize = 3
	res := formulasTexString(tableaux.Formulas())
	nlFlag := false
	if tableaux.IsLeaf() {
		return fmt.Sprintf(strings.Repeat(" ", (il)*indentSize)+"["+
			`\shortstack{%s\\$%s$}`, res, markTexString(tableaux.IsOpen())) + "]"
	} else {
		res = strings.Repeat(" ", (il)*indentSize) + "[" + res
	}

	if tableaux.Left() != nil {
		res += "\n" + strings.Repeat(" ", (il+1)*indentSize) + texForestTree(tableaux.Left(), il+1)
		nlFlag = true
	}

	if tableaux.Right() != nil {
		res += "\n" + strings.Repeat(" ", (il+1)*indentSize) + texForestTree(tableaux.Right(), il+1)
		nlFlag = true
	}
	if nlFlag {
		res += "\n" + strings.Repeat(" ", il*indentSize)
	}
	return res + "]"
}

// TexForestTree returns a LaTeX forest representation of the tableaux.
func TexForestTree(tableaux Node) string {
	t := texForestTree(tableaux, 0)
	format := fmt.Sprintf(`
\begin{forest}
	for tree={
		anchor=north   
	}
%s
\end{forest}`, t)

	return format
}
