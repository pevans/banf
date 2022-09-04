package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pevans/banf/bnf"
)

var grammarFile = flag.String("g", "",
	"a file that defines a set of rules for a BNF grammar")

func main() {
	flag.Parse()

	if *grammarFile == "" {
		fmt.Fprintf(os.Stderr, "a grammar file must be supplied\n")
		os.Exit(1)
	}

	gfile, err := os.Open(*grammarFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to read grammar file\n")
		os.Exit(1)
	}

	stream, err := bnf.Tokenize(gfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't tokenize grammar file: %v\n", err)
		os.Exit(1)
	}

	gram, err := bnf.NewGrammar(stream)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't define grammar: %v\n", err)
		os.Exit(1)
	}

	for _, infile := range flag.Args() {
		bytes, err := os.ReadFile(infile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "couldn't read file %s: %v\n", infile, err)
			os.Exit(1)
		}

		matches, err := gram.Match(string(bytes))
		if err != nil {
			fmt.Fprintf(os.Stderr, "match attempt for %s errored: %v\n", infile, err)
		}

		if matches {
			fmt.Printf("%s matches\n", infile)
		} else {
			fmt.Printf("%s does not match\n", infile)
		}
	}
}
