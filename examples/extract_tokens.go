package main

import (
	"fmt"
	textextractor "github.com/devalexandre/textextractor/pkg"
)

func main() {
	p := textextractor.NewTextExtractor()
	input := "Name 6: {Name}. DOB: {DOB}."
	tokens := p.ExtractTokens(input)
	fmt.Println("Extracted Tokens:", tokens)
	//output: Extracted Tokens: [Name DOB]
}
