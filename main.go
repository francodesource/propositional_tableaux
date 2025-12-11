package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"propositional_tableaux/formula"
	"propositional_tableaux/tableaux"
)

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

type nopReadCloser struct {
	io.Reader
}

func (nopReadCloser) Close() error { return nil }

func checkType(flag string) error {
	if flag != "semantic" && flag != "analytic" {
		return fmt.Errorf("error: invalid value for -type: '%s'", flag)
	}
	return nil
}

func checkFormat(flag string) error {
	if flag != "default" && flag != "ascii-tree" && flag != "ascii-tree-unicode" && flag != "tex-forest" {
		return fmt.Errorf("error: invalid value for -format: '%s'", flag)
	}
	return nil
}

func checkInput(flag string) (io.ReadCloser, error) {
	if flag == "stdin" {
		return nopReadCloser{os.Stdin}, nil
	}

	return os.Open(flag)
}

func checkOutput(flag string) (io.WriteCloser, error) {
	if flag == "stdout" {
		return nopWriteCloser{os.Stdout}, nil
	}

	if _, err := os.Stat(flag); errors.Is(err, os.ErrNotExist) {
		return os.Create(flag)
	}

	return os.Open(flag)
}

func errorHandler(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	flag.Usage()
	os.Exit(2)
}

func try(err error) {
	if err != nil {
		errorHandler(err)
	}
}

func main() {
	//flags
	var (
		tableauxType = flag.String("type", "semantic", " semantic | analytic")
		format       = flag.String("format", "default", "default | ascii-tree | ascii-tree-unicode | tex-forest")
		in           = flag.String("in", "stdin", "the name of the input file")
		out          = flag.String("out", "stdout", "the name of the output file")
	)

	flag.Parse()

	try(checkType(*tableauxType))
	try(checkFormat(*format))

	input, err := checkInput(*in)
	if err != nil {
		errorHandler(err)
	}

	output, err := checkOutput(*out)
	if err != nil {
		errorHandler(err)
	}

	defer func(input io.ReadCloser) {
		err := input.Close()
		if err != nil {
			panic(err)
		}
	}(input)
	defer func(output io.WriteCloser) {
		err := output.Close()
		if err != nil {
			panic(err)
		}
	}(output)

	s := bufio.NewScanner(input)
	var str string
	var printPrompt = input == nopReadCloser{os.Stdin}

	var prompt = ">> "

	if printPrompt {
		fmt.Println("Write a propositional formula, empty input to stop writing")
		fmt.Print(prompt)
	} else {
		fmt.Println("Reading from file")
	}
	for s.Scan() {
		line := s.Text()

		if line == "" {
			break
		}

		str += line

		if printPrompt {
			fmt.Print(prompt)
		}
	}

	f := formula.Parse(str)

	var tab tableaux.Node

	switch *tableauxType {
	case "semantic":
		tab = tableaux.BuildSemanticTableaux(f)
	case "analytic":
		tab = tableaux.BuildBufferedTableaux(f)
	}

	var stringRep string

	switch *format {
	case "default":
		stringRep = fmt.Sprint(tab)
	case "ascii-tree":
		stringRep = fmt.Sprint(tableaux.DefaultAsciiTree(tab))
	case "ascii-tree-unicode":
		stringRep = fmt.Sprint(tableaux.UnicodeAsciiTree(tab))
	case "tex-forest":
		stringRep = tableaux.TexForestTree(tab)
	}

	_, err = output.Write([]byte(stringRep))

	if err != nil {
		panic(err)
	}
}
