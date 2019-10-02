# Fn Flow pipelines Example

[Fn Flow](https://github.com/fnproject/flow) is an extension provided by the Fn Project to build workflows.
In this example we will create some apps creating a pipeline manually, and then creating it using Fn Flow.

## Starting the server

We asume that the `caddy` hostname points to the local machine, then we start the Fn and Fn Flow servers using:

```
make start
```

Configure the Fn CLI url:
```
export FN_API_URL=http://caddy:2020
```

In order to deploy function we must first create an application:
```
fn ls app
fn create app pipeline-demo
```

## Deploying pipelines manually

We provided a sequential function with generates N random numbers, one after another by manually invoking the function again, to
deploy this function we can execute:

```
cd simple-random/
fn deploy --app pipeline-demo --local
fn invoke pipeline-demo simple-random
```

The number of generated random numbers depends on a configuration variable named `times`, if we change it we can see control how many
random variables are created, and we can even start the function midway by giving it a `times` value:

```
fn config function pipeline-demo simple-random "times" "10"
fn inspect function pipeline-demo simple-random
echo '{"times": 7}' | fn invoke pipeline-demo simple-random --content-type application/json
```

The execution of this function is sequential, and the communication between the different stages is donw manually.


For a different example we deploy a function that solves equations, and another one that solves many equations and adds their results:

Solver will solve individual functions:
```
cd polishcalc-solver/
fn deploy --app pipeline-demo --local

echo -n '{"equation": "- * / 15 - 7 + 1 1 3 + 2 + 1 5"}' | fn invoke pipeline-demo polishcalc-solver --content-type application/json
```

Accumulator will call solvers in parallel and add the results:
```
cd polishcalc-accum/
fn deploy --app pipeline-demo --local

echo -n '{"equations": ["239", "* 2 + 9 512", "+ 3 2", "- * / 15 - 7 + 1 1 3 + 2 + 1 5"]}' | fn invoke pipeline-demo polishcalc-accum --content-type application/json
```

The stages of this function run in parallel.

## Deploying pipelines using Fn Flow

Now that we have implemented the pipelines manually, we can use Fn Flow to create the same functions and compare the
advantages or disadvantages of using it.

To the Fn Server functions that use Fn Flow are just normal functions, so they deploy and execute in the same fashion. As a first example
we deploy a function that uppercases and shifts the string we give it:

```
cd flow/
cd lettershifter/
fn deploy --app pipeline-demo --local

echo '{"value": "Shifted Hello World"}' | fn invoke pipeline-demo lettershifter
```

The invoke command probably did not work for you, that's because in order for Fn Flow to work we have to tell it where its server is located.
We do this by setting the config variable `COMPLETER_BASE_URL`:

```
fn config app pipeline-demo COMPLETER_BASE_URL http://caddy:8081
```


The random generator function is also implemented as an example using Fn Flow:

```
cd flow-random/
fn deploy --app pipeline-demo --local

fn invoke pipeline-demo flow-random
```

Regarding the implementation of the functions, we can see the Fn Flow function was a bit easier to do, and is easier to read, but we still have
create the logic for coordinating the different flows and gather their results.


Lastly we deploy the solver function which runs in parallel.
Note that this function depends on the polishcalc-solver function and before we invoke it we have to give it the polishcalc-solver ID on the line 17
of the func.go file.

Now we deploy and invoke:

```
cd flow-polishcalc/
fn deploy --app pipeline-demo --local

echo -n '{"equations": ["239", "* 2 + 9 512", "+ 3 2", "- * / 15 - 7 + 1 1 3 + 2 + 1 5"]}' | fn invoke pipeline-demo flow-polishcalc
```

In this case the implementation is not simpler than the one with Fn Flow, and we also have the added complexity of having to update the
polishcalc-solver ID every time we update it, read: we don't get much benefits from using Fn Flow, but we have some drawbacks
