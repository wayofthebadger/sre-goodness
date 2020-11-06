# sre-goodness 
SRE build fun times

This is my repo for the SRE funtimes test

Running from Command line locally using Docker

# Prerequisites

Export the required variables to your environment:
```
$ export Symbol=MSFT
$ export NDAY=3
```

Create an apikey file
```
$ echo -n '<apikey>' > confg/config.txt
```
Run using
```
go run stocks.go
```
NOTE: on a managed OSX machine you'll need to allow the app to open a port.

Go to http://localhost:8080

# Build container

To build this please first ensure you've got docker installed on your local system

Uncomment the variables fields in the Dockerfile (add the comments again for k8s deployment)

```
ENV Symbol=MSFT \
    NDAYS=4

... ...

# Create config for testing, removed for prod
COPY ./config/config.txt /dist/config/config.txt
... ...
```

Ensure you're in the root of this Git repo

Run the following command to build the image
```
$ docker build -t stocks-app .
```

Run the next command to spawn the container

```
docker run -p 8080:8080 stocks-app:latest
```

Go to http://localhost:8080

Be (very very) mildy impressed/horrified at an IT system's engineer's first ever real try at a GoLang app

Many thanks
