#########################
# Stage 1
FROM golang:latest as builder

# Copy in the go files 
COPY . /go/src/produce_demo

# Set up the working directory
WORKDIR /go/src/produce_demo

# Get Echo
RUN go get github.com/labstack/echo

# Build the files
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v -o /out/echo_app

# Run the tests with coverage
RUN go test --cover ./...

#########################
# Stage 2
#FROM alpine:latest as release # Use this if you want to view inside the container
#     alpine requires you not to use /bin/bash but rather /bin/ash
FROM scratch as release

COPY --from=builder /out/ /out/

ENTRYPOINT ["/out/echo_app"]

EXPOSE 8080


