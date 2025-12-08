package tableaux

import (
	"encoding/json"
	"fmt"
	"maps"
	"math/big"
	"math/rand"
	"propositional_tableaux/formula"
	tu "propositional_tableaux/internal/testutil"
	"reflect"
	"slices"
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

func generateFormula(rand *rand.Rand, size int) formula.Formula {
	if size == 0 {
		letters := "pqrstuvwxyz"
		randLetter := letters[rand.Intn(len(letters))]
		return formula.NewLetter(string(randLetter))
	}

	negation := rand.Intn(2) == 1

	if negation {
		return formula.NewNot(generateFormula(rand, size-1))
	} else {
		return formula.NewBinary(generateFormula(rand, size-1), generateFormula(rand, size-1),
			formula.Operator(rand.Intn(7))) // random binary operation
	}
}

// TestBuildSemanticTableaux3 checks that the assignments discovered by the semantic tableaux are the same as
// the one obtained by calculating all the truth-tables
func TestBuildSemanticTableaux3(t *testing.T) {
	f := func(f formula.Formula) bool {
		bfSat := bruteForceSat(f)
		tab := BuildSemanticTableaux(f)
		tabSat := len(tab.Eval()) > 0

		return bfSat == tabSat
	}

	maxSize := 12 // this is reasonable
	config := &quick.Config{
		MaxCount: 10,
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = reflect.ValueOf(generateFormula(r, r.Intn(maxSize)))
		},
	}

	if err := quick.Check(f, config); err != nil {
		t.Error(err)
	}
}

// TestBuildSemanticTableaux4 checks that the assignments obtained by a semantic tableaux, satisfy the formula.
func TestBuildSemanticTableaux4(t *testing.T) {
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

	maxSize := 12

	config := &quick.Config{
		MaxCount: 20,
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = reflect.ValueOf(generateFormula(r, r.Intn(maxSize)))
		},
	}

	if err := quick.Check(f, config); err != nil {
		t.Error(err)
	}
}

// TestBuildAnalyticTableaux checks that the assignments discovered by the analytic tableaux are the same of the semantic one
func TestBuildAnalyticTableaux(t *testing.T) {
	f := func(f formula.Formula) bool {
		semanticTab := BuildSemanticTableaux(f)
		analyticTab := BuildAnalyticTableaux(f)

		sAssignments := semanticTab.Eval()
		aAssignments := analyticTab.Eval()

		if len(sAssignments) != len(aAssignments) {
			t.Errorf("fail on formula: %v", f)
			return false
		}

		slice1 := make([]string, len(sAssignments))
		slice2 := make([]string, len(aAssignments))

		for i := range sAssignments {
			bytes1, _ := json.Marshal(sAssignments[i])
			bytes2, _ := json.Marshal(aAssignments[i])

			slice1[i] = string(bytes1)
			slice2[i] = string(bytes2)
		}
		slices.Sort(slice1)
		slices.Sort(slice2)
		res := slices.Compare(slice1, slice2) == 0
		if !res {
			t.Errorf("fail on formula: %v", f)
		}
		return res
	}
	maxSize := 10
	config := &quick.Config{
		MaxCount: 30,
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = reflect.ValueOf(generateFormula(r, r.Intn(maxSize)))
		},
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("%v", err)
	}
}

// TestBuildAnalyticTableaux4 checks that the assignments obtained by an analytic tableaux, satisfy the formula.
func TestBuildAnalyticTableaux2(t *testing.T) {
	f := func(f formula.Formula) bool {
		tab := BuildAnalyticTableaux(f)
		assignments := tab.Eval()

		for _, a := range assignments {
			if !evaluate(f, a) {
				fmt.Printf("test fail for formula %v with assignment %v", f, a)
				return false
			}
		}

		return true
	}

	maxSize := 10

	config := &quick.Config{
		MaxCount: 20,
		Values: func(values []reflect.Value, r *rand.Rand) {
			values[0] = reflect.ValueOf(generateFormula(r, r.Intn(maxSize)))
		},
	}

	if err := quick.Check(f, config); err != nil {
		t.Error(err)
	}
}
