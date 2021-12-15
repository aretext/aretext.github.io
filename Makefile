all: build

clean:
	git clean -xfd site/docs

build: clean
	mkdir -p ./site/docs
	mkdocs build

server:
	python -m http.server 8080 --directory ./site
