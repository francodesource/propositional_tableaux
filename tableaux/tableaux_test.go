package tableaux

import (
	"fmt"
	"maps"
	"math/big"
	"math/rand"
	"propositional_tableaux/formula"
	tu "propositional_tableaux/internal/testutil"
	"reflect"
	"slices"
	"strings"
	"testing"
	"testing/quick"
)

var (
	simpleUnsat = formula.NewAnd(
		tu.P,
		formula.Complement(tu.P),
	)

	simpleTaut = formula.NewOr(
		tu.P,
		formula.Complement(tu.P),
	)

	unsat = formula.NewAnd(
		formula.NewOr(tu.P, tu.Q),
		formula.NewAnd(formula.Complement(tu.P), formula.Complement(tu.Q)),
	)

	sat = formula.NewAnd(
		tu.P,
		formula.NewOr(formula.NewNot(tu.Q), formula.NewNot(tu.P)),
	)
)

func normalizeAssignments(assignments []Assignment) string {
	assignmentsString := make([]string, 0, len(assignments))
	for _, a := range assignments {
		assignmentsString = append(assignmentsString, normalizeAssignment(a))
	}
	slices.Sort(assignmentsString)
	return strings.Join(assignmentsString, " ")
}

func compareAssignments(a1, a2 []Assignment) bool {
	return normalizeAssignments(a1) == normalizeAssignments(a2)
}

// TestBuildSemanticTableaux tests whether the function produces tableaux that leads to correct assignments.
// Since the construction of a tableaux is deterministic it proves just some properties
func TestBuildSemanticTableaux(t *testing.T) {
	tests := []struct {
		name string
		f    formula.Formula
		want []Assignment
	}{
		{
			"simple unsat",
			simpleUnsat,
			[]Assignment{},
		},
		{
			"simple tautology",
			simpleTaut,
			[]Assignment{
				{"P": true},
				{"P": false},
			},
		},
		{
			"unsat",
			unsat,
			[]Assignment{},
		},
		{
			"sat",
			sat,
			[]Assignment{{
				"P": true,
				"Q": false,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tab := BuildSemanticTableaux(tt.f)
			got := tab.Eval()

			if len(got) != len(tt.want) {
				t.Errorf("want %v, got %v", tt.want, got)
			}

			equalCounter := 0
			for _, assignment1 := range tt.want {
				for _, assignment2 := range got {
					if maps.Equal(assignment2, assignment1) {
						equalCounter++
					}
				}
			}

			if equalCounter != len(tt.want) {
				t.Errorf("want %v, got %v", tt.want, got)
			}
		})
	}
}

// TestBuildSemanticTableaux2 checks the improvement of checking for complementary formula.
func TestBuildSemanticTableaux2(t *testing.T) {
	inner := formula.NewNand(
		formula.NewOr(tu.P, formula.NewBiconditional(tu.Q, tu.R)),
		formula.NewBiconditional(tu.P1, formula.NewXor(tu.R, tu.S)))

	f := formula.NewAnd(inner, formula.NewNot(inner))
	tab := BuildSemanticTableaux(f)
	assignment := tab.Eval()

	if len(assignment) > 0 {
		t.Errorf("expected formula to be unasat")
	}

	if tab.Height() != 2 {
		t.Errorf("expected tableaux to be of minimal length 2")
	}
}

// TestBuildSemanticTableaux_Tautologies checks that given a negated tautology the semantic tableaux produces a
// closed tableaux, which means that Eval() produces 0 assignments.
func TestBuildSemanticTableaux_Tautologies(t *testing.T) {
	tests := []struct {
		name string
		f    formula.Formula
	}{
		{
			name: "excluded middle",
			f:    formula.NewOr(tu.P, formula.NewNot(tu.P)),
		},
		{
			name: "contraposition",
			f: formula.NewBiconditional(
				formula.NewImplies(tu.P, tu.Q),
				formula.NewImplies(formula.NewNot(tu.Q), formula.NewNot(tu.P)),
			),
		},
		{
			name: "reductio ad absurdum",
			f: formula.NewImplies(
				formula.NewAnd(
					formula.NewImplies(formula.NewNot(tu.P), tu.Q),
					formula.NewImplies(formula.NewNot(tu.P), formula.NewNot(tu.Q)),
				),
				tu.P,
			),
		},
		{
			name: "De Morgan's law",
			f: formula.NewBiconditional(
				formula.NewNot(formula.NewAnd(tu.P, tu.Q)),
				formula.NewOr(formula.NewNot(tu.P), formula.NewNot(tu.Q)),
			),
		},
		{
			name: "hypothetical syllogism",
			f: formula.NewImplies(formula.NewAnd(
				formula.NewImplies(tu.P, tu.Q),
				formula.NewImplies(tu.Q, tu.R),
			), formula.NewImplies(tu.P, tu.R)),
		},
		{
			name: "proof by cases",
			f: formula.NewImplies(
				formula.NewAnd(
					formula.NewOr(tu.P, tu.R),
					formula.NewAnd(
						formula.NewImplies(tu.P, tu.R),
						formula.NewImplies(tu.Q, tu.R),
					),
				),
				tu.R,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tab := BuildSemanticTableaux(formula.NewNot(tt.f))
			if len(tab.Eval()) > 0 {
				t.Errorf("tableaux found assignments %v for tautology %v", tab.Eval(), tt.f)
			}
		})
	}

}

/*
	Property based testing
*/

func allLetters(f formula.Formula, letters map[string]struct{}) {
	switch f := f.(type) {
	case formula.Letter:
		letters[f.Name()] = struct{}{}
	case formula.Not:
		allLetters(f.Negated(), letters)
	case formula.Binary:
		allLetters(f.Left(), letters)
		allLetters(f.Right(), letters)
	}
}

func evaluate(f formula.Formula, assignment Assignment) bool {
	switch f := f.(type) {
	case formula.Letter:
		return assignment[f.Name()]
	case formula.Not:
		return !evaluate(f.Negated(), assignment)
	case formula.Binary:
		switch f.Op() {
		case formula.And:
			return evaluate(f.Left(), assignment) && evaluate(f.Right(), assignment)
		case formula.Or:
			return evaluate(f.Left(), assignment) || evaluate(f.Right(), assignment)
		case formula.Implies:
			return !evaluate(f.Left(), assignment) || evaluate(f.Right(), assignment)
		case formula.Nand:
			return !(evaluate(f.Left(), assignment) && evaluate(f.Right(), assignment))
		case formula.Nor:
			return !(evaluate(f.Left(), assignment) || evaluate(f.Right(), assignment))
		case formula.Biconditional:
			return evaluate(f.Left(), assignment) == evaluate(f.Right(), assignment)
		case formula.Xor:
			return evaluate(f.Left(), assignment) != evaluate(f.Right(), assignment)
		default:
			panic("unreachable")
		}
	default:
		panic(fmt.Errorf("%T is not a formula", f))
	}
}

func allAssignmentsAux(letters []string, assignments []Assignment) {
	one := big.NewInt(1)
	counter := big.NewInt(0)

	for i := range len(assignments) {
		for j := 0; j < len(letters); j++ {
			assignments[i][letters[j]] = counter.Bit(j) == 1
		}
		counter.Add(counter, one)
	}
}

func allAssignments(f formula.Formula) []Assignment {
	m := make(map[string]struct{})
	allLetters(f, m)
	letters := slices.Collect(maps.Keys(m))

	// initializing all assignments to 2^n
	res := make([]Assignment, 1<<len(letters))
	for i := range res {
		res[i] = make(Assignment, len(letters))
	}

	allAssignmentsAux(letters, res)
	return res
}

func bruteForceSat(f formula.Formula) bool {
	as := allAssignments(f)

	sat := false

	for _, a := range as {
		sat = sat || evaluate(f, a)
	}
	return sat
}

func compareTableauxWithTruthTables(t *testing.T, f formula.Formula, tab Node) bool {
	bfSat := bruteForceSat(f)
	tabSat := len(tab.Eval()) > 0

	res := bfSat == tabSat

	if !res {
		t.Errorf("failed for formula %v: truth-table = %v, tableaux = %v\n", f, bfSat, tabSat)
	}

	return res
}

// TestBuildSemanticTableaux_TruthTables checks that the assignments discovered by the semantic tableaux are the same as
// the one obtained by calculating all the truth-tables
func TestBuildSemanticTableaux_TruthTables(t *testing.T) {
	f := func(f formula.Formula) bool {
		tab := BuildSemanticTableaux(f)

		return compareTableauxWithTruthTables(t, f, tab)
	}

	maxSize := 50
	config := &quick.Config{
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = reflect.ValueOf(formula.GenerateRandom(r, r.Intn(maxSize)))
		},
	}

	if err := quick.Check(f, config); err != nil {
		t.Error(err)
	}
}

// TestBuildSemanticTableaux_Assignments checks that the assignments obtained by a semantic tableaux, satisfy the formula.
func TestBuildSemanticTableaux_Assignments(t *testing.T) {
	f := func(f formula.Formula) bool {
		tab := BuildSemanticTableaux(f)
		assignments := tab.Eval()

		for _, a := range assignments {
			if !evaluate(f, a) {
				fmt.Printf("test fail for formula %v with assignment %v", f, a)
				return false
			}
		}

		return true
	}

	maxSize := 50

	config := &quick.Config{
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = reflect.ValueOf(formula.GenerateRandom(r, r.Intn(maxSize)))
		},
	}

	if err := quick.Check(f, config); err != nil {
		t.Error(err)
	}
}

// testTableauxMarks this function tests that every leaf of the tableaux is either only literals and marks as closed
// or open, or it is mark as closed. The last option is for the case where the algorithm find a formula and its
// complement and so stops.
func testTableauxMarks(t *testing.T, tab Node) bool {
	if tab.IsLeaf() {
		allLiterals := true

		for f := range tab.Formulas() {
			if !formula.IsLiteral(f) {
				allLiterals = false

			}
		}

		return (allLiterals && (tab.IsOpen() || tab.IsClosed())) || tab.IsClosed()
	}

	var left, right = true, true
	if tab.Left() != nil {
		left = testTableauxMarks(t, tab.Left())
	}

	if tab.Right() != nil {
		right = testTableauxMarks(t, tab.Right())
	}

	return left && right
}

func TestSemanticTableaux_Marks(t *testing.T) {
	f := func(f formula.Formula) bool {
		tab := BuildSemanticTableaux(f)

		res := testTableauxMarks(t, tab)

		if !res {
			t.Errorf("Fail for %v\n", f)
		}
		return res
	}

	maxSize := 50

	config := &quick.Config{
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = reflect.ValueOf(formula.GenerateRandom(r, r.Intn(maxSize)))
		},
	}

	if err := quick.Check(f, config); err != nil {
		t.Error(err)
	}
}

/*
	Benchmarks
*/

var bigFormula = formula.Parse("!((!!!(((!!s !| ((y ^ w) <-> !z)) <-> (((q -> x) | (z ^ q)) -> !(t ^ u))) <-> ((((s <-> w) <-> (t !| x)) -> !!p) !| (!(s ^ p) <-> (!w !& !u)))) !& !!!((!(!!u <-> (!q !| (x -> r))) -> !(!!v ^ !(y !& q))) & !!((((t <-> s) !| (s & s)) ^ (!z | !p)) !& !!!!x))) !& !!((!(((!t !| (q ^ y)) <-> !(u -> s)) & !!(!x !& (t & r))) !| !!(((!s | (v !& x)) & !(s -> q)) & (!!t !| ((w ^ t) ^ (w <-> w))))) !& !((!!!(!u <-> !u) ^ (!(t !& t) ^ ((x & v) & !r))) | (!!((q & q) <-> !p) <-> (!(t <-> w) | (!p & (u !| r)))))))")

func BenchmarkBuildSemanticTableaux(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildSemanticTableaux(bigFormula)
	}
}

func BenchmarkBuildAnalyticTableaux(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildAnalyticTableaux(bigFormula)
	}
}

func BenchmarkBuildBufferedTableaux(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildBufferTableaux(bigFormula)
	}
}
