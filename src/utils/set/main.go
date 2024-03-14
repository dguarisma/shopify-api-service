package set

type Set[K comparable] struct {
	data map[K]struct{}
}

func New[K comparable](size ...int) *Set[K] {
	s := &Set[K]{
		data: map[K]struct{}{},
	}

	if len(size) != 0 {
		s.data = make(map[K]struct{}, size[0])
	}

	return s
}

func (s *Set[K]) Add(data K) {
	s.data[data] = struct{}{}
}

func (s *Set[K]) Exist(data K) bool {
	_, ok := s.data[data]
	return ok
}

func (s *Set[K]) Delete(data K) {
	delete(s.data, data)
}
func (s *Set[K]) Get() []K {
	info := make([]K, 0, len(s.data))
	for data := range s.data {
		info = append(info, data)
	}
	return info
}
