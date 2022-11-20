# generate wire dependency injection file
wire:
	cd ./cmd; \
	wire

# format all go files
fmt:
	go fmt ./internal/...

# build the binary and output to /dist directory
build:
	@mkdir -p ./dist
	go build -o ./dist ./cmd/...

# generate mocks for go interfaces
mocks:
	rm -rf ./mocks/
	mockery --dir=./internal -all -recursive=true -output=./mocks

# run all unit tests
unit:
	go test ./internal/...

# run the service locally in the current terminal
run-service:
	go run ./cmd --config dev-config.yaml

# build the docker image
docker-build:
	docker build -t ff-exercise .

# run the docker image
docker-run:
	docker run -dp 8080:8080 ff-exercise