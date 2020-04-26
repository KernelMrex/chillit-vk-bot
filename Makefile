build:
	go build -v ./cmd/chillitvkbot/.

run:
	go run -v ./cmd/chillitvkbot/. && ./chillitvkbot

run_dev:
	go build -v ./cmd/chillitvkbot/. && ./chillitvkbot

test:
	go test -v -race ./...

.DEFAULT_GOAL := run
