.PHONY: run build clean generate

dependencies:
	curl -L -o "./tailwindcss" "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64"
	chmod +x "./tailwindcss"
	curl -sLO https://github.com/saadeghi/daisyui/releases/latest/download/daisyui.js
	curl -sLO https://github.com/saadeghi/daisyui/releases/latest/download/daisyui-theme.js
	go install github.com/air-verse/air@latest
	go install github.com/a-h/templ/cmd/templ@latest


generate:
	templ generate

css:  
	npm run build

run: generate
	go run cmd/server/main.go

build: generate
	go build -o ./bin/server cmd/server/main.go

watch: 
	@trap 'kill $$(jobs -p) 2>/dev/null; exit' INT TERM; \
	templ generate --watch & \
	npm run watch & \
	air --build.cmd "go build -o bin/server cmd/server/main.go" --build.bin "./bin/server" & \
	wait

clean:
	rm -rf bin
	find . -name "*.templ.go" -type f -delete