// Code generated from formula/Formula.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"sync"
	"unicode"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type FormulaLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var FormulaLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func formulalexerLexerInit() {
	staticData := &FormulaLexerLexerStaticData
	staticData.ChannelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.ModeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.LiteralNames = []string{
		"", "'('", "')'", "'&'", "'|'", "'->'", "'<->'", "'!|'", "'!&'", "'^'",
		"'!'",
	}
	staticData.SymbolicNames = []string{
		"", "OP", "CP", "AND", "OR", "IMPLIES", "BICONDITIONAL", "NOR", "NAND",
		"XOR", "NOT", "VARIABLE", "WHITESPACE",
	}
	staticData.RuleNames = []string{
		"OP", "CP", "AND", "OR", "IMPLIES", "BICONDITIONAL", "NOR", "NAND",
		"XOR", "NOT", "VARIABLE", "WHITESPACE",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 12, 62, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 1, 0, 1, 0, 1, 1, 1, 1, 1, 2, 1, 2, 1, 3, 1, 3,
		1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7,
		1, 7, 1, 8, 1, 8, 1, 9, 1, 9, 1, 10, 4, 10, 52, 8, 10, 11, 10, 12, 10,
		53, 1, 11, 4, 11, 57, 8, 11, 11, 11, 12, 11, 58, 1, 11, 1, 11, 0, 0, 12,
		1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8, 17, 9, 19, 10, 21, 11,
		23, 12, 1, 0, 2, 4, 0, 48, 57, 65, 90, 95, 95, 97, 122, 3, 0, 9, 10, 13,
		13, 32, 32, 63, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0,
		7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0,
		0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 0, 21, 1, 0, 0,
		0, 0, 23, 1, 0, 0, 0, 1, 25, 1, 0, 0, 0, 3, 27, 1, 0, 0, 0, 5, 29, 1, 0,
		0, 0, 7, 31, 1, 0, 0, 0, 9, 33, 1, 0, 0, 0, 11, 36, 1, 0, 0, 0, 13, 40,
		1, 0, 0, 0, 15, 43, 1, 0, 0, 0, 17, 46, 1, 0, 0, 0, 19, 48, 1, 0, 0, 0,
		21, 51, 1, 0, 0, 0, 23, 56, 1, 0, 0, 0, 25, 26, 5, 40, 0, 0, 26, 2, 1,
		0, 0, 0, 27, 28, 5, 41, 0, 0, 28, 4, 1, 0, 0, 0, 29, 30, 5, 38, 0, 0, 30,
		6, 1, 0, 0, 0, 31, 32, 5, 124, 0, 0, 32, 8, 1, 0, 0, 0, 33, 34, 5, 45,
		0, 0, 34, 35, 5, 62, 0, 0, 35, 10, 1, 0, 0, 0, 36, 37, 5, 60, 0, 0, 37,
		38, 5, 45, 0, 0, 38, 39, 5, 62, 0, 0, 39, 12, 1, 0, 0, 0, 40, 41, 5, 33,
		0, 0, 41, 42, 5, 124, 0, 0, 42, 14, 1, 0, 0, 0, 43, 44, 5, 33, 0, 0, 44,
		45, 5, 38, 0, 0, 45, 16, 1, 0, 0, 0, 46, 47, 5, 94, 0, 0, 47, 18, 1, 0,
		0, 0, 48, 49, 5, 33, 0, 0, 49, 20, 1, 0, 0, 0, 50, 52, 7, 0, 0, 0, 51,
		50, 1, 0, 0, 0, 52, 53, 1, 0, 0, 0, 53, 51, 1, 0, 0, 0, 53, 54, 1, 0, 0,
		0, 54, 22, 1, 0, 0, 0, 55, 57, 7, 1, 0, 0, 56, 55, 1, 0, 0, 0, 57, 58,
		1, 0, 0, 0, 58, 56, 1, 0, 0, 0, 58, 59, 1, 0, 0, 0, 59, 60, 1, 0, 0, 0,
		60, 61, 6, 11, 0, 0, 61, 24, 1, 0, 0, 0, 3, 0, 53, 58, 1, 6, 0, 0,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// FormulaLexerInit initializes any static state used to implement FormulaLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewFormulaLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func FormulaLexerInit() {
	staticData := &FormulaLexerLexerStaticData
	staticData.once.Do(formulalexerLexerInit)
}

// NewFormulaLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewFormulaLexer(input antlr.CharStream) *FormulaLexer {
	FormulaLexerInit()
	l := new(FormulaLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &FormulaLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "Formula.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// FormulaLexer tokens.
const (
	FormulaLexerOP            = 1
	FormulaLexerCP            = 2
	FormulaLexerAND           = 3
	FormulaLexerOR            = 4
	FormulaLexerIMPLIES       = 5
	FormulaLexerBICONDITIONAL = 6
	FormulaLexerNOR           = 7
	FormulaLexerNAND          = 8
	FormulaLexerXOR           = 9
	FormulaLexerNOT           = 10
	FormulaLexerVARIABLE      = 11
	FormulaLexerWHITESPACE    = 12
)
