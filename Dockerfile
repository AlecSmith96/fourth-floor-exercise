FROM golang:alpine AS build


# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

# Setup folders
RUN mkdir /app
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY go.mod go.sum ./
COPY ./vendor ./vendor
COPY ./rest ./rest
COPY ./adapters ./adapters
COPY ./entities ./entities
COPY ./cmd ./cmd
COPY ./Makefile ./Makefile
COPY ./dev-config.yaml ./dev-config.yaml

RUN go mod download

RUN go build -o /build ./cmd/...

# # using distroless image for second stage as its very small and doesnt contain anything other 
# # than the binary to run
# FROM gcr.io/distroless/static-debian11 AS run-time

# COPY --from=build /build /build

EXPOSE 8080

CMD [ "/build", "--config", "dev-config.yaml" ]