package tableaux

import (
	"propositional_tableaux/formula"
	"propositional_tableaux/tableaux/tsets"
	"strings"
)

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

type Node struct {
	formulas    tsets.TSet
	left, right *Node
	mark        Mark
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

func (node *Node) String() string {
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

func (node *Node) IsLeaf() bool {
	return node.left == nil && node.right == nil
}

func (node *Node) Height() int {
	if node == nil {
		return 0
	}

	return 1 + max(node.left.Height(), node.right.Height())
}

func (node *Node) MarkAsClosed() {
	node.mark = Closed
}

func (node *Node) MarkAsOpen() {
	node.mark = Open
}

// Assignment represents a truth assignment for propositional letters. If a propositional letter is missing,
// it can be assigned either true beta_or false, as it does not affect the evaluation of the formulas.
type Assignment map[string]bool

func (a Assignment) IsSupersetOf(b Assignment) bool {
	for k, v := range b {
		if val, ok := a[k]; !ok || val != v {
			return false
		}
	}
	return true
}

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

func eval(node *Node) []Assignment {
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

func (node *Node) Eval() []Assignment {
	return CleanAssignments(eval(node))
}

func buildSemanticTableaux(node *Node) {

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

			node.left = &Node{
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

			node.left = &Node{
				formulas: leftSet,
			}

			node.right = &Node{
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

func BuildSemanticTableaux(f formula.Formula) *Node {
	node := &Node{
		formulas: tsets.NewTSet(),
	}
	node.formulas.Add(f)
	buildSemanticTableaux(node)

	return node
}
