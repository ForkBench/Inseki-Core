package inseki

type Stack[T any] struct {
	// Stack : We use a stack with a slice to store the elements
	elements []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Push(element T) {
	s.elements = append(s.elements, element)
}