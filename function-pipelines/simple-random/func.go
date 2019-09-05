package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"

	fdk "github.com/fnproject/fdk-go"
)

const ApiUrl = "http://localhost:8080"

type State struct {
	Times int `json:"times"`
}

type Msg struct {
	Msg []int `json:"msg"`
}

func main() {
	fdk.Handle(fdk.HandlerFunc(randomNTimes))
}

func requestNewNumber(state *State) *Msg {
	jsonBody, _ := json.Marshal(state)
	resp, err := http.Post(fmt.Sprintf("%v/t/demo/simple-random", ApiUrl),
		"application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		fmt.Println("Error ", err)
	}
	respMsg := &Msg{}

	// TOASK: Should I close body ??
	json.NewDecoder(resp.Body).Decode(respMsg)
	return respMsg
}

func randomNTimes(ctx context.Context, in io.Reader, out io.Writer) {
	fdkContext := fdk.GetContext(ctx)
	times, _ := strconv.Atoi(fdkContext.Config()["times"])
	state := &State{Times: 0}
	msg := &Msg{}

	json.NewDecoder(in).Decode(state)

	if state.Times < times {
		msg.Msg = append(msg.Msg, rand.Intn(10))
		state.Times += 1
		respMsg := requestNewNumber(state)

		msg.Msg = append(msg.Msg, respMsg.Msg...)
	}
	json.NewEncoder(out).Encode(&msg)
}
