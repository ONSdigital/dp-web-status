all: generate
	go get github.com/mitchellh/gox
	gox -osarch="linux/amd64" -output="dp_web_status_linux_amd64"
	zip dp-web-status.zip dp_web_status_linux_amd64 Dockerfile config.yml

run: generate
	go run main.go

generate:
	go generate ./...

.PHONY: all run generate
