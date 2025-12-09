// Code generated from formula/Formula.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // Formula

import "github.com/antlr4-go/antlr/v4"

// FormulaListener is a complete listener for a parse tree produced by FormulaParser.
type FormulaListener interface {
	antlr.ParseTreeListener

	// EnterStart is called when entering the start production.
	EnterStart(c *StartContext)

	// EnterBinary is called when entering the Binary production.
	EnterBinary(c *BinaryContext)

	// EnterNegation is called when entering the Negation production.
	EnterNegation(c *NegationContext)

	// EnterLetter is called when entering the Letter production.
	EnterLetter(c *LetterContext)

	// ExitStart is called when exiting the start production.
	ExitStart(c *StartContext)

	// ExitBinary is called when exiting the Binary production.
	ExitBinary(c *BinaryContext)

	// ExitNegation is called when exiting the Negation production.
	ExitNegation(c *NegationContext)

	// ExitLetter is called when exiting the Letter production.
	ExitLetter(c *LetterContext)
}
