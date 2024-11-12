package inseki

type Stack struct {
	Values []Target
}

// Push : Add a value to the stack
func (s *Stack) Push(value Target) {
	s.Values = append(s.Values, value)
}

// Pop : Remove a value from the stack, and return it
func (s *Stack) Pop() Target {
	if len(s.Values) == 0 {
		return Target{}
	}

	value := s.Values[len(s.Values)-1]
	s.Values = s.Values[:len(s.Values)-1]

	return value
}

// Peek : Get the first element of the stack without removing it
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

// Print : Show each value of the stack
func (s *Stack) Print() {
	for _, value := range s.Values {
		println(value.Filepath)
	}
}
