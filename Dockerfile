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

# Note: The CGO_ENABLED=0 is needed for the built image to run on alpine on the next stage. Otherwise it fails, probably
# with an "./bootstrap not found" error.
# Read more: https://www.reddit.com/r/golang/comments/pi97sp/what_is_the_consequence_of_using_cgo_enabled0/
# Turns out, apparently Go compiler assumes glibc, while alpine (image used in next stage) uses musl
# Is there a way to specify the libc to use during compilation?
# So apparently the glibc expectation vs musl in alpine is a very known issue, and basically depends on the
# libraries/dependencies used in the project - some things expect glibc, so unexpected behavior can happen
# (https://hub.docker.com/_/golang on the Image Variants section)
RUN ["/bin/bash", "-c", "GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap ./src/**.go"]

RUN ls .

# Run the tests in the container
#FROM build-stage AS run-test-stage
#RUN go test -v ./...

# Deploy the application binary into a lean image
FROM public.ecr.aws/lambda/provided:al2023 AS build-release-stage

WORKDIR /app

COPY --from=build-stage /app/bootstrap /app/bootstrap

RUN ls .

EXPOSE 8080

ENTRYPOINT ./bootstrap