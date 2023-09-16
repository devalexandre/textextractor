# NLPK: Natural Language Processing Toolkit

![NLPK Logo](https://example.com/nlpk-logo.png)

NLPK is a powerful toolkit for natural language processing tasks. It provides various functions to extract, manipulate, and analyze text data efficiently.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
  - [Extract Tokens](#extract-tokens)
  - [Get Word Before Token](#get-word-before-token)
  - [Get Word After Token](#get-word-after-token)
  - [Get Value Between Tokens](#get-value-between-tokens)
  - [Parse Value to Struct](#parse-value-to-struct)
  - [Learn from Training Data](#learn-from-training-data)
  - [Save and Load Tokens](#save-and-load-tokens)

## Installation

To use NLPK, you need to import it into your Go project:

```shell
go get github.com/devalexandre/nlpk

```
# Usage

## Extract Tokens

📜 Function: ExtractTokens(input string) []string

Extract tokens from a text.

```go
nlpk := nlpk.NewNLPK()
input := "Name: {Name}. DOB: {DOB}. Country: {Country}"
tokens := nlpk.ExtractTokens(input)

// Result: tokens = ["Name", "DOB", "Country"]
```

## Get Word Before Token

📜 Function: GetBeforeToken(input, token string) string

Get the word before a specific token in a text.

```go
nlpk := nlpk.NewNLPK()
input := "Name: John. Country: USA"
token := "Name"
wordBefore := nlpk.GetBeforeToken(input, fmt.Sprintf("{%s}", token))

// Result: wordBefore = ""
```
## Get Word After Token

Similarly, you can get the word following a specific token in the text.

📜 Function: GetAfterToken(input, token string) string

```go
nlpk := nlpk.NewNLPK()
input := "Name: John. Country: USA"
token := "Country"
wordAfter := nlpk.GetAfterToken(input, fmt.Sprintf("{%s}", token))

// Result: wordAfter = ": USA"
```

## Get Value Between Tokens

📜 Function: GetValueBetweenTokens(input string, train TokenTrain) (ParsedValue, bool)

Get the value between two tokens using a train configuration.

```go
nlpk := nlpk.NewNLPK()
trainWord := "Name: {NAME}. DOB: {DOB}"
tokens := nlpk.ExtractTokens(trainWord)
input := "Name: John. DOB: 01/01/1990"

values := make(map[string]string)

for _, token := range tokens {
	wordBefore := nlpk.GetBeforeToken(trainWord, fmt.Sprintf("{%s}", token))
	wordAfter := nlpk.GetAfterToken(trainWord, fmt.Sprintf("{%s}", token))
	train := nlpk.TokenTrain{
		Name:       token,
		WordBefore: wordBefore,
		WordAfter:  wordAfter,
	}
	value, found := nlpk.GetValueBetweenTokens(input, train)

	if found {
		values[token] = value
	}
}

// Result: values = {"NAME": "John", "DOB": "01/01/1990"}

```

## Parse Value to Struct

📜 Function: ParseValueToStruct(input string, output interface{}) bool

Parse values from a text into a struct using struct tags.

```go
nlpk := nlpk.NewNLPK()
input := "Name: John. DOB: 01/01/1990"

type Person struct {
	Name string `data:"NAME"`
	DOB  string `data:"DOB"`
}
person := Person{}
ok := nlpk.ParseValueToStruct(input, &person, "tokens.json")

// Result: person = {Name: "John", DOB: "01/01/1990"}

```

## Learn from Training Data

📜 Function: Learn(dataTrain []string) map[string]TokenTrain

Learn token patterns from a set of training data.

```go
nlpk := nlpk.NewNLPK()

dataTrain := []string{
    "Name 6: {NAME}. DOB:{DOB}",
    "Title: {TITLE} DOB: {DOB}",
}
learnedTokens := nlpk.Learn(dataTrain)
fmt.Println("Learned Tokens:", learnedTokens)
```

## Save and Load Tokens

📜 Functions: Save(tokens map[string]TokenTrain, filename string) error, Load(filename string) (map[string]TokenTrain, error)

Save learned tokens to a file and load them for future use.

```go
nlpk := nlpk.NewNLPK()

learnedTokens := map[string]TokenTrain{
    "NAME": {Name: "NAME", WordBefore: " 6: ", WordAfter: ". DOB:"},
    "DOB":  {Name: "DOB", WordBefore: "Title: ", WordAfter: ""},
}
err := nlpk.Save(learnedTokens, "tokens.json")
if err != nil {
    fmt.Println("Error saving tokens:", err)
}

loadedTokens, err := nlpk.Load("tokens.json")
if err != nil {
    fmt.Println("Error loading tokens:", err)
} else {
    fmt.Println("Loaded Tokens:", loadedTokens)
}
```