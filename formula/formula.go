package formula

import "fmt"

type Formula interface {
	String() string
}

type Letter struct {
	name string
}

func NewLetter(name string) Letter {
	return Letter{name: name}
}

func (l Letter) Name() string {
	return l.name
}

func (l Letter) String() string {
	return l.name
}

type Not struct {
	negated Formula
}

func NewNot(negated Formula) Not {
	return Not{negated: negated}
}

func (n Not) String() string {
	return "!" + n.negated.String()
}

// Enumeration for logical operators

// Operator represents logical operators
type Operator int

const (
	And Operator = iota
	Or
	Implies
	Nand
	Nor
	Biconditional
	Xor
)

func (o Operator) String() string {
	switch o {
	case And:
		return "&"
	case Or:
		return "|"
	case Implies:
		return "->"
	case Nand:
		return "!&"
	case Nor:
		return "!|"
	case Biconditional:
		return "<->"
	case Xor:
		return "^"
	default:
		panic(fmt.Errorf("unknown operator type %T", o)) // unreachable
	}
}

type Binary struct {
	left, right Formula
	op          Operator
}

func NewBinary(left, right Formula, op Operator) Binary {
	return Binary{left: left, right: right, op: op}
}

func NewAnd(left, right Formula) Binary {
	return NewBinary(left, right, And)
}

func NewOr(left, right Formula) Binary {
	return NewBinary(left, right, Or)
}

func NewImplies(left, right Formula) Binary {
	return NewBinary(left, right, Implies)
}

func NewNand(left, right Formula) Binary {
	return NewBinary(left, right, Nand)
}

func NewNor(left, right Formula) Binary {
	return NewBinary(left, right, Nor)
}

func NewBiconditional(left, right Formula) Binary {
	return NewBinary(left, right, Biconditional)
}
func NewXor(left, right Formula) Binary {
	return NewBinary(left, right, Xor)
}

func (b Binary) String() string {
	res := fmt.Sprintf("(%s %s %s)", b.left, b.op, b.right)

	return res
}
