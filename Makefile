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
# Ideally would also add calculation for total code coverage which would be printed to console
unit:
	go test ./internal/...

# Ideally we would also have a `make functional` make target here to run
# a set of functional tests that would test end to end functionality and 
# help catch integration issues between other services and APIs being used.
# due to time constraints I was unable to do this.

# run the service locally in the current terminal
run-service:
	go run ./cmd --config dev-config.yaml

# build the docker image
docker-build:
	docker build -t ff-exercise .

# run the docker image
docker-run:
	docker run -dp 8080:8080 ff-exercise