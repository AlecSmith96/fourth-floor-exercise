# build the binary and output to /dist directory
build:
	@mkdir -p ./dist
	go build -o ./dist ./cmd/...

# run the service locally in the current terminal
run-service:
	go run ./cmd --config dev-config.yaml

docker-build:
	docker build -t ff-exercise .

docker-run:
	docker run -dp 8080:8080 ff-exercise