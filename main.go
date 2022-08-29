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

	_, err = bnf.NewGrammar(stream)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't define grammar: %v\n", err)
		os.Exit(1)
	}

	/*
		err = gram.Build(stream)
		if err != nil {
			fmt.Fprintf(os.Stderr, "couldn't build grammar: %v\n", err)
			os.Exit(1)
		}

		for _, infile := range flag.Args() {
			err := gram.Validate(infile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "couldn't validate input %s: %v\n", infile, err)
				os.Exit(1)
			}

			fmt.Printf("%s is valid\n", infile)
		}
	*/
}
