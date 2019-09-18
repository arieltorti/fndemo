package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"time"

	fdk "github.com/fnproject/fdk-go"
	flows "github.com/fnproject/flow-lib-go"
)

func init() {
	flows.RegisterAction(GenerateRandomNumber)
	flows.RegisterAction(GenerateUntilN)
}

func main() {
	fdk.Handle(myHandler())
}

func myHandler() fdk.Handler {
	return flows.WithFlow(
		fdk.HandlerFunc(func(ctx context.Context, r io.Reader, w io.Writer) {
			cf := flows.CurrentFlow().CompletedValue("")
			cf = GenerateUntilN(cf, w)
			valueCh, errorCh := cf.Get()
			printResult(w, valueCh, errorCh)
		}))
}

func GenerateUntilN(cf flows.FlowFuture, w io.Writer) flows.FlowFuture {
	results := make([]int, 10)
	for i := 0; i < 10; i++ {
		cf := cf.ThenApply(GenerateRandomNumber)
		valueCh, _ := cf.Get()

		generatedInt := <- valueCh
		results[i] = generatedInt.(int)
	}

	return flows.CurrentFlow().CompletedValue(results)
}

func GenerateRandomNumber() int {
	return rand.Intn(10)
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