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

type PrecisionWeights struct {
	WordLengthWeight     float64
	TokenLengthWeight    float64
	CharacterCountWeight float64
}

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
	Precision int
	Weights   PrecisionWeights // Adicionado para armazenar os pesos de precisão
}

func NewTextExtractor() *TextExtractor {
	return &TextExtractor{
		ModelsDir: "models",
		Precision: 5,
	}
}

// ExtractTokens usando um analisador personalizado

// ExtractTokens usando um analisador personalizado
func (n TextExtractor) ExtractTokens(input string) []string {
	var tokens []string
	var tokenBuffer string
	insideToken := false

	for _, char := range input {
		switch char {
		case '{':
			insideToken = true
			tokenBuffer = ""
		case '}':
			if insideToken {
				tokens = append(tokens, tokenBuffer)
				insideToken = false
			}
		default:
			if insideToken {
				tokenBuffer += string(char)
			}
		}
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
	regexValue := fmt.Sprintf(`(.{%v})%s`, n.Precision, regexp.QuoteMeta(token))
	regex := regexp.MustCompile(regexValue)

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
	regexValue := fmt.Sprintf(`%s(.{%v})`, regexp.QuoteMeta(token), n.Precision)
	regex := regexp.MustCompile(regexValue)

	// Find the first match in the input string.
	match := regex.FindStringSubmatch(input)

	// If there is no match or the token is at the end of the string, return empty.
	if len(match) < 2 {
		return ""
	}

	// Return the 5 characters after the token.
	return match[1]
}

func (n TextExtractor) GetValueBetweenTokens(input string, model TokenTrain, weights PrecisionWeights) (Extracted, bool) {
	// Verifica se os campos WordBefore e WordAfter são válidos
	if model.WordBefore == "" && model.WordAfter == "" {
		return Extracted{}, false
	}

	escapedWordBefore := regexp.QuoteMeta(model.WordBefore)
	escapedWordAfter := regexp.QuoteMeta(model.WordAfter)
	var regexPattern string

	// Construindo o padrão da expressão regular com base no modelo
	if len(model.WordAfter) == 0 {
		regexPattern = escapedWordBefore + `(.+)`
	} else if len(model.WordBefore) == 0 {
		regexPattern = `(.+)` + escapedWordAfter
	} else {
		regexPattern = escapedWordBefore + `(.+?)` + escapedWordAfter
	}

	// Verifica se o padrão da expressão regular é válido
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return Extracted{}, false
	}

	// Encontrando a correspondência
	match := regex.FindStringSubmatch(input)
	if len(match) < 2 || match[1] == "" {
		return Extracted{}, false
	}

	result := strings.TrimSpace(match[1])

	// Calculando a precisão
	precision := calculatePrecision(result, len(model.Name), len(result), len(match[0]), weights)

	return Extracted{
		Token:     model.Name,
		Value:     result,
		Precision: precision,
	}, true
}

// GetValue extracts values using a trained model, and if not found, it tries the next token using recursion.
// GetValue extrai valores usando um modelo treinado, e se não encontrado, tenta o próximo token usando recursão.
func (n TextExtractor) GetValue(input string, model []TokenTrain) (Extracted, bool) {
	if len(model) == 0 {
		return Extracted{}, false
	}

	extracted, have := n.GetValueBetweenTokens(input, model[0], n.Weights)
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

// Save saves tokens to a .gob file in the "models" folder.
func (n TextExtractor) Save(tokens []TokenTrain, filename string) error {
	// Determine the absolute path to the "models" directory at the project's root.
	dir, err := n.GetModelsDir()
	if err != nil {
		return err
	}

	// Ensure that the "models" directory exists; create it if it doesn't.
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Create the full file path within the "models" directory.
	filePath := filepath.Join(dir, fmt.Sprintf("%s.gob", filename))

	// Open the file for writing (or create it if it doesn't exist).
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a Gob encoder.
	encoder := gob.NewEncoder(file)

	// Encode the tokens into Gob and write to the file.
	if err := encoder.Encode(tokens); err != nil {
		return err
	}

	return nil
}

// Load loads tokens from a .gob file in the "models" folder.
func (n TextExtractor) Load(filename string) ([]TokenTrain, error) {
	// Determine the absolute path to the "models" directory at the project's root.
	modelsDir, err := n.GetModelsDir()
	if err != nil {
		return nil, err
	}

	// Create the full file path within the "models" directory.
	filePath := filepath.Join(modelsDir, fmt.Sprintf("%s.gob", filename))

	// Open the file for reading.
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a Gob decoder.
	decoder := gob.NewDecoder(file)

	// Decode the Gob into a slice of tokens.
	tokens := []TokenTrain{}
	if err := decoder.Decode(&tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (n TextExtractor) ParseValueToStruct(input string, output interface{}, pathFile string) error {
	tagsToFields := make(map[string]string)
	t := reflect.TypeOf(output).Elem()
	tokens, errLoad := n.Load(pathFile)

	if errLoad != nil {
		return errLoad
	}

	// Mapeia tags para campos
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("data")
		if tag != "" {
			tagsToFields[tag] = field.Name
		}
	}

	valueMap := make(map[string]Extracted)

	for _, token := range tokens {
		train := TokenTrain{
			Name:       token.Name,
			WordBefore: strings.Trim(token.WordBefore, "{}"),
			WordAfter:  strings.Trim(token.WordAfter, "{}"),
		}

		extracted, have := n.GetValueBetweenTokens(input, train, n.Weights)
		if have {
			fieldName, tagExists := tagsToFields[extracted.Token]
			if tagExists {
				existingValue, found := valueMap[fieldName]
				if !found || extracted.Precision > existingValue.Precision {
					valueMap[fieldName] = extracted
				}
			}
		}
	}

	// Preenche a estrutura de saída usando os valores do mapa
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
func calculatePrecision(value string, tokenLength, characterCount, tokenCount int, weights PrecisionWeights) float64 {
	// Garantir que os pesos não sejam zero
	if weights.WordLengthWeight == 0 || weights.TokenLengthWeight == 0 || weights.CharacterCountWeight == 0 {
		return 0
	}

	wordLengthPrecision := float64(len(value)) * weights.WordLengthWeight
	tokenLengthPrecision := float64(tokenLength) * weights.TokenLengthWeight
	characterCountPrecision := float64(characterCount) * weights.CharacterCountWeight

	totalWeight := weights.WordLengthWeight + weights.TokenLengthWeight + weights.CharacterCountWeight
	totalPrecision := (wordLengthPrecision + tokenLengthPrecision + characterCountPrecision) / totalWeight

	return totalPrecision
}

// GetModelsDir returns the absolute path to the "models" directory at the project's root.
func (n TextExtractor) GetModelsDir() (string, error) {
	return filepath.Abs(n.ModelsDir)
}
