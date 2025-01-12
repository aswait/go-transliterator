package lexicalanalyzer

import (
	"fmt"
	"strings"

	sourcetext "gitgub.com/aswait/go-transliterator/pkg/source-text"
)

type LexicalAnalyzerer interface {
	Transliterate() (string, error)
	Validate() string
	SourceLoadFromInput(input string)
}

type LexicalAnalyzer struct {
	Source   sourcetext.SourceTexter
	alphabet map[rune]string
}

func NewLexicalAnalyzer(source sourcetext.SourceTexter) *LexicalAnalyzer {
	alphabet := make(map[rune]string)

	for r := 'a'; r <= 'z'; r++ {
		alphabet[r] = "Letter"
	}
	for r := 'A'; r <= 'Z'; r++ {
		alphabet[r] = "Letter"
	}

	for r := '0'; r <= '9'; r++ {
		alphabet[r] = "Digit"
	}

	alphabet[' '] = "Space"
	alphabet['\n'] = "EndRow"
	alphabet['\t'] = "Tab"
	alphabet['/'] = "Comment"
	alphabet['*'] = "Comment"

	return &LexicalAnalyzer{
		Source:   source,
		alphabet: alphabet,
	}
}

func (la *LexicalAnalyzer) Transliterate() (string, error) {
	var result strings.Builder

	for la.Source.HasMoreSymbols() {
		symbol, err := la.Source.ReadNextSymbol()
		if err != nil {
			return "", fmt.Errorf("Error reading text: %v", err)
		}

		class, exists := la.alphabet[symbol]
		if !exists {
			return "", fmt.Errorf("Символ '%c' не принадлежит алфавиту", symbol)
		}
		if symbol == '\n' || symbol == '\t' {
			result.WriteString(fmt.Sprintf("(%s)\n", class))
		} else {
			result.WriteString(fmt.Sprintf("(%s, %c)\n", class, symbol))
		}

	}

	return result.String(), nil
}

func (la *LexicalAnalyzer) Validate() string {
	if la.Source.HasMoreSymbols() {
		return "Error: Text not fully processed."
	}
	return "Текст верен"
}

func (la *LexicalAnalyzer) SourceLoadFromInput(input string) {
	la.Source.LoadFromInput(input)
}
