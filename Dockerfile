
# Build the application from source
FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /key-server

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM alpine:3.20

WORKDIR /

COPY --from=build-stage /key-server /key-server

ENTRYPOINT ["/key-server"]