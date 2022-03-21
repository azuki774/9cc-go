package parser

var (
	ND_NUM = 1
	ND_ADD = 11
	ND_SUB = 12
)

type abstSyntaxNode struct {
	nodeKind      int
	leftHandNode  *abstSyntaxNode
	rightHandNode *abstSyntaxNode
}
