# Propositional Tableaux

This program provides some functions to build and use Semantic and Analytic Tableaux.
Three kind of tableau can be found in this module:
- `SemanticNode`
- `AnalyticNode`
- `BufferNode`

All the implementations about the tableaux are collected in the `tableaux` package.
This package provides a common interface called `Node` that is implemented by all three kinds of tableaux.

A tableau can be built from a formula with its specific constructor:
- `BuildSemanticTableaux` for semantic tableau;
- `BuildAnalyticTableaux` for basic analytic tableau;
- `BuildBufferTableaux` for the analytic tableau that uses the buffer.

Then the function `Eval` can be used to produce a slice of `Assignment`, which is a map from a string that identify a letter, to a bool which identify the value of the assignment.

Since one of the most important property of tableau calculus is to produce human-readable proves, 
This module provides three ways to visualize a tableau: 
- a very basic one, which is the default string representation. This represents the tree as nested objects that are enclosed by curly brackets;
- a more readable one that can produce a tree drawn with ascii characters. In this version the formulas can be represented with just ascii characters or with Unicode characters, this function is really slow for bigger formulas;
- a latex string that compiles to a tree using forest package.

Here are some examples with code:

```go
f := formula.Parse("((p | !q) & !q)")
t := tableaux.BuildSemanticTableaux(f)

	fmt.Println(t)

	assignments := t.Eval()
	fmt.Println(assignments)
```
Which produces the following output:
```plaintext
{
values: { literals: {}, alpha: {((p | !q) & !q)}, beta: {} }
left:    {
values: { literals: {!q}, alpha: {}, beta: {(p | !q)} }
left:    {
values: { literals: {p, !q}, alpha: {}, beta: {} }
mark: Open
}
right:    {
values: { literals: {!q}, alpha: {}, beta: {} }
mark: Open
}
}
}
```
This is the default print function of a tableau.
And the output of the assignments:
```plaintext
[map[q:false]]
```
As we can see the assignments are cleaned of redundant elements: the tableau discovers two assignments where one of them gives a value to $p$,
but the value of $p$ does not matter for formula satisfiability.

Since the first way is not very easy to read we can change it to:
```go
fmt.Println(tableaux.UnicodeAsciiTree(t))
```
output:
```plaintext
╭─────────────────╮
│{((p ∨ ¬q) ∧ ¬q)}│
╰────────┬────────╯
         │         
 ╭───────┴──────╮  
 │{¬q, (p ∨ ¬q)}│  
 ╰───────┬──────╯  
     ╭───┴────╮    
 ╭───┴───╮ ╭──┴─╮  
 │{¬q, p}│ │{¬q}│  
 │-------│ │----│  
 │   ○   │ │  ○ │  
 ╰───────╯ ╰────╯  
```

Finally, the last representation uses LaTeX.
```go
fmt.Println(tableaux.TexForestTree(t))
```
This will print the following LaTeX string:

```latex
\begin{forest}
    for tree={
        anchor=north   
    }
    [{$\left\{\left(\left(p \lor \neg q\right) \land \neg q\right)\right\}$}
        [{$\left\{\neg q, \left(p \lor \neg q\right)\right\}$}
            [\shortstack{{$\left\{\neg q, p\right\}$}\\$\odot$}]
            [\shortstack{{$\left\{\neg q\right\}$}\\$\odot$}]
        ]
    ]
\end{forest}
```

This code can be compiled with latex to obtain a pdf representation of the tableau.

## Command line interface
The software provides a command line interface for visualizing tableaux.
The user can call the program with different flags for different options.
- `type` select the type of the tableau, it can be either semantic or analytic;
- `format` select the format of the tableau: 
  - `default` the default format as nested objects;
  - `ascii-tree` print the tableaux as an ascii tree, using ascii characters to write formulas;
  - `ascii-tree-unicode` print the tableaux as an ascii tree, using Unicode characters to write formulas;
  - `tex-forest` print the tableaux as a LaTeX string that can compile into a forest package tree;
- `in` define the input file path from where the formula can be read;
- `out` define an output file path where the tableau will be written.

An example of usage is:
```bash
proptab -type=semantic -format=ascii-tree-unicode -in=path/to/input -out=path/to/output
```
Every flag can be omitted. If `in` is omitted, the user will be asked to insert the formula from `stdin`.
If `out` is omitted the tableau will be printed on `stdout`.

The syntax used for formulas must follow the grammar defined in [Formula.g4](https://github.com/francodesource/propositional_tableaux/blob/master/formula/Formula.g4)

## Installation
The module can be imported in a project via:
```bash
go get github.com/francodesource/propositional_tableaux
```
### Command line interface
The program can be simply installed via `go install`:
```bash
go install github.com/francodesource/propositional_tableaux/cmd/proptab@latest
```

or it can be built from source code by running:
```bash
go build .
```
in the project folder.
