package nlpk

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestExtractTokens(t *testing.T) {
	t.Run("extract tokens", func(t *testing.T) {
		nlpk := NewNLPK()
		input := "Name 6: {Name}. Name (non-Latin script): {NameNonLatin}. DOB: {DOB}. POB: {POB} a.k.a: {GoodQualityAKA}  Other Information: {OtherInformation} Listed on: {Listed} UK Sanctions List Date Designated: 04/10/2011 Last Updated: 01/02/2021 Group ID: 12156."
		want := []string{"Name", "NameNonLatin", "DOB", "POB", "GoodQualityAKA", "OtherInformation", "Listed"}
		got := nlpk.ExtractTokens(input)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestGetBeforeToken(t *testing.T) {

	t.Run("get word before token", func(t *testing.T) {
		nlpk := NewNLPK()
		input := "Name 6: {Name}. Brazil"
		want := "e 6: "
		token := nlpk.ExtractTokens(input)[0]
		got := nlpk.GetBeforeToken(input, fmt.Sprintf("{%s}", token))
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestGetAfterToken(t *testing.T) {

	t.Run("get word after token", func(t *testing.T) {
		nlpk := NewNLPK()
		input := "Name 6: {Name}. Brazil"
		want := ". Bra"
		token := nlpk.ExtractTokens(input)[0]
		got := nlpk.GetAfterToken(input, fmt.Sprintf("{%s}", token))
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestGetValueBetweenTokens(t *testing.T) {

	t.Run("get value between tokens", func(t *testing.T) {
		nlpk := NewNLPK()
		trainWord := `Name 6: {NAME}.
		Name (non-Latin script): عبد العزيز عباسین
		DOB: {DOB}`

		token := nlpk.ExtractTokens(trainWord)

		input := `Name 6: ABBASIN 1: ABDUL AZIZ 2: n/a 3: n/a 4: n/a 5: n/a.
		Name (non-Latin script): عبد العزيز عباسین
		DOB: --/--/1969`

		want := map[string]string{
			"NAME": "ABBASIN 1: ABDUL AZIZ 2: n/a 3: n/a 4: n/a 5: n/a.",
			"DOB":  "--/--/1969",
		}

		for _, token := range token {
			wordBefore := nlpk.GetBeforeToken(trainWord, fmt.Sprintf("{%s}", token))
			wordAfter := nlpk.GetAfterToken(trainWord, fmt.Sprintf("{%s}", token))
			train := TokenTrain{
				Name:       token,
				WordBefore: wordBefore,
				WordAfter:  wordAfter,
			}
			got, have := nlpk.GetValueBetweenTokens(input, train)

			if !have {
				t.Errorf("got %v want %v", have, true)
			}
			if got.Value != want[token] {
				t.Errorf("got %v want %v", got.Value, want[token])
			}
		}
	})
}

func TestParseValueToStruct(t *testing.T) {
	t.Run("parse value to struct", func(t *testing.T) {
		nlpk := NewNLPK()
		input := `Name 6: ABBASIN 1: ABDUL AZIZ 2: n/a 3: n/a 4: n/a 5: n/a.
	Name (non-Latin script): عبد العزيز عباسین
	DOB: --/--/1969`

		want := map[string]string{
			"NAME": "ABBASIN 1: ABDUL AZIZ 2: n/a 3: n/a 4: n/a 5: n/a.",
			"DOB":  "",
		}

		type Person struct {
			Name string `data:"NAME"`
			DOB  string `data:"DOB"`
		}
		person := Person{}
		ok := nlpk.ParseValueToStruct(input, &person, "tokens.json")

		if !ok {
			t.Errorf("got %v want %v", ok, true)
		}

		if person.Name != want["NAME"] {
			t.Errorf("got %v want %v", person.Name, want["NAME"])
		}

		if person.DOB != want["DOB"] {
			t.Errorf("got %v want %v", person.DOB, want["DOB"])
		}

	})
}

func TestLearn(t *testing.T) {
	t.Run("learn", func(t *testing.T) {
		nlpk := NewNLPK()
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

		ln := nlpk.Learn(dataTrain)

		if len(ln) <= 10 {
			t.Errorf("got %v want %v", len(ln), 10)
		}
	})

}

func TestSave(t *testing.T) {
	t.Run("learn and save", func(t *testing.T) {
		nlpk := NewNLPK()
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

		ln := nlpk.Learn(dataTrain)

		if len(ln) <= 10 {
			t.Errorf("got %v want %v", len(ln), 10)
		}

		err := nlpk.Save(ln, "tokens.json")
		if err != nil {
			t.Errorf("got %v want %v", err, nil)
		}

		//verify if file exists
		_, err = os.Stat("tokens.json")
		if err != nil {
			t.Errorf("got %v want %v", err, nil)
		}

	})

}

// test Load
func TestLoad(t *testing.T) {
	t.Run("load", func(t *testing.T) {
		nlpk := NewNLPK()
		model, err := nlpk.Load("tokens.json")
		if err != nil {
			t.Errorf("got %v want %v", err, nil)
		}
		if len(model) <= 10 {
			t.Errorf("got %v want %v", len(model), 10)
		}
	})
}