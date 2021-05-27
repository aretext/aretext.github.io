all: build

clean:
	rm -rf site

build: clean
	go run build.go --docsPath ./aretext/docs --outputDirPath site

server:
	go run server.go site

fmt:
	goimports -w *.go
