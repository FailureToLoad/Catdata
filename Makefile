.PHONY: run build clean generate

dependencies:
	curl -L -o "./tailwindcss" "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64"
	chmod +x "./tailwindcss"
	curl -sLO https://github.com/saadeghi/daisyui/releases/latest/download/daisyui.js
	curl -sLO https://github.com/saadeghi/daisyui/releases/latest/download/daisyui-theme.js
	go install github.com/air-verse/air@latest
	go install github.com/a-h/templ/cmd/templ@latest


generate: css
	templ generate

css:  
	npm run build

run: generate
	go run cmd/server/main.go

build: generate
	go build -o ./bin/server cmd/server/main.go

watchcss: generate
	npm run watch

watchtempl: generate
	templ generate --watch

watch: generate
	air --build.cmd "go build -o bin/server cmd/server/main.go" --build.bin "./bin/server"

clean:
	rm -rf bin
	find . -name "*.templ.go" -type f -delete