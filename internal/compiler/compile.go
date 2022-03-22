package compiler

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/azuki774/9cc-go/internal/parser"
)

var (
	SourceFile *os.File
	OutputFile *os.File
)

func CompileMain(OutputFileName string, SourceFileName string) (err error) {
	b, err := ioutil.ReadFile(SourceFileName)
	defer SourceFile.Close()
	if err != nil {
		return fmt.Errorf("CompileMain: SourceFile can't open")
	}
	OutputFile, err = os.Create(OutputFileName)
	defer OutputFile.Close()
	if err != nil {
		return fmt.Errorf("CompileMain: OutputFile can't open")
	}

	tokens, err := parser.TokenizeMain(string(b))
	if err != nil {
		return fmt.Errorf("CompileMain : Tokenize error : %w", err)
	}

	topNode := parser.ParserMain(tokens)
	codes, _ := parser.GenAssembleMain(topNode)

	err = stringsWriter(OutputFile, codes)
	if err != nil {
		return fmt.Errorf("CompileMain : stringsWriter error : %w", err)
	}

	return nil
}

// Tokenizeのみ実行する、アウトプットは標準出力
func TokenizeOnly(SourceFileName string) (err error) {
	b, err := ioutil.ReadFile(SourceFileName)
	defer SourceFile.Close()
	if err != nil {
		return fmt.Errorf("TokenizeOnly: SourceFile can't open")
	}

	fmt.Println("Tokenize Only")
	tokens, err := parser.TokenizeMain(string(b))
	if err != nil {
		return fmt.Errorf("TokenizeOnly : %w", err)
	}

	for _, token := range tokens {
		token.Show()
	}

	return nil
}
