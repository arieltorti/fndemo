package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	fdk "github.com/fnproject/fdk-go"
	flows "github.com/fnproject/flow-lib-go"
)

const funcID  = "01DMK1Q7JKNGAG00GZJ0000011"

type UserInput struct {
	Equations []string `json:"equations"`
}

type Response struct {
	Result int `json:"result"`
	Error string `json:"error"`
}

type Equation struct {
	Equation string `json:"equation"`
}

func init() {
	flows.RegisterAction(SplitEquations)
	flows.RegisterAction(SolveEquation)
}

func main() {
	flows.Debug(true)
	fdk.Handle(myHandler())
}

func myHandler() fdk.Handler {
	return flows.WithFlow(
		fdk.HandlerFunc(func(ctx context.Context, r io.Reader, w io.Writer) {
			equations := new(UserInput)
			if err := json.NewDecoder(r).Decode(equations); err != nil {
				errorResponse(err, w)
				return
			}

			if len(equations.Equations) == 0 {
				errorResponse(errors.New("no equations or invalid equations were provided"), w)
				return
			}

			cf := flows.CurrentFlow().CompletedValue(equations.Equations).ThenApply(SplitEquations)
			valueCh, errorCh := cf.Get()
			printResult(w, valueCh, errorCh)
		}))
}

func SplitEquations(equations []string) int {
	cf := flows.CurrentFlow()

	ch := make(chan int)
	for _, v := range equations {
		go func(cf flows.Flow, ch chan <- int) {
			ff := cf.CompletedValue(v).ThenApply(SolveEquation)
			valueCh, _ := ff.Get()
			value := <- valueCh
			intValue, _ := value.(int)

			ch <- intValue
		}(cf, ch)
	}

	result := 0
	for range equations {
		result += <- ch
	}
	return result
}

func SolveEquation(equation string) int {
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(&Equation{Equation: equation})

	httpRequest := &flows.HTTPRequest{Method: "POST", Body: buf.Bytes()}

	cf := flows.CurrentFlow().InvokeFunction(funcID, httpRequest)
	valueCh, _ := cf.Get()

	response := <- valueCh
	solverResponse := new(Response)
	httpResponse := response.(*flows.HTTPResponse)

	json.Unmarshal(httpResponse.Body, &solverResponse)
	fmt.Println("SOLVER RESPONSE IS", solverResponse)
	return solverResponse.Result
}

func printResult(w io.Writer, valueCh chan interface{}, errorCh chan error) {
	select {
	case value := <-valueCh:
		fmt.Fprintf(w, "Flow succeeded with value %v", value)
	case err := <-errorCh:
		fmt.Fprintf(w, "Flow failed with error %v", err)
	case <-time.After(time.Minute * 1):
		fmt.Fprintf(w, "Timed out!")
	}
}

func errorResponse(err error, out io.Writer) {
	log.Print(err)
	errMsg := &Response{Error: err.Error()}
	if err := json.NewEncoder(out).Encode(&errMsg); err != nil {
		log.Panic(err)
	}
}