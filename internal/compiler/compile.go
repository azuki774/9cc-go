package compiler

import (
	"os"
)

// Temporary
func ConstCompile(OutputFileName string, SourceFileName string) {
	ofile, err := os.Create(OutputFileName)
	if err != nil {
		panic(err)
	}
	defer ofile.Close()

	_, err = ofile.WriteString(".intel_syntax noprefix\n")
	_, err = ofile.WriteString(".globl main\n")
	_, err = ofile.WriteString("main:\n")
	_, err = ofile.WriteString("  mov rax, 42\n")
	_, err = ofile.WriteString("  ret\n")
	if err != nil {
		panic(err)
	}
}
