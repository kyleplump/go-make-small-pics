package main

type Color struct {
	r uint16
	g uint16
	b uint16
	a uint16
	run uint16
}

type Stack struct {
	items []Color
}

func (s *Stack) Push(v Color) {
	s.items = append(s.items, v);
}

func (s *Stack) Pop() (Color, bool) {
	popped_value := s.items[len(s.items) - 1]
	s.items = s.items[:len(s.items) - 1]
	return popped_value, true;
}
