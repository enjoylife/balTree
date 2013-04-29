package gotree

type CompareFunc func(interface{}, interface{}) int
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
