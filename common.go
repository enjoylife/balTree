package gotree

type CompareFunc func(interface{}, interface{}) int
type IterFunc func(interface{}, interface{})

type Node interface {
	Next() *Node
	Prev() *Node
	Key() interface{}
	Value() interface{}
}

type Tree interface {
	Search(interface{}) (interface{}, bool)
	Insert(interface{}) interface{}
	Delete(interface{}) interface{}
	Traverse() interface{}
	Compare(interface{}, interface{}) int
}
