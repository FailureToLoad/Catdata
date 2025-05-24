.PHONY: run build clean generate

dependencies:
	curl -L -o "./tailwindcss" "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64"
	chmod +x "./tailwindcss"
	
	curl -sLO "https://github.com/a-h/templ/releases/latest/download/templ_Linux_x86_64.tar.gz"
	tar -xzf "templ_Linux_x86_64.tar.gz"
	rm "templ_Linux_x86_64.tar.gz"
	chmod +x "templ"

	curl -sLO https://github.com/saadeghi/daisyui/releases/latest/download/daisyui.js
	curl -sLO https://github.com/saadeghi/daisyui/releases/latest/download/daisyui-theme.js

	go install github.com/air-verse/air@latest

templ:
	$(error templ missing, run make dependencies)

tailwindcss:
	$(error tailwind missing, run make dependencies)

daisyui.js:
	$(error daisyui missing, run make dependencies)

daisyui-theme.js:
	$(error daisyui-theme missing, run make dependencies)

generate: templ
	./templ generate

css: tailwindcss daisyui.js daisyui-theme.js
	./tailwindcss -i ./input.css -o ./static/styles.css --minify

run: generate
	go run cmd/server/main.go

build: generate
	go build -o ./bin/server cmd/server/main.go

watch: templ tailwindcss daisyui.js daisyui-theme.js
	@trap 'kill $$(jobs -p) 2>/dev/null; exit' INT TERM; \
	./templ generate --watch & \
	./tailwindcss -i ./input.css -o ./static/styles.css --watch & \
	air --build.cmd "go build -o bin/api cmd/server/main.go" --build.bin "./bin/api" & \
	wait

clean:
	rm -rf bin
	find . -name "*.templ.go" -type f -delete