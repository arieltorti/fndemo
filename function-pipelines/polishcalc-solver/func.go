package main

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"func/polishcalc"
	fdk "github.com/fnproject/fdk-go"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

type Equation struct {
	Equation string `json:"equation"`
}

type Response struct {
	Result int `json:"result"`
	Error string `json:"error"`
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	equation := new(Equation)
	if err := json.NewDecoder(in).Decode(equation); err != nil {
		log.Fatal(err)
		return
	}

	if equation.Equation == "" {
		if err := json.NewEncoder(out).Encode(&Response{Error: "Invalid equation"}); err != nil {
			log.Fatal(err)
		}
		return
	}

	result, err := solve(equation.Equation)
	if err != nil {
		if err := json.NewEncoder(out).Encode(&Response{Error: "Invalid equation"}); err != nil {
			log.Fatal(err)
		}
	}

	response := &Response{Result: result}

	if err := json.NewEncoder(out).Encode(response); err != nil {
		log.Fatal(err)
	}
}

func solve(equation string) (int, error) {
	res, err := polishcalc.Parse(equation)

	if err != nil {
		return 0, err
	}
	return res, nil
}