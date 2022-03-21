package compiler

import "os"

// 与えられた文字列のスライスを１つずつ書き出す
func stringsWriter(ofile *os.File, strings []string) (err error) {
	for _, s := range strings {
		_, err = ofile.WriteString(s)
		if err != nil {
			return err
		}
	}
	return nil
}

// アセンブリのprefix部分を書き込む
func prefixWriter(ofile *os.File) (err error) {
	slist := []string{".intel_syntax noprefix\n", ".globl main\n"}
	err = stringsWriter(ofile, slist)
	return err
}
