package tools

type StackValue struct {
	Filepath    string
	Association Association
}

type Stack struct {
	Values []StackValue
}

func (s *Stack) Push(value StackValue) {
	s.Values = append(s.Values, value)
}

func (s *Stack) Pop() StackValue {
	if len(s.Values) == 0 {
		return StackValue{}
	}

	value := s.Values[len(s.Values)-1]
	s.Values = s.Values[:len(s.Values)-1]

	return value
}

func (s *Stack) Peek() StackValue {
	if len(s.Values) == 0 {
		return StackValue{}
	}

	return s.Values[len(s.Values)-1]
}

func (s *Stack) IsEmpty() bool {
	return len(s.Values) == 0
}

func (s *Stack) Len() int {
	return len(s.Values)
}

func (s *Stack) Clear() {
	s.Values = []StackValue{}
}
