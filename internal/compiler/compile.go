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

func CompileMain(OutputFileName string, SourceFileName string, noMain bool, showTokenize bool) (err error) {
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

	if showTokenize {
		for _, token := range tokens {
			fmt.Printf("%s\n", token.ShowString())
		}
	}
	topNode, err := parser.ParserMain(tokens, noMain)
	if err != nil {
		return err
	}

	codes, _ := parser.GenAssembleMain(topNode, noMain)

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

	for i, token := range tokens {
		fmt.Printf("%d: %s\n", i, token.ShowString())
	}

	return nil
}
