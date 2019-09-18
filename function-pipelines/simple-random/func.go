package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	fdk "github.com/fnproject/fdk-go"
)

const ApiUrl = "http://caddy:2020"

type State struct {
	Times int `json:"times"`
}

type Msg struct {
	Msg []int `json:"msg"`
}

func _getMapWithDefault(m map[string]string, key string, default_value string) string {
	v, ok := m[key]
	if !ok {
		v = default_value
	}
	return v
}

func main() {
	fdk.Handle(fdk.HandlerFunc(randomNTimes))
}

func requestNewNumber(state *State) *Msg {
	jsonBody, _ := json.Marshal(state)
	resp, err := http.Post(fmt.Sprintf("%v/t/pipeline-demo/simple-random", ApiUrl),
		"application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		fmt.Println("Error ", err)
	}
	defer resp.Body.Close()
	respMsg := &Msg{}

	err = json.NewDecoder(resp.Body).Decode(respMsg)
	if err != nil {
		log.Fatal(err)
	}
	return respMsg
}

func randomNTimes(ctx context.Context, in io.Reader, out io.Writer) {
	fdkContext := fdk.GetContext(ctx)
	times, _ := strconv.Atoi(_getMapWithDefault(fdkContext.Config(), "times", "5"))
	state := &State{Times: 0}
	msg := &Msg{}

	err := json.NewDecoder(in).Decode(state)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	if state.Times < times {
		msg.Msg = append(msg.Msg, rand.Intn(10))
		state.Times += 1
		respMsg := requestNewNumber(state)

		msg.Msg = append(msg.Msg, respMsg.Msg...)
	}
	err = json.NewEncoder(out).Encode(&msg)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
