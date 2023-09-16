package prose

import (
	"encoding/json"
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
}

type Extracted struct {
	Token string
	Value string
}

type PROSE struct{}

func NewPROSE() *PROSE {
	return &PROSE{}
}

// get tokens from input using {}
func (n PROSE) ExtractTokens(input string) []string {
	tokens := []string{}

	//use regex to extract tokens from input in {}
	regex := regexp.MustCompile(`\{([^\}]+)\}`)
	for _, match := range regex.FindAllStringSubmatch(input, -1) {
		tokens = append(tokens, match[1])
	}

	return tokens
}

// get word before token and after token and create regex to get value between them
func (n PROSE) GenerateRegex(tokens []string) []string {
	regex := []string{}
	for _, token := range tokens {
		regex = append(regex, `(?P<`+token+`>[^\s]+)`)
	}

	return regex
}

// GetBeforeToken retorna os 5 caracteres antes do token na string de entrada.
func (n PROSE) GetBeforeToken(input string, token string) string {
	// Define a expressão regular para encontrar o token e os 5 caracteres antes.
	regex := regexp.MustCompile(`(.{5})` + regexp.QuoteMeta(token))

	// Encontra a primeira correspondência na string de entrada.
	match := regex.FindStringSubmatch(input)

	// Se não houver correspondência ou o token estiver no início da string, retorna vazio.
	if len(match) < 2 {
		return ""
	}

	// Retorna os 5 caracteres antes do token.
	return match[1]
}

// GetAfterToken retorna os 5 caracteres após o token na string de entrada.
func (n PROSE) GetAfterToken(input string, token string) string {
	// Define a expressão regular para encontrar o token e os 5 caracteres após.
	regex := regexp.MustCompile(regexp.QuoteMeta(token) + `(.{5})`)

	// Encontra a primeira correspondência na string de entrada.
	match := regex.FindStringSubmatch(input)

	// Se não houver correspondência ou o token estiver no final da string, retorna vazio.
	if len(match) < 2 {
		return ""
	}

	// Retorna os 5 caracteres após o token.
	return match[1]
}

// get value between word before and word after
func (n PROSE) GetValueBetweenTokens(input string, model TokenTrain) (Extracted, bool) {
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
	if len(match) >= 2 {
		result = match[1]
		//remove espaço no inicio e no fim
		result = strings.TrimSpace(result)
		return Extracted{
			Token: model.Name,
			Value: result,
		}, true
	}
	return Extracted{}, false
}

// get value using trained model, and if not found value tray next token using recursion
func (n PROSE) GetValue(input string, model []TokenTrain) (Extracted, bool) {
	if len(model) == 0 {
		return Extracted{}, false
	}
	extracted, have := n.GetValueBetweenTokens(input, model[0])
	if have {
		return extracted, true
	}
	return n.GetValue(input, model[1:])
}

func (n PROSE) Learn(input []string) []TokenTrain {
	tokens := []TokenTrain{}

	for _, i := range input {
		t := TokenTrain{}
		//can have more than one token in the same string
		for _, token := range n.ExtractTokens(i) {
			t.Name = token
			t.WordBefore = n.GetBeforeToken(i, fmt.Sprintf("{%s}", token))
			t.WordAfter = n.GetAfterToken(i, fmt.Sprintf("{%s}", token))
			tokens = append(tokens, t)
		}

	}

	return tokens
}

func (n PROSE) Save(tokens []TokenTrain, filename string) error {
	// Abre o arquivo para escrita (ou cria se não existir)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Cria um encoder JSON
	encoder := json.NewEncoder(file)

	// Codifica os tokens em JSON e escreve no arquivo
	if err := encoder.Encode(tokens); err != nil {
		return err
	}

	return nil
}

// Load carrega os tokens de um arquivo JSON
func (n PROSE) Load(filename string) ([]TokenTrain, error) {
	// Abre o arquivo para leitura
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Cria um decoder JSON
	decoder := json.NewDecoder(file)

	// Decodifica o JSON para um slice de tokens
	tokens := []TokenTrain{}
	if err := decoder.Decode(&tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (n PROSE) ParseValueToStruct(input string, output interface{}, pathFile string) bool {
	tagsToFields := make(map[string]string)
	t := reflect.TypeOf(output).Elem()
	tokens, _ := n.Load(pathFile)

	// Mapear tags para campos
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("data")
		if tag != "" {
			tagsToFields[tag] = field.Name
		}
	}

	valueMap := make(map[string]string)

	for _, token := range tokens {
		train := TokenTrain{
			Name:       token.Name,
			WordBefore: token.WordBefore,
			WordAfter:  token.WordAfter,
		}
		p, have := n.GetValueBetweenTokens(input, train)
		if have {
			// Verifique se a tag está mapeada para um campo na estrutura
			fieldName, tagExists := tagsToFields[p.Token]
			if tagExists {
				// Preencha o mapa temporário com os valores extraídos
				valueMap[fieldName] = p.Value
			}
		}
	}

	// Preencha a estrutura de saída usando os valores do mapa
	outputValue := reflect.ValueOf(output).Elem()
	for fieldName, value := range valueMap {
		field := outputValue.FieldByName(fieldName)
		if field.IsValid() && field.CanSet() {
			field.SetString(value)
		}
	}

	return true
}
