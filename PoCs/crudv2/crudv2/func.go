package main

import (
	"context"
	"encoding/json"
	"io"
	"log"

	fdk "github.com/fnproject/fdk-go"
)
type Message struct {
	Trigger string `json:"trigger"`
	Method string `json:"post"`
}

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	msg := new(Message)

	log.Println(fdk.GetContext(ctx))
	method := fdk.GetContext(ctx).(fdk.HTTPContext).RequestMethod()
	trigger := fdk.GetContext(ctx).(fdk.HTTPContext).RequestURL()
	if err := json.NewDecoder(in).Decode(&msg); err != nil {
		log.Println("The was an error while parsing the request ", err)
	}

	msg.Method = method
	msg.Trigger = trigger
	json.NewEncoder(out).Encode(&msg)
}
