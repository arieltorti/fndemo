package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	fdk "github.com/fnproject/fdk-go"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tags"
)

var MNameLength = stats.Int64("demo/name_length", "Custom metric description", "")

func main() {
	err := view.Register(
                &view.View{
                        Name:        MNameLength.Name(),
                        Description: MNameLength.Description(),
                        Measure:     MNameLength,
                        TagKeys:     []tag.Key{},
                        Aggregation: view.Sum(),
                },
        )
        if err != nil {
                log.Fatal("Error while registering METRICS VIEW:", err)
        }


	fdk.Handle(fdk.HandlerFunc(myHandler))
}

type Person struct {
	Name string `json:"name"`
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	p := &Person{Name: "World"}
	json.NewDecoder(in).Decode(p)

	msg := struct {
		Msg string `json:"message"`
	}{
		Msg: fmt.Sprintf("Hello %s", p.Name),
	}

	stats.Record(ctx, MNameLength.M(int64(len(p.Name))))

	row, _ := view.RetrieveData("demo/name_length")
	json.NewEncoder(out).Encode(&msg)
}
