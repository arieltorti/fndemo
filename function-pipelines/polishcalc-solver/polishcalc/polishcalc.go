package polishcalc

import (
	"errors"
	"fmt"
	"func/polishcalc/stack"
	"strconv"
	"strings"
)

// TODO: Implement log levels

type Parser struct {
	tokens []string
	stack *stack.Stack
}

func New() *Parser {
	p := new(Parser)
	p.stack = stack.New()
	return p
}

func invalidEquationError(equation string) error {
	return errors.New(fmt.Sprintf("the equation `%v` is invalid", equation))
}

func Parse(equation string) (int, error) {
	if equation == "" {
		return 0, invalidEquationError(equation)
	}

	tokens := strings.Fields(equation)
	if len(tokens) < 1 {
		return 0, invalidEquationError(equation)
	}

	parser := New()
	parser.tokens = tokens

	res, err := parseTokens(parser)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func parseTokens(p *Parser) (int, error) {
	for _, v := range p.tokens {
		token, err := getTokenValue(v)
		if err != nil {
			return 0, err
		}

		err = pushToken(p, token)
		if err != nil {
			return 0, err
		}
	}

	// If the equation is valid, we have to end up with only
	// one element in the stack: the result.
	if p.stack.Size() != 1 {
		return 0, errors.New("the stack has more than one element")
	}
	res, ok := p.stack.Pop().(int)
	if !ok {
		return 0, errors.New("invalid equation, response is not a number")
	}
	return res, nil
}

func getTokenValue(token string) (interface{}, error) {
	v, err := strconv.Atoi(token)
	if err != nil {
		//	Token is not an int, check if its an operation
		op, ok := tokenToOperation(token)
		if !ok {
			return nil, errors.New(fmt.Sprintf("invalid token %v", token))
		}

		return op, nil
	}
	return v, nil
}

func pushToken(p *Parser, token interface{}) error {
	operation, ok := token.(Operation)
	if !ok {
		// Case: Token is an int
		value, _ := token.(int)
		_, isInt := p.stack.Peek().(int)

		//	If the last item is an int, and we are pushing another int
		//  we can compute an operation, because all of them are binary
		//  this process repeats until we cant do it anymore
		for isInt {
			op1 := p.stack.Pop().(int)
			operation, ok = p.stack.Pop().(Operation)
			if !ok {
				// This case can happen with invalid equations such as `+ 1 2 4`
				// Here we will push +, then 1, when we get to 2 we compute 1 + 2 = 3 and push it
				// by the time we receive the number 4 the stack only has a single token, 3,
				// we pop it out, and when we pop the operation it is nil because the stack is empty
				// and the conversion fails.
				return errors.New("invalid equation")
			}

			value = operation.Do(op1, value)
			_, isInt = p.stack.Peek().(int)
		}

		p.stack.Push(value)
	} else {
		// Case: Token is an operation
		p.stack.Push(operation)
	}

	return nil
}
