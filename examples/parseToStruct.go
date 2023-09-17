package main

import textextractor "github.com/devalexandre/textextractor/pkg"

func main() {
	extractAndSave()
	p := textextractor.NewTextExtractor()
	input := "Name 6: ABDUL BASEER 1: ABDUL QADEER 2: BASIR 3: n/a 4: n/a 5: n/a.\nName (non-Latin script): عبدالقدیر بصیر عبد البصير\nTitle: (1) General (2) Maulavi DOB: --/--/1964. POB: (1) Hisarak District, Nangarhar Province. (2) Surkh Rod District, Nangarhar Province, (1) Afghanistan (2) Afghanistan Good quality a.k.a: (1) BASIR, Abdul Qadir (2) HAJI, Ahmad (3) HAQQANI, Abdul Qadir (4) QADIR, Abdul  Nationality: Afghanistan Passport Number: D 000974 Passport Details: Afghanistan number Position: (1) Head of Taliban Peshawar Financial Commission. (2) Military Attache, Taliban Embassy, Islamabad, Pakistan Other Information: (UK Sanctions List Ref):AFG0098. (UN Ref):TAi.128. Financial advisor to Taliban Peshawar Military Council and Head of Taliban Peshawar Financial Commission. Believed to be in Afghanistan/Pakistan border area. Review pursuant to Security Council resolution 1822 (2008) was concluded on 21 Jul. 2010. INTERPOL-UN Security Council Special Notice web link: https://www.interpol.int/en/How-we-work/Notices/View-UN-Notices-Individuals click here Listed on: 23/02/2001 UK Sanctions List Date Designated: 25/01/2001 Last Updated: 01/02/2021 Group ID: 6911"

	type Entity struct {
		Name        string `data:"NAME"`
		DOB         string `data:"DOB"`
		Title       string `data:"TITLE"`
		AKA         string `data:"AKA"`
		Nationality string `data:"NATIONALITY"`
	}

	entity := Entity{}
	ok := p.ParseValueToStruct(input, &entity, "model.json")
	if !ok {
		panic("error")
	}
}
