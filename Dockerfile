# Build the application from source
FROM golang:1.22.2 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

#USER root

RUN mkdir -p ./src
RUN mkdir -p ./src/db
RUN mkdir -p ./src/middleware
RUN mkdir -p ./src/util

COPY src/db/ ./src/db
COPY src/middleware/ ./src/middleware
COPY src/util/ ./src/util
COPY src/*.go ./src/


#USER nonroot:nonroot

RUN ["/bin/bash", "-c", "GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap ./src/**.go"]

# Run the tests in the container
#FROM build-stage AS run-test-stage
#RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/bootstrap /bootstrap

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/bootstrap"]