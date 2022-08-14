package main

import "fmt"

func main() {
	fmt.Println("hello world")

	// broad thoughts
	// - set of rules to match
	// - set of possible tokens that can be matched
	// - parse a reader into a set of those tokens
	// - if we have things that don't match any tokens, error
	// - match tokens against rules
	// - if we have token sequences that don't match any rules, error
	// - produce an AST of the input
	// - all done
}
