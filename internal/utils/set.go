package utils

type Set[T comparable] struct {
	m map[T]struct{}
}

func NewSet[T comparable](elements ...T) Set[T] {
	set := Set[T]{m: make(map[T]struct{})}
	for _, element := range elements {
		set.m[element] = struct{}{}
	}
	return set
}

func (s Set[T]) Contains(element T) bool {
	_, ok := s.m[element]
	return ok
}

func (s Set[T]) Add(element T) (existed bool) {
	_, existed = s.m[element]
	s.m[element] = struct{}{}
	return
}

func (s Set[T]) Remove(element T) (existed bool) {
	_, existed = s.m[element]
	delete(s.m, element)
	return
}

func (s Set[T]) ToSlice() []T {
	elements := make([]T, 0, len(s.m)+1)
	for element := range s.m {
		elements = append(elements, element)
	}
	return elements
}
