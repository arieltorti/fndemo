package main

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	fdk "github.com/fnproject/fdk-go"
	flows "github.com/fnproject/flow-lib-go"
)

// Who calls me ?
func init() {
	flows.RegisterAction(returnHelloWorld)
	flows.RegisterAction(toUpperCase)
	flows.RegisterAction(shiftLetters)
}

func main() {
	fdk.Handle(myHandler())
}

type Person struct {
	Name string `json:"name"`
}

func myHandler() fdk.Handler {
	return flows.WithFlow(
		fdk.HandlerFunc(func(ctx context.Context, r io.Reader, w io.Writer) {
			cf := flows.CurrentFlow().CompletedValue("I am not used")
			valueCh, errorCh := cf.ThenApply(returnHelloWorld).ThenApply(toUpperCase).ThenApply(shiftLetters).Get()
			printResult(w, valueCh, errorCh)
		}))
}

func returnHelloWorld() string {
	return "Hello World"
}

func toUpperCase(s string) string {
	return strings.ToUpper(s)
}

func shiftLetters(s string) string {
	tmp := []rune(s)
	for k, v := range tmp {
		tmp[k] = v + 2
	}
	return string(tmp)
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