build:
	go build -v ./cmd/chillitvkbot/.

run:
	go run -v ./cmd/chillitvkbot/. -config_path=configs/config.yaml

run_dev:
	go run -v ./cmd/chillitvkbot/. -config_path=configs/config.yaml.devel

test:
	go test -v -race ./...

.DEFAULT_GOAL := run
