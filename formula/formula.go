package formula

import (
	"fmt"
	"math/rand"
)

type Classification int

const (
	Alpha                       = -1
	LiteralClass Classification = 0
	Beta                        = 1
)

func (c Classification) String() string {
	switch c {
	case LiteralClass:
		return "LiteralClass"
	case Alpha:
		return "Alpha"
	case Beta:
		return "Beta"
	default:
		panic("unknown Classification")
	}
}

type Formula interface {
	Class() Classification
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

func (l Letter) Class() Classification {
	return LiteralClass
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

func (n Not) Negated() Formula {
	return n.negated
}

func (n Not) Class() Classification {
	switch inner := n.negated.(type) {
	case Letter:
		return LiteralClass
	case Not:
		return Alpha
	case Binary:
		return inner.Class() * -1

	default:
		panic(fmt.Errorf("%v: %T is not a Formula", inner, inner))
	}
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

func (b Binary) Left() Formula {
	return b.left
}

func (b Binary) Right() Formula {
	return b.right
}

func (b Binary) Op() Operator {
	return b.op
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

func (b Binary) Class() Classification {
	switch b.op {
	case And, Nor, Biconditional:
		return Alpha
	default:
		return Beta
	}
}

func (b Binary) String() string {
	res := fmt.Sprintf("(%s %s %s)", b.left, b.op, b.right)

	return res
}

func IsLiteral(formula Formula) bool {
	if _, ok := formula.(Letter); ok {
		return true
	}

	if not, ok := formula.(Not); ok {
		_, ok := not.Negated().(Letter)
		return ok
	}

	return false
}

// Literal is either a letter or its negation
type Literal struct {
	Name string
	Neg  bool
}

// AsLiteral converts a Formula to a Literal. It panics if the formula is not a literal.
// It should be used only after checking with IsLiteral.
func AsLiteral(formula Formula) Literal {
	if letter, ok := formula.(Letter); ok {
		return Literal{
			Name: letter.Name(),
			Neg:  false,
		}
	}

	if not, ok := formula.(Not); ok {
		if letter, ok := not.Negated().(Letter); ok {
			return Literal{
				Name: letter.Name(),
				Neg:  true,
			}
		}
	}

	panic(fmt.Errorf("%v is not a literal", formula))
}

func Complement(formula Formula) Formula {
	if not, ok := formula.(Not); ok {
		return not.Negated()
	} else {
		return NewNot(formula)
	}
}

// GenerateRandom generates a random formula of the given size.
func GenerateRandom(rand *rand.Rand, size int) Formula {
	if size <= 0 {
		letters := "pqrstuvwxyz"
		randLetter := letters[rand.Intn(len(letters))]
		return NewLetter(string(randLetter))
	}

	negation := rand.Intn(2) == 1

	if negation {
		return NewNot(GenerateRandom(rand, size-1))
	} else {
		return NewBinary(GenerateRandom(rand, size/2), GenerateRandom(rand, size/2),
			Operator(rand.Intn(7))) // random binary operation
	}
}
