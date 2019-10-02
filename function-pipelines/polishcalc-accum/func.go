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

type ChannelResult struct {
	Result int
	Error error
}

type Equation struct {
	Equation string `json:"equation"`
}

type UserInput struct {
	Equations []string `json:"equations"`
}

type Response struct {
	Result int `json:"result"`
	Error string `json:"error"`
}

const ApiURL = "http://caddy:2020"
var solveURL string

func main() {
	solveURL = fmt.Sprintf("%v/t/pipeline-demo/polishcalc-solver", ApiURL)
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	// Get the equations from the user
	equations := new(UserInput)
	if err := json.NewDecoder(in).Decode(equations); err != nil {
		errorResponse(err, out)
		return
	}

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
		errorResponse(err, out)
	}
	return
}

func solveEquation(equation string, ch chan<-*ChannelResult) {
	// Encode the equation in an io.Reader
	buf := bytes.NewBuffer(nil)
	payload := &Equation{Equation: equation}
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		ch <- &ChannelResult{Error: err}
		return
	}

	resp, err := http.Post(solveURL, "application/json", buf)
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

	// If the solver sends an error it will encode it in the response
	// We need to grab it from the response and pass it along
	if solverResponse.Error != "" {
		err = errors.New(solverResponse.Error)
	}
	ch <- &ChannelResult{Result: solverResponse.Result, Error: err}
}

func errorResponse(err error, out io.Writer) {
	log.Print(err)
	errMsg := &Response{Error: err.Error()}
	if err := json.NewEncoder(out).Encode(&errMsg); err != nil {
		log.Panic(err)
	}
}
