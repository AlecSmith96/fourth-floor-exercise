# Ideally a multi-stage build would be used here to allow for the
# binary to be built in a separate image that contains the necessary
# tools, before then copying the binary over into a separate image such
# as a scratch or distroless image. This would lead to a much smaller image
# and be safer as no uneeded tools would be included in it.
FROM golang:alpine AS build


# install git
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

# setup folder
RUN mkdir /app
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY go.mod go.sum ./
COPY ./vendor ./vendor
COPY ./internal/rest ./internal/rest
COPY ./internal/adapters ./internal/adapters
COPY ./internal/entities ./internal/entities
COPY ./cmd ./cmd
COPY ./Makefile ./Makefile
COPY ./dev-config.yaml ./dev-config.yaml

RUN go mod download

RUN go build -o /build ./cmd/...

EXPOSE 8080

CMD [ "/build", "--config", "dev-config.yaml" ]