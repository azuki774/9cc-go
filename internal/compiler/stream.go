package compiler

type stringStream struct {
	str   string
	index int
}

func (ss *stringStream) ok() bool {
	return ss.index+1 <= len(ss.str)
}

func (ss *stringStream) nextChar() (b byte, err error) {
	b = ss.str[ss.index]
	ss.index++
	return b, nil
}

func (ss *stringStream) nextPeekChar() (b byte, err error) {
	b = ss.str[ss.index]
	return b, nil
}

func newStream(str string) *stringStream {
	return &stringStream{str: str, index: 0}
}
