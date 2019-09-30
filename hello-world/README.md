# Hello World Example

This example shows to create a fn application and deploy a function to it.
First we create the app, then the function and finally we will register a trigger for it.

## Create Fn App
In order to create an application we have to be running a Fn Server. To do this locally you can run
```
fn start
```
after installing Fn.

Now we lets create an app called `hello-demo` with:

```
fn create app hello-demo
```

Applications can contain many functions which in turn contain triggers.

## Create Fn function

The command `fn init` lets you create function boilerplate code, the --runtime arguments tell it which language to use.
We will create a function called `hello-world` by running:

```
fn init --runtime go hello-world
```

## Deploy the function

Now that we have a function we can tell Fn to deploy it by using the `fn deploy` command. The deploy command builds, tags and registers the docker image (All Fn functions are docker images)
and all of their configurations and metadata are stored in a database, by default that database is sqlite3. Lets deploy our hello-world function:

```
cd hello-world/
fn deploy --app hello-demo  --local
```

The --local flag tells Fn not to push the generated docker image to the registry.

## Invoking the function

Our functions is now registered in Fn. To check this we can run:
```
> fn ls functions hello-demo

NAME		IMAGE			ID
hello-world	demo/hello-world:0.0.1	01DKF29G46NG8G00GZJ0000002
```

If the run the commands with DEBUG=1 env variable we can see the process of the cli commands.
For example if we run

```
fn invoke hello-demo hello-world
```

It will run the function hello-world inside the hello-demo app.
In the output of the command we can see it makes many requests. The flow of a function invoke is :

* Get the id of the function `hello-world` inside the app `hello-demo`
* Using that id get its trigger id
* Make a POST request to /invoke/{trigger_id}

The request should return a message saying "Hello World", but we also have access to the requests headers and parameters, for example if we do:

```
echo -n '{"name": "Mario"}' | fn invoke hello-demo hello-world --content-type application/json
{"message":"Hello Mario"}
```

Next we will look at the func.yaml file and how to create function triggers
