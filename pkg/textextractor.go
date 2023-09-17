package textextractor

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

type TokenTrain struct {
	Name       string
	WordBefore string
	WordAfter  string
}

type Extracted struct {
	Token     string
	Value     string
	Precision float64
}

type TextExtractor struct {
	ModelsDir string
}

func NewTextExtractor() *TextExtractor {
	return &TextExtractor{
		ModelsDir: "models",
	}
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

// GetBeforeToken returns the 5 characters before the token in the input string.
func (n TextExtractor) GetBeforeToken(input string, token string) string {
	// Define the regular expression to find the token and the 5 characters before it.
	regex := regexp.MustCompile(`(.{5})` + regexp.QuoteMeta(token))

	// Find the first match in the input string.
	match := regex.FindStringSubmatch(input)

	// If there is no match or the token is at the beginning of the string, return empty.
	if len(match) < 2 {
		return ""
	}

	// Return the 5 characters before the token.
	return match[1]
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

// Save salva os tokens em um arquivo .gob na pasta "models".
func (n TextExtractor) Save(tokens []TokenTrain, filename string) error {
	// Determine o caminho absoluto para o diretório "models" na raiz do projeto.
	dir, err := n.GetModelsDir()
	if err != nil {
		return err
	}

	// Garanta que o diretório "models" exista; crie-o se não existir.
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Crie o caminho completo do arquivo no diretório "models".
	filePath := filepath.Join(dir, fmt.Sprintf("%s.gob", filename))

	// Abra o arquivo para escrita (ou crie se não existir).
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Crie um codificador Gob.
	encoder := gob.NewEncoder(file)

	// Codifique os tokens em Gob e escreva no arquivo.
	if err := encoder.Encode(tokens); err != nil {
		return err
	}

	return nil
}

// Load carrega tokens de um arquivo .gob na pasta "models".
func (n TextExtractor) Load(filename string) ([]TokenTrain, error) {
	// Determine o caminho absoluto para o diretório "models" na raiz do projeto.
	modelsDir, err := n.GetModelsDir()
	if err != nil {
		return nil, err
	}

	// Crie o caminho completo do arquivo no diretório "models".
	filePath := filepath.Join(modelsDir, fmt.Sprintf("%s.gob", filename))

	// Abra o arquivo para leitura.
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Crie um decodificador Gob.
	decoder := gob.NewDecoder(file)

	// Decodifique o Gob em uma fatia de tokens.
	tokens := []TokenTrain{}
	if err := decoder.Decode(&tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

// ParseValueToStruct parses values from input and populates a struct based on data tags.
func (n TextExtractor) ParseValueToStruct(input string, output interface{}, pathFile string) error {
	tagsToFields := make(map[string]string)
	t := reflect.TypeOf(output).Elem()
	tokens, errLoad := n.Load(pathFile)

	if errLoad != nil {
		return errLoad
	}

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

	return nil
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

// GetModelsDir retorna o caminho absoluto para o diretório "models" na raiz do projeto.
func (n TextExtractor) GetModelsDir() (string, error) {
	return filepath.Abs(n.ModelsDir)
}
