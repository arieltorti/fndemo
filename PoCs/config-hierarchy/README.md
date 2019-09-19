# Configuration Hierarchy

This small proof of concept shows how fn configuration hierarchy works.
The function is a hello world that returns `Hello World` by default or `Hello <input>` if an input is given.

The name of the input parameter can be configured, and as you will see, function configuration has priority over app configuration.

## Running

First start the fn server and create an app.

```
fn start
fn create app demo-python
```

And deploy the function

```
cd config-hierarchy/
fn deploy --app demo-python --local
```

Now lets invoke the function without any configuration

```
fn invoke demo-python config-hierarchy
```

returns `Hello World`, if we were to give it a name parameter:

```
echo '{"name": "Moon"}' | fn invoke demo-python config-hierarchy
```

it would return `Hello Moon`.
Now to the configuration, lets add a config key `param` to the app.

```
fn config app demo-python "param" "planet"
```

Now the function will not read the `name` param, instead it will read the `planet` one

```
echo '{"planet": "Moon"}' | fn invoke demo-python config-hierarchy
```

But remember that function config has priority over app config, if we change the configuration on the function as such:

```
fn config fn demo-python config-hierarchy "param" "lastname"
```

Then the app config will be overwritten by the function config resulting in 

```
echo '{"lastname": "Mars"}' | fn invoke demo-python config-hierarchy
```

returning `Hello Mars`
