package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pevans/banf/blog"
	"github.com/pevans/banf/bnf"
	"github.com/rotisserie/eris"
)

var grammarFile = flag.String("g", "",
	"a file that defines a set of rules for a BNF grammar")

var verboseFlag = flag.Bool("v", false,
	"provide verbose description of what's happening")

func main() {
	flag.Parse()

	if *grammarFile == "" {
		blog.Fail(eris.New("a grammar file must be supplied"))
	}

	if *verboseFlag {
		blog.Verbosity++
	}

	gfile, err := os.Open(*grammarFile)
	if err != nil {
		blog.Fail(eris.Wrap(err, "couldn't open grammar file"))
	}

	stream, perr := bnf.Tokenize(gfile)
	if perr != nil {
		perr.File = *grammarFile
		perr.Err = eris.Wrap(perr.Err, "couldn't tokenize grammar")
		blog.Fail(perr)
	}

	gram, err := bnf.NewGrammar(stream)
	if err != nil {
		blog.Fail(eris.Wrap(err, "couldn't define grammar"))
	}

	for _, infile := range flag.Args() {
		bytes, err := os.ReadFile(infile)
		if err != nil {
			blog.Fail(eris.Wrapf(err, "couldn't read file %s", infile))
		}

		perr := gram.Match(string(bytes))
		if perr != nil && perr.Err != nil {
			perr.File = infile
			perr.Err = eris.Wrapf(perr.Err, "match attempt for %s failed with error", infile)
			blog.Warn(perr)
		}

		if perr == nil {
			fmt.Printf("%s matches\n", infile)
		} else {
			fmt.Printf("%s does not match\n", infile)
			blog.Warn(perr)
		}
	}
}
