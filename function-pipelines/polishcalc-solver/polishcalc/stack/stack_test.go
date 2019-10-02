package stack

import (
	"testing"
)

func assertEqual(expected interface{}, got interface{}, t *testing.T) {
	if expected != got {
		t.Errorf("Expected `%v` but got `%v`", expected, got)
	}
}

func TestPopEmptyStackReturnsNil(t *testing.T) {
	var expected interface{} = nil

	s := New()
	got := s.Pop()

	assertEqual(expected, got, t)
}

func TestPushPop(t *testing.T) {
	expected := 10
	s := New()

	s.Push(expected)
	got := s.Pop()

	assertEqual(expected, got, t)
}

func TestPeekEmptyStackReturnsNil(t *testing.T) {
	var expected interface{} = nil

	s := New()
	got := s.Peek()
	assertEqual(expected, got, t)
}

func TestPeekReturnsLastPushElement(t *testing.T) {
	expected := 10
	s := New()

	s.Push(expected)
	got := s.Peek()

	assertEqual(expected, got, t)
}

func TestSizeShouldUpdate(t *testing.T) {
	// Size starts at 0
	var expected uint = 0
	s := New()

	got := s.Size()
	assertEqual(expected, got, t)

	// When elements are pushed it increases
	expected = 2
	s.Push(1)
	s.Push(2)
	got = s.Size()
	assertEqual(expected, got, t)

	// When elements are popped it decreases, if there are no more elements
	// it stays at 0
	expected = 0
	s.Pop()
	s.Pop()
	got = s.Size()
	assertEqual(expected, got, t)

	// Should not decrease on empty stack
	s.Pop()
	s.Pop()
	got = s.Size()
	assertEqual(expected, got, t)
}

func TestStackShouldAutoExpand(t *testing.T) {
	var expected uint = 512
	s := New()

	assertEqual(expected, s.capacity, t)

	itemQuantity := int(2*chunkSize + 1)
	for i := 0; i < itemQuantity; i++ {
		s.Push(i)
	}

	assertEqual(itemQuantity, int(s.Size()), t)
	expected = 3*chunkSize
	assertEqual(expected, s.capacity, t)

	for i := int(s.Size() - 1); i >= 0; i-- {
		got := s.Pop()
		assertEqual(i, got, t)
	}
}