// Code generated from formula/Formula.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // Formula

import "github.com/antlr4-go/antlr/v4"

// BaseFormulaListener is a complete listener for a parse tree produced by FormulaParser.
type BaseFormulaListener struct{}

var _ FormulaListener = &BaseFormulaListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseFormulaListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseFormulaListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseFormulaListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseFormulaListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStart is called when production start is entered.
func (s *BaseFormulaListener) EnterStart(ctx *StartContext) {}

// ExitStart is called when production start is exited.
func (s *BaseFormulaListener) ExitStart(ctx *StartContext) {}

// EnterBinary is called when production Binary is entered.
func (s *BaseFormulaListener) EnterBinary(ctx *BinaryContext) {}

// ExitBinary is called when production Binary is exited.
func (s *BaseFormulaListener) ExitBinary(ctx *BinaryContext) {}

// EnterNegation is called when production Negation is entered.
func (s *BaseFormulaListener) EnterNegation(ctx *NegationContext) {}

// ExitNegation is called when production Negation is exited.
func (s *BaseFormulaListener) ExitNegation(ctx *NegationContext) {}

// EnterLetter is called when production Letter is entered.
func (s *BaseFormulaListener) EnterLetter(ctx *LetterContext) {}

// ExitLetter is called when production Letter is exited.
func (s *BaseFormulaListener) ExitLetter(ctx *LetterContext) {}
