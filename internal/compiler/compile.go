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

// Temporary
func ConstCompile(OutputFileName string, SourceFileName string) {
	ofile, err := os.Create(OutputFileName)
	if err != nil {
		panic(err)
	}
	defer ofile.Close()

	prefixWriter(ofile)
	_, err = ofile.WriteString("main:\n")
	_, err = ofile.WriteString("  mov rax, 42\n")
	_, err = ofile.WriteString("  ret\n")
	if err != nil {
		panic(err)
	}
}

func CompileMain(OutputFileName string, SourceFileName string) (err error) {
	SourceFile, err = os.Create(SourceFileName)
	defer SourceFile.Close()
	if err != nil {
		return fmt.Errorf("SourceFile can't open")
	}
	OutputFile, err = os.Create(OutputFileName)
	defer OutputFile.Close()
	if err != nil {
		return fmt.Errorf("OutputFile can't open")
	}

	err = prefixWriter(OutputFile)
	if err != nil {
		return err
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
