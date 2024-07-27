package tools

type Stack struct {
	Values []Target
}

func (s *Stack) Push(value Target) {
	s.Values = append(s.Values, value)
}

func (s *Stack) Pop() Target {
	if len(s.Values) == 0 {
		return Target{}
	}

	value := s.Values[len(s.Values)-1]
	s.Values = s.Values[:len(s.Values)-1]

	return value
}

func (s *Stack) Peek() Target {
	if len(s.Values) == 0 {
		return Target{}
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
	s.Values = []Target{}
}

func (s *Stack) Print() {
	for _, value := range s.Values {
		println(value.Filepath)
	}
}
