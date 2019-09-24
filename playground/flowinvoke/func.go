package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	fdk "github.com/fnproject/fdk-go"
	flows "github.com/fnproject/flow-lib-go"
)

func main() {
	fdk.Handle(myHandler())
}

type Response struct {
	Msg string `json:"message"`
}

func myHandler() fdk.Handler {
	return flows.WithFlow(
		fdk.HandlerFunc(func(ctx context.Context, r io.Reader, w io.Writer) {
			req := &flows.HTTPRequest{Method: "POST", Body: []byte("{}")}

			cf := flows.CurrentFlow().InvokeFunction("01DJAVB4T9NG8G00GZJ0000002", req)
			valueCh, errorCh := cf.Get()

			select {
			case value := <-valueCh:
				resp, ok := value.(*flows.HTTPResponse)
				if !ok {
					panic("received unexpected value from the server")
				}
				var gr Response
				json.Unmarshal(resp.Body, &gr)
				fmt.Fprintf(w, "Got HTTP status %v and payload %v", resp.StatusCode, gr)
			case err := <-errorCh:
				fmt.Fprintf(w, "Flow failed with error %v", err)
			case <-time.After(time.Minute * 1):
				fmt.Fprintf(w, "Timed out!")
			}
		}))
}
