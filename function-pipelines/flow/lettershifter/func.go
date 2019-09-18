package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fnproject/fdk-go"
	flows "github.com/fnproject/flow-lib-go"
)

func init() {
	flows.RegisterAction(toUpperCase)
	flows.RegisterAction(shiftLetters)
}

func main() {
	fdk.Handle(myHandler())
}

type UserInput struct {
	Value string `json:"value"`
}

func myHandler() fdk.Handler {
	return flows.WithFlow(
		fdk.HandlerFunc(func(ctx context.Context, r io.Reader, w io.Writer) {
			input := &UserInput{Value: "Hello World"}
			json.NewDecoder(r).Decode(&input)

			cf := flows.CurrentFlow().CompletedValue(input.Value)
			valueCh, errorCh := cf.ThenApply(toUpperCase).ThenApply(shiftLetters).Get()
			printResult(w, valueCh, errorCh)
		}))
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