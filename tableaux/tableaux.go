package tableaux

import (
	"propositional_tableaux/formula"
	"propositional_tableaux/tableaux/fsets"
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
	formulas    fsets.FormulaSet
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

func eval(node *Node) []map[string]bool {
	if node.IsLeaf() {
		if node.mark == Closed {
			return []map[string]bool{}
		}
		assignment := make(map[string]bool)

		if node.mark == Open {
			for literal := range node.formulas.Iter() {
				switch literal := literal.(type) {
				case formula.Letter:
					assignment[literal.Name()] = true
				case formula.Not:
					letter, ok := literal.Negated().(formula.Letter)
					if !ok {
						panic("malformed tableaux")
					}
					assignment[letter.Name()] = false
				default:
					panic("malformed tableaux")
				}
			}
			return []map[string]bool{assignment}
		}
	}

	var res []map[string]bool
	if node.left != nil {
		res = eval(node.left)
	}

	if node.right != nil {
		res = append(res, eval(node.right)...)
	}
	return res
}

func (node *Node) Eval() []map[string]bool {
	return eval(node)
}

func buildSemanticTableaux(node *Node) {
	allLiterals := true
	complementaryPair := false

	for f := range node.formulas.Iter() {
		if !formula.IsLiteral(f) {
			allLiterals = false
			break
		} else {
			if node.formulas.HasComplementaryOf(f) {
				complementaryPair = true
			}
		}
	}

	if allLiterals && complementaryPair {
		node.mark = Closed
		return
	}

	// Here the condition is that the set is composed of all literals and
	// there is not a complementary pair of literals
	if allLiterals {
		node.mark = Open
		return
	}

	// If it reaches here U(l) contains non-literals
	for f := range node.formulas.Iter() {
		switch f.Class() {
		case formula.Literal:
			continue
		case formula.Alpha:
			left, right := ApplyRule(f)
			newSet := fsets.Remove(node.formulas, f).Add(left, right)
			node.left = &Node{
				formulas: newSet,
			}
			buildSemanticTableaux(node.left)
		case formula.Beta:
			left, right := ApplyRule(f)
			leftSet := fsets.Remove(node.formulas, f).Add(left)
			rightSet := fsets.Remove(node.formulas, f).Add(right)

			node.left = &Node{
				formulas: leftSet,
			}

			node.right = &Node{
				formulas: rightSet,
			}

			buildSemanticTableaux(node.left)
			buildSemanticTableaux(node.right)
		}
	}
}

func BuildSemanticTableaux(f formula.Formula) *Node {
	node := &Node{
		formulas: fsets.New(f),
	}

	buildSemanticTableaux(node)

	return node
}
