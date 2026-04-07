package codegen

// NodeType AST 节点类型 (用于代码生成)
type NodeType int

const (
	NodeElement NodeType = iota
	NodeText
	NodeExpression
	NodeFragment
)

// AttrValue 属性值
type AttrValue struct {
	IsExpression bool
	Value        string
}

// Node AST 节点 (用于代码生成)
type Node struct {
	Type       NodeType
	TagName    string
	Attributes map[string]AttrValue
	Children   []Node
	Value      string // Text, Expression
	Spread     []string
}