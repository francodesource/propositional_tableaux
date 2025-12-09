package formula

import "testing"

var p = NewLetter("p")
var q = NewLetter("q")

func TestParseFormula(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   Formula
	}{
		{
			name:   "letter",
			source: "p",
			want:   NewLetter("p"),
		},
		{
			name:   "literal",
			source: "!p",
			want:   NewNot(NewLetter("p")),
		},
		{
			name:   "and",
			source: "(p & q)",
			want:   NewAnd(NewLetter("p"), NewLetter("q")),
		},
		{
			name:   "or",
			source: "(p | q)",
			want:   NewOr(p, q),
		},
		{
			name:   "implication",
			source: "(p -> q)",
			want:   NewImplies(p, q),
		},
		{
			name:   "biconditional",
			source: "(p <-> q)",
			want:   NewBiconditional(p, q),
		},
		{
			name:   "nand",
			source: "(p !& q)",
			want:   NewNand(p, q),
		},
		{
			name:   "nor",
			source: "(p !| q)",
			want:   NewNor(p, q),
		},
		{
			name:   "xor",
			source: "(p ^ q)",
			want:   NewXor(p, q),
		},
		{
			name:   "complex formula",
			source: "!((a !& !b) ^ (a <-> !(b | c)))",
			want: NewNot(NewXor(
				NewNand(NewLetter("a"), NewNot(NewLetter("b"))),
				NewBiconditional(NewLetter("a"), NewNot(NewOr(NewLetter("b"), NewLetter("c")))))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.source)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
