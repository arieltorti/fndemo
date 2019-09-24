# Metrics

WARNING: THIS IS A WIP

The fnproject exposes metrics about the docker environment, cpu and memory consumption, number of running, queue, completed and failed function executions and many more.

In this demo I tried to add my own custom metrics, in order to export a metric it has to registered on a View and the View needs a go Context object, given that the functions
can't access it's parent context the only way I could add metrics was to register them inside listeners as such

```
	var MCustom = stats.Int64("demo/custom_metric", "Custom metric description", "")

	err := view.Register(
		&view.View{
			Name:        MCustom.Name(),
			Description: MCustom.Description(),
			Measure:     MCustom,
			TagKeys:     []tag.Key{},
			Aggregation: view.Sum(),
		},
	)
	if err != nil {
		log.Fatal("Error while registering METRICS VIEW:", err)
	}
	stats.Record(ctx, MCustom.M(int64(1)))

	return nil
```

One could register a listener that runs after app execution and send the app metrics on the response. The listener intercepts the response of the app, grabs the metrics, removes them from the response
and logs them.

The function presented at function-with-metrics was an attempt to add metrics from the fn function itself. It doesn't work
