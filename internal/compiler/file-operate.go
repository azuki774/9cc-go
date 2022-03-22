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
