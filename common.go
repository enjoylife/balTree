package gotree

type TravOrder int

const (
	PreOrder = iota
	PostOrder
	InOrder
	LevelOrder
)

func (order TravOrder) String() string {
	var s string
	switch order {
	case PreOrder:
		s = "pre-order traversal"
	case PostOrder:
		s = "post-order traversal"
	case InOrder:
		s = "in-order traversal"
	case LevelOrder:
		s = "level-order traversal"
	default:
		s = "unkown traversal"
	}
	return s
}

type Direction int

const (
	GT Direction = iota
	EQ
	LT
)

func (d Direction) String() string {
	var s string
	switch d {
	case GT:
		s = "greater than"
	case EQ:
		s = "equal to"
	case LT:
		s = "less than"
	}
	return s
}

type CompareFunc func(interface{}, interface{}) Direction
type IterFunc func(interface{}, interface{})

type Node interface {
	//New(key interface{}, value interface{}, extra ...interface{})
	Key() interface{}
	Value() interface{}
	MinChild() *Node
	MaxChild() *Node
	Children() []*Node
}

type Tree interface {
	Search(interface{}) (interface{}, bool)
	Insert(interface{}, interface{}) (interface{}, bool)
	Remove(interface{}) (interface{}, bool)
	Next(*Node) *Node
	Prev(*Node) *Node
	Traverse(IterFunc)
}
