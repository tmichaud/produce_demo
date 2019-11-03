#! /bin/bash

# Useful notes/command/etc. for developers

# To Run local (ie: not in a docker container)
#go get github.com/labstack/echo
#go build api/handlers/producehandler.go
#go build api/produce.go
#go build routers/router.go
#go run main.go

# To build the container
#! /bin/bash
#docker build -t produce_demo .

# To run the applicaton in the container
#docker run -it --rm produce_demo /bin/ash  # To Log into the machine
#docker run -p 8080:8080 produce_demo
