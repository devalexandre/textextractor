# TextExtractor: Text Data Extraction Toolkit

![Prose Logo](https://example.com/Prose-logo.png)

TextExtractor is a versatile toolkit designed for efficient text data extraction tasks. It simplifies the process of extracting valuable information enclosed between tokens in a given input text.

### Simplified Text Extraction

TextExtractor streamlines text extraction using a straightforward approach. Consider this practical example

**Example Input Text**: "Product: {PRODUCT_NAME}. Price: {PRICE} USD"

**Extraction Process:**

```
Input Text: My name is {<------->} and I'm a developer.
```

1. The library searches for the first token marker '{' and '}' pair in the input text.
2. Once the opening '{' is found, the algorithm starts reading characters to find the value inside the token.
3. The library identifies the closing '}' of the first token.
4. The extracted value 'NAME' is stored for later use.
5. The algorithm continues searching for the next token.
6. If additional tokens are present, the process repeats to extract their values.

**Result:**

- The value 'PRODUCT_NAME' is extracted from the token pair '{PRODUCT_NAME}'.
- The value 'PRICE' is extracted from the token pair '{PRICE}'.



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

To use TextExtractor in your Go project, simply import it using:

```shell
go get github.com/devalexandre/TextExtractor

```
# Usage

## Extract Tokens

📜 Function: ExtractTokens(input string) []string

Extract tokens from a text.

```go
p := textextractor.NewTextExtractor()
input := "Name: {Name}. DOB: {DOB}. Country: {Country}"
tokens := p.ExtractTokens(input)

// Result: tokens = ["Name", "DOB", "Country"]
```

## Get Word Before Token

📜 Function: GetBeforeToken(input, token string) string

Get the word before a specific token in a text.

```go
p := textextractor.NewTextExtractor()
input := "Name: John. Country: USA"
token := "Name"
wordBefore := p.GetBeforeToken(input, fmt.Sprintf("{%s}", token))

// Result: wordBefore = ""
```
## Get Word After Token

Similarly, you can get the word following a specific token in the text.

📜 Function: GetAfterToken(input, token string) string

```go
p := textextractor.NewTextExtractor()
input := "Name: John. Country: USA"
token := "Country"
wordAfter := p.GetAfterToken(input, fmt.Sprintf("{%s}", token))

// Result: wordAfter = ": USA"
```

## Get Value Between Tokens

📜 Function: GetValueBetweenTokens(input string, train TokenTrain) (ParsedValue, bool)

Get the value between two tokens using a train configuration.

```go
p := textextractor.NewTextExtractor()
trainWord := "Name: {NAME}. DOB: {DOB}"
tokens := p.ExtractTokens(trainWord)
input := "Name: John. DOB: 01/01/1990"

values := make(map[string]string)

for _, token := range tokens {
	wordBefore := p.GetBeforeToken(trainWord, fmt.Sprintf("{%s}", token))
	wordAfter := p.GetAfterToken(trainWord, fmt.Sprintf("{%s}", token))
	train := p.TokenTrain{
		Name:       token,
		WordBefore: wordBefore,
		WordAfter:  wordAfter,
	}
	value, found := p.GetValueBetweenTokens(input, train)

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
p := textextractor.NewTextExtractor()
input := "Name: John. DOB: 01/01/1990"

type Person struct {
	Name string `data:"NAME"`
	DOB  string `data:"DOB"`
}
person := Person{}
ok := p.ParseValueToStruct(input, &person, "tokens.json")

// Result: person = {Name: "John", DOB: "01/01/1990"}

```

## Learn from Training Data

📜 Function: Learn(dataTrain []string) map[string]TokenTrain

Learn token patterns from a set of training data.

```go
p := textextractor.NewTextExtractor()

dataTrain := []string{
    "Name 6: {NAME}. DOB:{DOB}",
    "Title: {TITLE} DOB: {DOB}",
}
learnedTokens := p.Learn(dataTrain)
fmt.Println("Learned Tokens:", learnedTokens)
```

## Save and Load Tokens

📜 Functions: Save(tokens map[string]TokenTrain, filename string) error, Load(filename string) (map[string]TokenTrain, error)

Save learned tokens to a file and load them for future use.

```go
p := textextractor.NewTextExtractor()

learnedTokens := map[string]TokenTrain{
    "NAME": {Name: "NAME", WordBefore: " 6: ", WordAfter: ". DOB:"},
    "DOB":  {Name: "DOB", WordBefore: "Title: ", WordAfter: ""},
}
err := p.Save(learnedTokens, "tokens.json")
if err != nil {
    fmt.Println("Error saving tokens:", err)
}

loadedTokens, err := p.Load("tokens.json")
if err != nil {
    fmt.Println("Error loading tokens:", err)
} else {
    fmt.Println("Loaded Tokens:", loadedTokens)
}
```