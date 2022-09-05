package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pevans/banf/bnf"
	"github.com/rotisserie/eris"
)

var grammarFile = flag.String("g", "",
	"a file that defines a set of rules for a BNF grammar")

func main() {
	flag.Parse()

	if *grammarFile == "" {
		fail(eris.New("a grammar file must be supplied"))
	}

	gfile, err := os.Open(*grammarFile)
	if err != nil {
		fail(eris.Wrap(err, "couldn't open grammar file"))
	}

	stream, perr := bnf.Tokenize(gfile)
	if perr != nil {
		perr.File = *grammarFile
		perr.Err = eris.Wrap(perr.Err, "couldn't tokenize grammar")
		fail(perr)
	}

	gram, err := bnf.NewGrammar(stream)
	if err != nil {
		fail(eris.Wrap(err, "couldn't define grammar"))
	}

	for _, infile := range flag.Args() {
		bytes, err := os.ReadFile(infile)
		if err != nil {
			fail(eris.Wrapf(err, "couldn't read file %s", infile))
		}

		perr := gram.Match(string(bytes))
		if perr != nil && perr.Err != nil {
			perr.File = infile
			perr.Err = eris.Wrapf(perr.Err, "match attempt for %s failed with error", infile)
			warn(perr)
		}

		if perr == nil {
			fmt.Printf("%s matches\n", infile)
		} else {
			fmt.Printf("%s does not match\n", infile)
			warn(perr)
		}
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, eris.ToString(err, true))
	os.Exit(1)
}

func warn(err error) {
	fmt.Fprintln(os.Stderr, eris.ToString(err, true))
}
