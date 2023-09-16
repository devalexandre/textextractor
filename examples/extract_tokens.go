package main

import (
	"fmt"

	nlpk "github.com/devalexandre/nlpk/pkg"
)

func main() {
	ner := nlpk.NewNLPK()
	input := "Name 6: {Name}. DOB: {DOB}."
	tokens := ner.ExtractTokens(input)
	fmt.Println("Extracted Tokens:", tokens)
	//output: Extracted Tokens: [Name DOB]
}
