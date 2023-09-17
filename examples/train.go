package main

import (
	"fmt"
	"github.com/devalexandre/textextractor/pkg"
)

func extractAndSave() {

	p := textextractor.NewTextExtractor()

	text := []string{
		"Name 6: ABDUL GHANI 1: ABDUL GHAFAR 2: QURISHI 3: n/a 4: n/a 5: n/a.\nName (non-Latin script): عبدالغفار قریشی عبد الغنی\nTitle: Maulavi DOB:  (1) --/--/1970. (2) --/--/1967. POB: Turshut village, Wursaj District, Takhar Province, Afghanistan Good quality a.k.a: QURESHI, Abdul Ghaffar  Nationality: {NATIONALITY} Passport Number: D 000933 Passport Details: (Afghan) Issued in Kabul on 13 Sep 1998 National Identification Number: 55130 National Identification Details: (Afghan) (tazkira) Address: Khairkhana Section, Number 3, Kabul, Afghanistan.Position: Repatriation Attache, Taliban Embassy, Islamabad, Pakistan Other Information: (UK Sanctions List Ref):AFG0100. (UN Ref):TAi.130. Involved in drug trafficking. Belongs to Tajik ethnic group. Review pursuant to Security Council resolution 1822 (2008) was concluded on 29 Jul. 2010. INTERPOL-UN Security Council Special Notice web link: https://www.interpol.int/en/How-we-work/Notices/View-UN-Notices-Individuals click here Listed on: 23/02/2001 UK Sanctions List Date Designated: 25/01/2001 Last Updated: 01/02/2021 Group ID: 7405.",
		"Name 6: {NAME}.\nName (non-Latin script): عبد العزيز عباسین\nDOB: {DOB}. POB: {POB} Good quality a.k.a: {AKA}  Other Information: (UK Sanctions List Ref):AFG0121. (UN Ref):TAi.155. Key commander in the Haqqani Network (TAe.012) under Sirajuddin Jallaloudine Haqqani (TAi.144). Taliban Shadow Governor for Orgun District, Paktika Province as of early 2010. Operated a training camp for nonAfghan fighters in Paktika Province. Has been involved {CRIME}. INTERPOLUN Security Council Special Notice web link: https://www.interpol.int/en/How-we-work/Notices/View-UN-NoticesIndividuals click here Listed on: 21/10/2011 UK Sanctions List Date Designated: 04/10/2011 Last Updated: 01/02/2021 Group ID: 12156.",
		"Name 6:{NAME}a.\nName (non-Latin script): عزیز الرحمان عبد الاحد\nTitle: {TITLE} DOB: {DOB}. POB: {POB} Nationality: {NATIONALITY} National Identification Number: 44323 National Identification Details: (Afghan) (tazkira) Position: {POSITION}, Taliban Embassy, Abu Dhabi, United Arab Emirates Other Information: (UK Sanctions List Ref):AFG0094. (UN Ref):TAi.121. Belongs to Hotak tribe. Review pursuant to Security Council resolution 1822 (2008) was concluded on 29 Jul. 2010. INTERPOL-UN Security Council Special Notice web link: https://www.interpol.int/en/How-we-work/ Notices/View-UN-Notices-Individuals click here Listed on: 23/02/2001 UK Sanctions List Date Designated: 25/01/2001 Last Updated: 01/02/2021 Group ID: 7055.\n",
		"Name 6:{NAME}.\nName (non-Latin script): عبدالغنی برادر عبد الاحمد ترک\nTitle: {TITLE} DOB: {DOB}. POB: {POB} Good quality a.k.a: {AKA} Nationality: {NATIONALITY} Position: {POSITION} Other Information: (UK Sanctions List Ref):AFG0024. (UN Ref):TAi.024. Arrested in Feb. 2010 and in custody in Pakistan. Extradition request to Afghanistan pending in Lahore High Court, Pakistan as of June 2011. Belongs to Popalzai tribe. Senior Taliban military commander and member of Taliban Quetta Council as of May 2007. Review pursuant to Security Council resolution 1822 (2008) was concluded on 1 Jun. 2010. INTERPOL-UN Security Council Special Notice web link: https://www.interpol.int/en/How-we-work/Notices/View-UN-Notices-Individuals click here Listed on: 02/04/2001 UK Sanctions List Date Designated: 23/02/2001 Last Updated: 01/02/2021 Group ID: 7060.",
	}

	tk := p.Learn(text)

	if len(tk) <= 10 {
		panic("error")
	}

	p.Save(tk, "model")

	fmt.Println("Tokens", tk)
}
