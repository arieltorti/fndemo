# [Fn server behind TLS issue](https://github.com/fnproject/fn/issues/1537)

There is an issue by the user prologic that states it isn't possible to use `fn invoke` when the fn server is running behind a TLS load balancer.

## Is this true ?

We are going to try and reproduce this issue. In the current folder there's a `docker-compose.yml` file which provides a simple architecture, Caddy as the LB, and a single fnserver.
The `docker-compose.full.yml` contains the complete architecture which should be running on a real production server, a postgres database where the data is stored, Caddy still doing the load balancing, and 
fnserver split into runners, load balancer and api.

We suppose that the caddy hostname points to localhost and that we're using certificates that we trust (If not true see [Notes](#notes)).
First we will run the simpler architecture, the fn services and caddy need SSL Server Certificates to work, if you want to use lets encrypt feel free to set it up, but in this example we are going to 
create our own self-signed certificates by running:

```
./create_certs.sh
```

And deploy the services using:

```
make start
```

To tell the fn cli to make requests to our server we have to set the following variables:

```
export FN_API_URL=https://caddy:2020
export FN_INSECURE=1  # Ignore tls errors
```

### Create an app and deploy a function

```
fn create app demo
fn deploy --app demo --local  # Inside the directory of a function
```

Now that there's an app and a function, lets invoke it and see if we can reproduce the issue:

```
fn invoke demo hello-world
```

Success ! It works. If we set `DEBUG=1` we can see the requests are indeed being made over https.

### What about the full architecture ?

Before we deploy the full architecture we have to stop be old one:

```
make stop
make start-deploy
```

This version uses a postgres database stored in a volume, which means our data will persist as long as we don't purge the volume.
Caddy acts as a load balancer and a reverse proxy redirecting requests to either the fn load balancer or the fn api. In theory we would have many instances of the fn lb running and 
Caddy would load balance them, which in turn load balance the runners.
Replaying the old commands we get:

```
fn create app demo
fn deploy --app demo --local
fn invoke demo hello-world
```

And it works as well. On a side note you can try and not take my word for it, but if you were to install fncli inside the containers and set the correct `FN_API_URL` it would also work.

## Why did the issue happen then ?

Most probably the user prologic made a configuration error. When we deploy a function, the trigger is stored as an annotation inside the function ([More on annotations](#todo)), if the function was deployed without tls then the trigger endpoint is saved using `http://` instead of `https://`.
This means that if you deploy a function without tls and later try to use it behind tls then it won't work.
