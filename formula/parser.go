package formula

import (
	"github.com/antlr4-go/antlr/v4"
	"propositional_tableaux/formula/parser"
)

type formulaListener struct {
	*parser.BaseFormulaListener
	stack []Formula
}

func (f *formulaListener) pop() Formula {
	n := len(f.stack)
	elem := f.stack[n-1]
	f.stack = f.stack[:n-1]
	return elem
}

func (f *formulaListener) ExitLetter(ctx *parser.LetterContext) {
	f.stack = append(f.stack, Letter{name: ctx.GetText()})
}

func (f *formulaListener) ExitNegation(ctx *parser.NegationContext) {
	f.stack = append(f.stack, Not{negated: f.pop()})
}

func (f *formulaListener) ExitBinary(ctx *parser.BinaryContext) {
	right := f.pop()
	left := f.pop()
	var formula Formula
	switch ctx.GetOp().GetTokenType() {
	case parser.FormulaParserAND:
		formula = NewAnd(left, right)
	case parser.FormulaParserOR:
		formula = NewOr(left, right)
	case parser.FormulaParserIMPLIES:
		formula = NewImplies(left, right)
	case parser.FormulaLexerBICONDITIONAL:
		formula = NewBiconditional(left, right)
	case parser.FormulaLexerNAND:
		formula = NewNand(left, right)
	case parser.FormulaLexerNOR:
		formula = NewNor(left, right)
	case parser.FormulaLexerXOR:
		formula = NewXor(left, right)
	}
	f.stack = append(f.stack, formula)
}

func Parse(input string) Formula {
	is := antlr.NewInputStream(input)
	lexer := parser.NewFormulaLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewFormulaParser(stream)
	listener := &formulaListener{}
	antlr.ParseTreeWalkerDefault.Walk(listener, p.Start_())
	return listener.pop()
}
