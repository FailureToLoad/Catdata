.PHONY: run build clean generate

setup:
	chmod +x ./install-tailwind.sh
	chmod +x ./install-templ.sh
	./install-tailwind.sh
	./install-templ.sh
	go mod tidy

generate:
	templ generate

css:
	npx tailwindcss -i ./input.css -o ./static/styles.css --minify

css/watch:
	npx tailwindcss -i ./input.css -o ./assets/css/output.css --watch

run: generate
	go run cmd/server/main.go

build: generate
	go build -o ./bin/server cmd/server/main.go

clean:
	rm -rf bin
	find . -name "*.templ.go" -type f -delete