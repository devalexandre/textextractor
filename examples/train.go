package main

import (
	"fmt"
	"github.com/devalexandre/textextractor/pkg"
)

func extractAndSave() {

	p := textextractor.NewTextExtractor()

	text := []string{
		"Hello {USER}, sorry for the delay. Since KSQL aims at being simple it avoids as much as possible to hide the actual SQL that is being generated.\n\nYou should be able to use a query builder to get this behavior you wanted.\n\nThat said I think that keeping this WHERE deleted_at IS NULL clause might actually make the query more readable, since a reader that is not aware of this default value could get very confused.\n\nBut that's just my personal preference.",
	}

	tk := p.Learn(text)

	if len(tk) == 0 {
		panic("error")
	}

	p.Save(tk, "tokens")

	fmt.Println("Tokens", tk)
}
