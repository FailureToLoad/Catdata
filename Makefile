.PHONY: run build clean generate

generate:
	templ generate

run: generate
	go run cmd/server/main.go

build: generate
	go build -o ./bin/server cmd/server/main.go

clean:
	rm -rf bin
	find . -name "*.templ.go" -type f -delete