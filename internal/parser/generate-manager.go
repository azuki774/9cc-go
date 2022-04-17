package parser

import "fmt"

type CodeManager struct {
	codes     []string
	noMain    bool // true なら ソースファイル全体をmain関数とする
	jumpLabel int
}

func newCodeManager(noMain bool) *CodeManager {
	return &CodeManager{noMain: noMain}
}

func (cm *CodeManager) getCodes() []string {
	return cm.codes
}

func (cm *CodeManager) getNoMain() bool {
	return cm.noMain
}

func (cm *CodeManager) getJumpLabel() int {
	return cm.jumpLabel
}

func (cm *CodeManager) AddJumpLabel(a int) {
	cm.jumpLabel += a
	return
}

func (cm *CodeManager) AddCode(addCodeFormat string, a ...interface{}) {
	// 置換フォーマットは１つのみ対応
	if a == nil {
		cm.codes = append(cm.codes, addCodeFormat+"\n")
		return
	}
	addCode := fmt.Sprintf(addCodeFormat+"\n", a[0])
	cm.codes = append(cm.codes, addCode)
}

func GenAssembleMain(nodes []*abstSyntaxNode, noMain bool) (codes []string, err error) {
	cm := newCodeManager(noMain)
	cm.genInitCode()
	for _, node := range nodes {
		if err := cm.genCode(node); err != nil {
			return nil, err
		}

		// 各式の計算結果をスタックからraxにpop
		if cm.getNoMain() {
			cm.AddCode("pop rax")
		}
	}
	cm.genEndCode()
	return cm.getCodes(), nil
}
