package inseki

import (
	"errors"
	"math"
)

type Set[T any] struct {
	// We use a Set with a slice to store the elements
	// More memory efficient than a hashmap for huge Sets
	elements []T
}

func NewSet[T any]() *Set[T] {
	return &Set[T]{}
}

func (s *Set[T]) Add(element T, comparator func(T, T) int8) bool {
	if s.elements == nil {
		s.elements = make([]T, 0)
	}

	if index := s.Contains(element, comparator); index == math.MaxUint16 {
		s.insert(element, comparator)
		return true
	}

	return false

}

func (s *Set[T]) insert(element T, comparator func(T, T) int8) {
	// The set is sorted, insert the element at the right position
	// Use dichotomy search
	left := 0
	right := len(s.elements) - 1
	for left <= right {
		mid := (left + right) / 2
		switch comparator(s.elements[mid], element) {
		case 1:
			right = mid - 1
		case -1:
			left = mid + 1
		}
	}

	// Insert the element
	s.elements = append(s.elements, element)
	copy(s.elements[left+1:], s.elements[left:])
	s.elements[left] = element
}

func (s *Set[T]) Get() (T, error) {
	if s.IsEmpty() {
		var result T
		return result, errors.New("set is empty")
	}

	element := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]

	return element, nil
}

func (s *Set[T]) IsEmpty() bool {
	return len(s.elements) == 0
}

func (s *Set[T]) Size() int {
	return len(s.elements)
}

func (s *Set[T]) Clear() {
	s.elements = make([]T, 0)
}

/*
Contains checks if the element is in the Set

The comparator function must return:
  - 0 if the elements are equal
  - 1 if the first element is greater than the second
  - -1 if the first element is less than the second

Returns the position of the element in the set or -1 if the element is not in the Set
*/
func (s *Set[T]) Contains(element T, comparator func(T, T) int8) uint16 {

	// Using dichotomy search
	left := 0
	right := len(s.elements) - 1
	for left <= right {
		mid := (left + right) / 2
		switch comparator(s.elements[mid], element) {
		case 0:
			return uint16(mid)
		case 1:
			right = mid - 1
		case -1:
			left = mid + 1
		}
	}

	return math.MaxUint16
}
