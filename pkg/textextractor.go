package textextractor

import (
	"encoding/gob"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type TokenTrain struct {
	Name       string
	WordBefore string
	WordAfter  string
	Last       bool
	Order      int
}

type Extracted struct {
	Token     string
	Value     string
	Precision float64
}

type TextExtractor struct{}

func NewTextExtractor() *TextExtractor {
	return &TextExtractor{}
}

// ExtractTokens extracts tokens from input using {}
func (n TextExtractor) ExtractTokens(input string) []string {
	tokens := []string{}

	// Use regex to extract tokens from input in {}
	regex := regexp.MustCompile(`\{([^\}]+)\}`)
	for _, match := range regex.FindAllStringSubmatch(input, -1) {
		tokens = append(tokens, match[1])
	}

	return tokens
}

// GenerateRegex generates regex patterns for tokens
func (n TextExtractor) GenerateRegex(tokens []string) []string {
	regex := []string{}
	for _, token := range tokens {
		regex = append(regex, `(?P<`+token+`>[^\s]+)`)
	}

	return regex
}

// GetBeforeToken retorna os 5 caracteres antes do token na string de entrada.
func (n TextExtractor) GetBeforeToken(input string, token string) string {
	// Define a expressão regular para encontrar o token e os 5 caracteres antes dele.
	regex := regexp.MustCompile(`(.{0,5})` + regexp.QuoteMeta(token))

	// Encontra a primeira correspondência na string de entrada.
	match := regex.FindStringSubmatch(input)

	// Se não houver correspondência, retorna vazio.
	if len(match) < 2 {
		return ""
	}

	// Pega os 5 caracteres antes do token.
	beforeToken := match[1] //  na string ... para tonken COUNTRY, r o beforeToken é "me}. "

	// Se não houver chaves {} nos 5 caracteres antes do token, retorne esses caracteres.
	return beforeToken
}

// GetAfterToken returns the 5 characters after the token in the input string.
func (n TextExtractor) GetAfterToken(input string, token string) string {
	// Define the regular expression to find the token and the 5 characters after it.
	regex := regexp.MustCompile(regexp.QuoteMeta(token) + `(.{5})`)

	// Find the first match in the input string.
	match := regex.FindStringSubmatch(input)

	// If there is no match or the token is at the end of the string, return empty.
	if len(match) < 2 {
		return ""
	}

	// Return the 5 characters after the token.
	return match[1]
}

func (n TextExtractor) GetTokenNameBeforeToken(input string, partialTokenName string) string {
	// Find the partialTokenName in the input string.
	index := strings.Index(input, partialTokenName) + (len(partialTokenName))

	// If the partialTokenName is not found or does not contain '}', return empty.
	if index == -1 || !strings.Contains(partialTokenName, "}") {
		return ""
	}

	// Find the '{' immediately before the partialTokenName.
	openingBraceIndex := strings.LastIndex(input[:index], "{")
	if openingBraceIndex == -1 {
		return ""
	}

	value := input[openingBraceIndex:index] // "{Name}. "

	// Extract everything between the curly braces {} and add it to the results.
	re := regexp.MustCompile(`\{([^{}]+)\}`)
	matches := re.FindStringSubmatch(value)
	if len(matches) >= 2 {
		return matches[0]
	}

	return ""
}

func (n TextExtractor) GetTokenNameAfterToken(input string, partialTokenName string) string {
	// Find the partialTokenName in the input string.
	index := strings.Index(input, partialTokenName) + len(partialTokenName)

	// If the partialTokenName is not found or does not contain '{', return empty.
	if index == -1 || !strings.Contains(partialTokenName, "{") {
		return ""
	}

	// Find the '}' immediately after the partialTokenName.
	closingBraceIndex := strings.Index(input[index:], "}")
	if closingBraceIndex == -1 {
		return ""
	}

	// Adjust the closing brace index to the correct position in the input string.
	closingBraceIndex += index

	// Find the '{' immediately before the partialTokenName.
	openingBraceIndex := strings.LastIndex(input[:index], "{")
	if openingBraceIndex == -1 {
		return ""
	}

	value := input[openingBraceIndex : closingBraceIndex+1] // "{Name}"

	// Extract everything between the curly braces {} and add it to the results.
	re := regexp.MustCompile(`\{([^{}]+)\}`)
	matches := re.FindStringSubmatch(value)
	if len(matches) >= 2 {
		return matches[0]
	}

	return ""
}

func (n TextExtractor) Normalize(input string, tokens []TokenTrain) []TokenTrain {

	// Normalize the input string by replacing tokens with their names
	var normalizedTrain []TokenTrain
	for i, token := range tokens {
		normalized := TokenTrain{}
		normalized.Name = token.Name

		if strings.Contains(token.WordBefore, "}") {
			normalized.WordBefore = n.GetTokenNameBeforeToken(input, token.WordBefore)
		} else {
			normalized.WordBefore = token.WordBefore
		}

		if strings.Contains(token.WordAfter, "{") {
			normalized.WordAfter = n.GetTokenNameAfterToken(input, token.WordAfter)
		} else {
			normalized.WordAfter = token.WordAfter
		}

		if strings.Contains(token.WordBefore, "}") || strings.Contains(token.WordAfter, "{") {
			normalized.Last = true
			normalized.Order = i
		}

		if normalized.WordAfter == "" && strings.Contains(token.WordBefore, "{") {
			normalized.Last = true
			normalized.Order = 0
			normalizedTrain[0].Order = i
		}

		normalizedTrain = append(normalizedTrain, normalized)
	}

	return normalizedTrain

}

// GetValueBetweenTokens extracts the value between tokens using a regex pattern.
func (n TextExtractor) GetValueBetweenTokens(input string, model TokenTrain) (Extracted, bool) {
	var regex *regexp.Regexp
	if len(model.WordAfter) == 0 {
		regex = regexp.MustCompile(model.WordBefore + `(.+)`)
	}

	if len(model.WordBefore) == 0 {
		regex = regexp.MustCompile(`(.+)` + model.WordAfter)
	}

	if len(model.WordAfter) > 0 && len(model.WordBefore) > 0 {
		regex = regexp.MustCompile(model.WordBefore + `(.+?)` + model.WordAfter)
	}

	match := regex.FindStringSubmatch(input)
	var result string
	var precision float64
	if len(match) >= 2 {
		result = match[1]
		// Remove leading and trailing spaces
		result = strings.TrimSpace(result)

		// Calculate precision based on the ratio of characters in the word to the context
		contextLength := len(match[0]) // Length of context between the tokens
		wordLength := len(result)      // Length of the extracted word
		if contextLength > 0 {
			precision = calculatePrecision(result, len(model.Name), wordLength, contextLength)
		}

		return Extracted{
			Token:     model.Name,
			Value:     result,
			Precision: precision,
		}, true
	}
	return Extracted{}, false
}

// GetValue extracts values using a trained model, and if not found, it tries the next token using recursion.
func (n TextExtractor) GetValue(input string, model []TokenTrain) (Extracted, bool) {
	if len(model) == 0 {
		return Extracted{}, false
	}
	extracted, have := n.GetValueBetweenTokens(input, model[0])
	if have {
		return extracted, true
	}
	return n.GetValue(input, model[1:])
}

// Learn generates token training data from input strings.
func (n TextExtractor) Learn(input []string) []TokenTrain {
	tokens := []TokenTrain{}

	for _, i := range input {
		t := TokenTrain{}
		// Can have more than one token in the same string
		for _, token := range n.ExtractTokens(i) {
			t.Name = token
			t.WordBefore = n.GetBeforeToken(i, fmt.Sprintf("{%s}", token))
			t.WordAfter = n.GetAfterToken(i, fmt.Sprintf("{%s}", token))
			tokens = append(tokens, t)
		}
	}

	return tokens
}

// Save saves tokens to a .gob file.
func (n TextExtractor) Save(tokens []TokenTrain, filename string) error {
	// Open the file for writing (or create if it doesn't exist)
	file, err := os.Create(fmt.Sprintf("%s.gob", filename))
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a Gob encoder
	encoder := gob.NewEncoder(file)

	// Encode tokens into Gob and write to the file
	if err := encoder.Encode(tokens); err != nil {
		return err
	}

	return nil
}

// Load loads tokens from a .gob file.
func (n TextExtractor) Load(filename string) ([]TokenTrain, error) {
	// Open the file for reading
	file, err := os.Open(fmt.Sprintf("%s.gob", filename))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a Gob decoder
	decoder := gob.NewDecoder(file)

	// Decode the Gob into a slice of tokens
	tokens := []TokenTrain{}
	if err := decoder.Decode(&tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

// ParseValueToStruct parses values from input and populates a struct based on data tags.
func (n TextExtractor) ParseValueToStruct(input string, output interface{}, pathFile string) bool {
	tagsToFields := make(map[string]string)
	t := reflect.TypeOf(output).Elem()
	tokens, _ := n.Load(pathFile)

	// Map tags to fields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("data")
		if tag != "" {
			tagsToFields[tag] = field.Name
		}
	}

	valueMap := make(map[string]Extracted) // Change to store Extracted instead of string

	for _, token := range tokens {
		train := TokenTrain{
			Name:       token.Name,
			WordBefore: token.WordBefore,
			WordAfter:  token.WordAfter,
		}
		p, have := n.GetValueBetweenTokens(input, train)
		if have {
			// Check if the tag is mapped to a field in the structure
			fieldName, tagExists := tagsToFields[p.Token]
			if tagExists {
				// Check if there is already a value for the same key
				existingValue, found := valueMap[fieldName]
				if !found || p.Precision < existingValue.Precision {
					// If there is no existing value or the new precision is higher, update the map
					valueMap[fieldName] = p
				}
			}
		}
	}

	// Populate the output structure using the values from the map
	outputValue := reflect.ValueOf(output).Elem()
	for fieldName, extracted := range valueMap {
		field := outputValue.FieldByName(fieldName)
		if field.IsValid() && field.CanSet() {
			field.SetString(extracted.Value)
		}
	}

	return true
}

// calculatePrecision calculates precision of the extracted value.
func calculatePrecision(value string, tokenLength, characterCount, tokenCount int) float64 {
	// Adjust these weights according to your preference
	wordLengthWeight := 0.4
	tokenLengthWeight := 0.3
	characterCountWeight := 0.3

	// Calculate precision based on provided weights and values
	wordLengthPrecision := float64(len(value)) / float64(wordLengthWeight)
	tokenLengthPrecision := float64(tokenLength) / float64(tokenLengthWeight)
	characterCountPrecision := float64(characterCount) / float64(characterCountWeight)

	// Combine weighted precision values
	totalPrecision := (wordLengthPrecision + tokenLengthPrecision + characterCountPrecision) / 3.0

	// Convert precision to a scale of 0 to 100
	scaledPrecision := totalPrecision * 100.0

	return scaledPrecision
}
