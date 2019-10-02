package polishcalc

// TODO: Tests

type Operation interface {
	Do(op1 int, op2 int) int
}

type Summation struct { }

func (s *Summation) Do(op1 int, op2 int) int {
	return op1 + op2
}

type Difference struct { }

func (s *Difference) Do(op1 int, op2 int) int {
	return op1 - op2
}

type Multiplication struct {
	op1 int
	op2 int
}

func (s *Multiplication) Do(op1 int, op2 int) int {
	return op1 * op2
}

type Division struct { }

func (s *Division) Do(op1 int, op2 int) int {
	return op1 / op2
}

func tokenToOperation(token string) (Operation, bool) {
	switch token {
	case "+":
		return new(Summation), true
	case "*":
		return new(Multiplication), true
	case "-":
		return new(Difference), true
	case "/":
		return new(Division), true
	default:
		return nil, false
	}
}