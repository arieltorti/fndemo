package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"

	"func/polishcalc"
	fdk "github.com/fnproject/fdk-go"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

type UserInput struct {
	Equation string `json:"equation"`
}

type Response struct {
	Result int `json:"result"`
	Error string `json:"error"`
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	// Get the equation from the user
	equation := new(UserInput)
	if err := json.NewDecoder(in).Decode(equation); err != nil {
		errorResponse(err, out)
		return
	}

	if equation.Equation == "" {
		errorResponse(errors.New("invalid equation"), out)
		return
	}

	// Solve the equation
	result, err := solve(equation.Equation)
	if err != nil {
		errorResponse(err, out)
		return
	}

	response := &Response{Result: result}
	if err := json.NewEncoder(out).Encode(response); err != nil {
		errorResponse(err, out)
	}
}

func errorResponse(err error, out io.Writer) {
	log.Print(err)
	errMsg := &Response{Error: err.Error()}
	if err := json.NewEncoder(out).Encode(&errMsg); err != nil {
		log.Panic(err)
	}
}

func solve(equation string) (int, error) {
	return polishcalc.Parse(equation)
}