all: build

clean:
	git clean -xfd site/docs

build: clean
	mkdir -p ./site/docs
	uv run -- mkdocs build

server:
	uv run -- python -m http.server 8080 --directory ./site
