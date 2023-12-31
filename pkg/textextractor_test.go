package textextractor_test

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	textextractor "github.com/devalexandre/textextractor/pkg"

	"testing"
)

func TestSave(t *testing.T) {
	t.Run("learn and save", func(t *testing.T) {
		p := textextractor.NewTextExtractor()
		dataTrain := []string{
			"play to {MUSIC}",
			"play the {MUSIC}",
			"play {MUSIC}",
			"play {MUSIC} for me",
			"play {MUSIC} please",
			"play {MUSIC} now",
		}

		ln := p.Learn(dataTrain)

		if len(ln) <= 4 {
			t.Errorf("got %v want %v", len(ln), 10)
		}

		err := p.Save(ln, "tokens")
		if err != nil {
			t.Errorf("got %v want %v", err, nil)
		}

		//verify if file exists
		_, err = os.Stat("models/tokens.gob")
		if err != nil {
			t.Errorf("got %v want %v", err, nil)
		}

	})

}

// test Load
func TestLoad(t *testing.T) {
	t.Run("load", func(t *testing.T) {
		p := textextractor.NewTextExtractor()
		model, err := p.Load("tokens")
		if err != nil {
			t.Errorf("got %v want %v", err, nil)
		}
		if len(model) <= 4 {
			t.Errorf("got %v want %v", len(model), 10)
		}
	})
}

func TestExtractTokens(t *testing.T) {
	t.Run("extract tokens", func(t *testing.T) {
		p := textextractor.NewTextExtractor()
		input := "Name 6: {Name}. Name (non-Latin script): {NameNonLatin}. DOB: {DOB}. POB: {POB} a.k.a: {GoodQualityAKA}  Other Information: {OtherInformation} Listed on: {Listed} UK Sanctions List Date Designated: 04/10/2011 Last Updated: 01/02/2021 Group ID: 12156."
		want := []string{"Name", "NameNonLatin", "DOB", "POB", "GoodQualityAKA", "OtherInformation", "Listed"}
		got := p.ExtractTokens(input)

		fmt.Println(got)
		fmt.Println(want)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestGetBeforeToken(t *testing.T) {

	t.Run("get word before token", func(t *testing.T) {
		p := textextractor.NewTextExtractor()
		input := "Name 6: {Name}. Brazil"
		want := "e 6: "
		token := p.ExtractTokens(input)[0]
		got := p.GetBeforeToken(input, fmt.Sprintf("{%s}", token))
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestGetAfterToken(t *testing.T) {

	t.Run("get word after token", func(t *testing.T) {
		p := textextractor.NewTextExtractor()
		input := "Name 6: {Name}. Brazil"
		want := ". Bra"
		token := p.ExtractTokens(input)[0]
		got := p.GetAfterToken(input, fmt.Sprintf("{%s}", token))
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestLearnSaveLoad(t *testing.T) {
	extractor := textextractor.NewTextExtractor()
	trainingData := []string{"Play {Song}"}

	// Test Learn
	learnedTokens := extractor.Learn(trainingData)
	if len(learnedTokens) == 0 {
		t.Errorf("Learn() = %v, want at least one token", len(learnedTokens))
	}

	// Test Save
	err := extractor.Save(learnedTokens, "test_tokens")
	if err != nil {
		t.Errorf("Save() error = %v", err)
	}

	// Test Load
	loadedTokens, err := extractor.Load("test_tokens")
	if err != nil {
		t.Errorf("Load() error = %v", err)
	}
	if len(loadedTokens) == 0 {
		t.Errorf("Load() = %v, want at least one token", len(loadedTokens))
	}

	// Cleanup
	os.Remove("models/test_tokens.gob")
}

func TestGetValueBetweenTokens(t *testing.T) {
	t.Run("get value between tokens", func(t *testing.T) {
		p := textextractor.NewTextExtractor()
		p.Weights = textextractor.PrecisionWeights{
			WordLengthWeight:     0.4,
			TokenLengthWeight:    0.3,
			CharacterCountWeight: 0.3,
		}

		trainWord := `Name 6: {NAME}. Name (non-Latin script): عبد العزيز عباسین DOB: {DOB}`
		tokens := p.ExtractTokens(trainWord)

		input := `Name 6: ABBASIN 1: ABDUL AZIZ 2: n/a 3: n/a 4: n/a 5: n/a. Name (non-Latin script): عبد العزيز عباسین DOB: --/--/1969`

		want := map[string]string{
			"NAME": "ABBASIN 1: ABDUL AZIZ 2: n/a 3: n/a 4: n/a 5: n/a",
			"DOB":  "--/--/1969",
		}

		for _, token := range tokens {
			wordBefore := p.GetBeforeToken(trainWord, fmt.Sprintf("{%s}", token))
			wordAfter := p.GetAfterToken(trainWord, fmt.Sprintf("{%s}", token))
			train := textextractor.TokenTrain{
				Name:       token,
				WordBefore: wordBefore,
				WordAfter:  wordAfter,
			}
			got, have := p.GetValueBetweenTokens(input, train, p.Weights)

			if !have {
				t.Errorf("expected to have value for token %s, but didn't", token)
			}

			// Removendo espaços em branco e comparando
			gotValueTrimmed := strings.TrimSpace(got.Value)
			wantValueTrimmed := strings.TrimSpace(want[token])
			if gotValueTrimmed != wantValueTrimmed {
				t.Errorf("for token %s, got value %q, want %q", token, gotValueTrimmed, wantValueTrimmed)
			}
		}
	})
}

func TestLearn(t *testing.T) {
	t.Run("learn", func(t *testing.T) {
		p := textextractor.NewTextExtractor()
		dataTrain := []string{
			`Name 6: {NAME}.
			Name (non-Latin script): عبد العزيز عباسین
			DOB:{DOB}. POB: {POB} Good quality a.k.a: {AKA}  Other Information: (UK Sanctions List Ref):AFG0121. (UN Ref):TAi.155. Key commander in the Haqqani Network (TAe.012) under Sirajuddin Jallaloudine Haqqani (TAi.144). Taliban Shadow Governor for Orgun District, Paktika Province as of early 2010. Operated a training camp for nonAfghan fighters in Paktika Province. Has been involved in the transport of weapons to Afghanistan. INTERPOLUN Security Council Special Notice web link: https://www.interpol.int/en/How-we-work/Notices/View-UN-NoticesIndividuals click here Listed on: 21/10/2011 UK Sanctions List Date Designated: 04/10/2011 Last Updated: 01/02/2021 Group ID: 12156.`,
			`Name 6: {NAME}
			Name (non-Latin script): عزیز الرحمان عبد الاحد
			Title: {TITLE} DOB: {DOB}. POB: {POB} Nationality: Afghanistan National Identification Number: 44323 National Identification Details: (Afghan) (tazkira) Position: Third Secretary, Taliban Embassy, Abu Dhabi, United Arab Emirates Other Information: (UK Sanctions List Ref):AFG0094. (UN Ref):TAi.121. Belongs to Hotak tribe. Review pursuant to Security Council resolution 1822 (2008) was concluded on 29 Jul. 2010. INTERPOL-UN Security Council Special Notice web link: https://www.interpol.int/en/How-we-work/ Notices/View-UN-Notices-Individuals click here Listed on: 23/02/2001 UK Sanctions List Date Designated: 25/01/2001 Last Updated: 01/02/2021 Group ID: 7055.`,
			`Name 6: {NAME}
			Name (non-Latin script): عبدالغنی برادر عبد الاحمد ترک
			Title: {TITLE} DOB: {DOB}. POB: {POB} Good quality a.k.a:{AKA}  Nationality: Afghanistan Position: Deputy Minister of Defence under the Taliban regime Other Information: (UK Sanctions List Ref):AFG0024. (UN Ref):TAi.024. Arrested in Feb. 2010 and in custody in Pakistan. Extradition request to Afghanistan pending in Lahore High Court, Pakistan as of June 2011. Belongs to Popalzai tribe. Senior Taliban military commander and member of Taliban Quetta Council as of May 2007. Review pursuant to Security Council resolution 1822 (2008) was concluded on 1 Jun. 2010. INTERPOL-UN Security Council Special Notice web link: https://www.interpol.int/en/How-we-work/Notices/View-UN-Notices-Individuals click here Listed on: 02/04/2001 UK Sanctions List Date Designated: 23/02/2001 Last Updated: 01/02/2021 Group ID: 7060.`,
		}

		ln := p.Learn(dataTrain)

		if len(ln) <= 10 {
			t.Errorf("got %v want %v", len(ln), 10)
		}

		err := p.Save(ln, "tokens_names")
		if err != nil {
			t.Errorf("got %v want %v", err, nil)
		}
	})

}

func TestParseValueToStruct(t *testing.T) {
	t.Run("parse value to struct", func(t *testing.T) {
		p := textextractor.NewTextExtractor()
		input := "Name 6: ABBASIN 1: ABDUL AZIZ 2: n/a 3: n/a 4: n/a 5: n/a.\nName (non-Latin script): عبد العزيز عباسین\nDOB: --/--/1969. POB: Sheykhan Village, Pirkowti Area, Orgun District, Paktika Province, Afghanistan Good quality a.k.a: MAHSUD, Abdul Aziz  Other Information: (UK Sanctions List Ref):AFG0121. (UN Ref):TAi.155. Key commander in the Haqqani Network (TAe.012) under Sirajuddin Jallaloudine Haqqani (TAi.144). Taliban Shadow Governor for Orgun District, Paktika Province as of early 2010. Operated a training camp for nonAfghan fighters in Paktika Province. Has been involved in the transport of weapons to Afghanistan. INTERPOLUN Security Council Special Notice web link: https://www.interpol.int/en/How-we-work/Notices/View-UN-NoticesIndividuals click here Listed on: 21/10/2011 UK Sanctions List Date Designated: 04/10/2011 Last Updated: 01/02/2021 Group ID: 12156."

		want := map[string]string{
			"NAME": "ABBASIN 1: ABDUL AZIZ 2: n/a 3: n/a 4: n/a 5: n/a.",
			"DOB":  "--/--/1969",
		}

		type Person struct {
			Name string `data:"NAME"`
			DOB  string `data:"DOB"`
		}
		person := Person{}
		err := p.ParseValueToStruct(input, &person, "tokens_names")

		if err != nil {
			t.Errorf("got %v want %v", err, true)
		}

		if person.Name != want["NAME"] {
			t.Errorf("got %v want %v", person.Name, want["NAME"])
		}

		if person.DOB != want["DOB"] {
			t.Errorf("got %v want %v", person.DOB, want["DOB"])
		}

	})
}

func TestGenerateRegex(t *testing.T) {
	extractor := textextractor.NewTextExtractor()
	tokens := []string{"Name", "DOB"}
	expectedRegex := []string{`(?P<Name>[^\s]+)`, `(?P<DOB>[^\s]+)`}

	regexPatterns := extractor.GenerateRegex(tokens)
	if !reflect.DeepEqual(regexPatterns, expectedRegex) {
		t.Errorf("GenerateRegex() = %v, want %v", regexPatterns, expectedRegex)
	}
}

func TestGetModelsDir(t *testing.T) {
	extractor := textextractor.NewTextExtractor()
	modelsDir, err := extractor.GetModelsDir()
	if err != nil {
		t.Errorf("GetModelsDir() error = %v", err)
	}
	if modelsDir == "" {
		t.Errorf("GetModelsDir() returned an empty string")
	}
}

func TestExtractTokensError(t *testing.T) {
	extractor := textextractor.NewTextExtractor()
	input := "Name 6: {Name Name (non-Latin script): {NameNonLatin DOB: {DOB POB: {POB a.k.a: {GoodQualityAKA Other Information: {OtherInformation Listed on: {Listed UK Sanctions List Date Designated: 04/10/2011 Last Updated: 01/02/2021 Group ID: 12156."

	got := extractor.ExtractTokens(input)
	if len(got) != 0 {
		t.Errorf("ExtractTokens() with malformed input = %v, want %v", got, []string{})
	}
}

func TestLoadError(t *testing.T) {
	extractor := textextractor.NewTextExtractor()
	_, err := extractor.Load("non_existent_file.gob")
	if err == nil {
		t.Errorf("Load() with non-existent file, want error")
	}
}

func TestParseValueToStructError(t *testing.T) {
	extractor := textextractor.NewTextExtractor()
	input := "Name: John Doe, Age: 30"

	type Person struct {
		Name string `data:"Name"`
		Age  string `data:"Age"`
	}
	var person Person

	err := extractor.ParseValueToStruct(input, &person, "non_existent_model")
	if err == nil {
		t.Errorf("ParseValueToStruct() with non-existent model file, want error")
	}
}
