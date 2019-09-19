# Ability to choose build image

By setting the `build_from` configuration on the `func.yaml` file we can specify which docker image to use in the building stage

```
build_image: python:3.6
```

If we start the server and deploy

```
fn start
fn create app demo-build

cd build-from/
fn deploy --app demo-build --local
```

On the generate docker file we can see python:3.6 is being used as build the image:

```
FROM python:3.6 as build-stage
```
