package compiler

import (
	"fmt"
	"os"
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

func TokenizeOnly(OutputFileName string, SourceFileName string) (err error) {
	SourceFile, err = os.Open(SourceFileName)
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

	fmt.Println("Tokenize Only")

	return nil
}
