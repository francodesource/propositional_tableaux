package formula

import (
	"fmt"
	"math/rand"
)

// Classification represent the classification of a formula: it can be alpha, beta or literal.
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

// Formula represent a formula in propositional logic.
type Formula interface {
	Class() Classification
	String() string
}

// Letter is a letter of a propositional formula.
type Letter struct {
	name string
}

// NewLetter returns a new Letter with the given name.
func NewLetter(name string) Letter {
	return Letter{name: name}
}

// Name returns the name of the Letter.
func (l Letter) Name() string {
	return l.name
}

// Class returns the Classification of the formula.
func (l Letter) Class() Classification {
	return LiteralClass
}

func (l Letter) String() string {
	return l.name
}

// Not is a negation of a Formula.
type Not struct {
	negated Formula
}

// NewNot creates a new negation of the given formula.
func NewNot(negated Formula) Not {
	return Not{negated: negated}
}

// Negated returns the inner formula inside the negation.
func (n Not) Negated() Formula {
	return n.negated
}

// Class returns the Classification of the formula.
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

// Operator represents logical operators.
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

// Binary represent any binary logic operation.
type Binary struct {
	left, right Formula
	op          Operator
}

// Left returns the left side formula of the Binary operation.
func (b Binary) Left() Formula {
	return b.left
}

// Right returns the right side formula of the Binary operation.
func (b Binary) Right() Formula {
	return b.right
}

// Op returns the Operator of the Binary operation.
func (b Binary) Op() Operator {
	return b.op
}

// NewBinary returns a new Binary formula between left and right with the given Operator.
func NewBinary(left, right Formula, op Operator) Binary {
	return Binary{left: left, right: right, op: op}
}

// NewAnd returns a new Binary formula representing the conjunction of left and right.
func NewAnd(left, right Formula) Binary {
	return NewBinary(left, right, And)
}

// NewOr returns a new Binary formula representing the disjunction of left and right.
func NewOr(left, right Formula) Binary {
	return NewBinary(left, right, Or)
}

// NewImplies returns a new Binary formula representing the implication from left to right.
func NewImplies(left, right Formula) Binary {
	return NewBinary(left, right, Implies)
}

// NewNand returns a new Binary formula representing the NAND of left and right.
func NewNand(left, right Formula) Binary {
	return NewBinary(left, right, Nand)
}

// NewNor returns a new Binary formula representing the NOR of left and right.
func NewNor(left, right Formula) Binary {
	return NewBinary(left, right, Nor)
}

// NewBiconditional returns a new Binary formula representing the biconditional of left and right.
func NewBiconditional(left, right Formula) Binary {
	return NewBinary(left, right, Biconditional)
}

// NewXor returns a new Binary formula representing the exclusive or of left and right.
func NewXor(left, right Formula) Binary {
	return NewBinary(left, right, Xor)
}

// Class returns the Classification of the formula.
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

// IsLiteral checks if the given formula is a literal (either a letter or its negation).
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

// Complement returns the complement of the given formula.
// If the formula is a negation, it returns the inner formula; otherwise, it returns the negation of the formula.
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
