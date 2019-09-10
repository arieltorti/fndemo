package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	fdk "github.com/fnproject/fdk-go"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

type ChannelResult struct {
	Result int
	Error error
}

type Equation struct {
	Equation string `json:"equation"`
}

type Equations struct {
	Equations []string `json:"equations"`
}

type Response struct {
	Result int `json:"result"`
	Error string `json:"error"`
}

const solveUrl = "http://localhost:8080/t/demo/polishcalc-solver"

func solveEquation(equation string, ch chan<-*ChannelResult) {
	buf := bytes.NewBuffer(nil)
	payload := &Equation{Equation: equation}
	json.NewEncoder(buf).Encode(payload)

	resp, err := http.Post(solveUrl, "application/json", buf)
	if err != nil {
		ch <- &ChannelResult{Error: err}
		return
	}
	defer resp.Body.Close()
	solverResponse := new(Response)

	if err := json.NewDecoder(resp.Body).Decode(solverResponse); err != nil {
		ch <- &ChannelResult{Error: err}
		return
	}

	err = nil
	if solverResponse.Error != ""{
		err = errors.New(solverResponse.Error)
	}
	ch <- &ChannelResult{Result: solverResponse.Result, Error: err}
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	// Get the equations from the user
	equations := new(Equations)
	json.NewDecoder(in).Decode(equations)

	if len(equations.Equations) == 0 {
		errorResponse(errors.New("no equations or invalid equations were provided"), out)
		return
	}

	results := make(chan *ChannelResult, len(equations.Equations))
	for _, v := range equations.Equations {
		go solveEquation(v, results)
	}

	// For now we only implement sum of values
	accumResult := 0
	for range equations.Equations {
		res := <-results

		if res.Error != nil {
			errorResponse(res.Error, out)
			return
		}
		accumResult += res.Result
	}

	respMsg := Response{Result: accumResult}
	if err := json.NewEncoder(out).Encode(respMsg); err != nil {
		log.Panic(err)
	}
	return
}

func errorResponse(err error, out io.Writer) {
	log.Print(err)
	errMsg := &Response{Error: err.Error()}
	if err := json.NewEncoder(out).Encode(&errMsg); err != nil {
		log.Panic(err)
	}
}