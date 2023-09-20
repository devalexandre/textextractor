package main

import textextractor "github.com/devalexandre/textextractor/pkg"

func main() {
	extractAndSave()
	p := textextractor.NewTextExtractor()
	input := "Hello @olivatooo, sorry for the delay. Since KSQL aims at being simple it avoids as much as possible to hide the actual SQL that is being generated.\n\nYou should be able to use a query builder to get this behavior you wanted.\n\nThat said I think that keeping this WHERE deleted_at IS NULL clause might actually make the query more readable, since a reader that is not aware of this default value could get very confused.\n\nBut that's just my personal preference."
	type Entity struct {
		User string `data:"USER"`
	}

	entity := Entity{}
	if err := p.ParseValueToStruct(input, &entity, "tokens"); err != nil {
		panic(err)
	}

}
