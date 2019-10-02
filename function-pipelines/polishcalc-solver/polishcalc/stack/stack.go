package stack

var chunkSize uint = 512

type Stack struct {
	size uint
	capacity uint

	stack []interface{}
}

func New() *Stack {
	s := new(Stack)
	s.capacity = chunkSize

	s.stack = make([]interface{}, chunkSize)
	return s
}

func (s *Stack) Push(v interface{}) {
	// Extend the stack capacity if necessary
	if s.size == s.capacity {
		s.capacity += chunkSize
		newChunk := make([]interface{}, chunkSize)

		// How efficient is this ?
		s.stack = append(s.stack, newChunk...)
	}
	s.stack[s.size] = v
	s.size++
}

func (s *Stack) Pop() interface{} {
	if s.size < 1 {
		return nil
	}
	s.size--

	// Set the value to nil to speed up garbage collection
	v := s.stack[s.size]
	s.stack[s.size] = nil
	return v
}

func (s *Stack) Peek() interface{} {
	if s.size < 1 {
		return nil
	}

	return s.stack[s.size - 1]
}

func (s *Stack) Size() uint {
	return s.size
}