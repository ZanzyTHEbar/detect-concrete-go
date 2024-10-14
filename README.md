# Detect Cracks in Concrete

This repo shows that working in go with opencv is possible. The code is a simple example of how to detect cracks in concrete.

## Requirements

- Docker
- GoCV (https://gocv.io/)
- Go 1.23

### How to install GoCV

Ensure that you build the gocv docker image, follow the instructions in the [gocv repo](https://github.com/hybridgroup/gocv).

You will need to override the Makefile target `docker` to add the following line:

```bash
docker build -t gocv --build-arg OPENCV_VERSION=$(OPENCV_VERSION) --build-arg GOVERSION=$(GOVERSION) .
```

Then run the following command:

```bash
make docker
```

## How to run

The docker engine must be running, then run the following command:

```bash
make build
```

This will build a docker image with the name `concrete`, containing the code to detect cracks in concrete.

To run the code, execute the following command:

```bash
make run
```

This will run docker container in interactive mode.

Inside of the container, run the following command:

```bash
./bin/main
```
