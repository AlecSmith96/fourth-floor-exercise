# generate wire dependency injection file
wire:
	cd ./cmd \
	wire

# build the binary and output to /dist directory
build:
	@mkdir -p ./dist
	go build -o ./dist ./cmd/...

# run the service locally in the current terminal
run-service:
	go run ./cmd --config dev-config.yaml

# build the docker image
docker-build:
	docker build -t ff-exercise .

# run the docker image
docker-run:
	docker run -p 8080:8080 ff-exercise