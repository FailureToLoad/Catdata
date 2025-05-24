.PHONY: run build clean generate

dependencies:
	curl -L -o "./tailwindcss" "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64"
	chmod +x "./tailwindcss"
	
	curl -sLO "https://github.com/a-h/templ/releases/latest/download/templ_Linux_x86_64.tar.gz"
	tar -xzf "templ_Linux_x86_64.tar.gz"
	rm "templ_Linux_x86_64.tar.gz"
	chmod +x "templ"
	go mod tidy

templ:
	$(error templ missing, run make dependencies)

tailwindcss:
	$(error tailwind missing, run make dependencies)

generate: templ
	./templ generate

css: tailwindcss
	./tailwindcss -i ./input.css -o ./static/styles.css --minify

css/watch: tailwindcss
	./tailwindcss -i ./input.css -o ./assets/css/output.css --watch

run: generate
	go run cmd/server/main.go

build: generate
	go build -o ./bin/server cmd/server/main.go

clean:
	rm -rf bin
	find . -name "*.templ.go" -type f -delete