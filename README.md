# TextExtractor: A Powerful Text Processing Tool in Go

Welcome to `TextExtractor`, a versatile Go package designed for extracting and manipulating text with ease and precision. Whether you're dealing with structured or unstructured text, `TextExtractor` offers a suite of tools to tokenize, extract, and transform text data efficiently.


## Key Features

- **Token Extraction**: Seamlessly extract tokens from text within curly braces `{}`. Ideal for parsing templates or structured documents.
- **Regex Generation**: Automatically generate regular expressions from tokens for advanced pattern matching.
- **Contextual Extraction**: Retrieve text segments before or after specific tokens, enabling contextual analysis.
- **Value Extraction**: Utilize trained models to extract values with precision, considering the surrounding context.
- **Data Mapping**: Effortlessly map extracted values to struct fields, streamlining data processing workflows.
- **Precision Calculation**: Evaluate the accuracy of extracted data with built-in precision scoring.
- **Model Persistence**: Save and load your token models, making your data processing repeatable and reliable.


## Installation

Install `TextExtractor` with ease using Go's package manager:

```bash
go get github.com/devalexandre/textextractor
```

```go
package main

import (
    "fmt"
    "github.com/devalexandre/textextractor"
)

func main() {
    // Initialize TextExtractor
    extractor := textextractor.NewTextExtractor()

    // Sample text with tokens
    input := "Hello, {Name}! Your appointment is on {Date}."

    // Extract and process tokens
    tokens := extractor.ExtractTokens(input)
    regexPatterns := extractor.GenerateRegex(tokens)

    fmt.Println("Extracted Tokens:", tokens)
    fmt.Println("Regex Patterns:", regexPatterns)
}
```


### Contribuições e Suporte


## Contributing

Contributions to `TextExtractor` are welcome! Whether it's bug reports, feature requests, or code contributions, feel free to open an issue or submit a pull request.

## Support

If you encounter any problems or have questions, please open an issue on GitHub. We're here to help!

