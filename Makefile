.PHONY: build install test clean

build:
	cd go && go build -o ../jqmas-core ./cmd/jqmas-core

install: build
	cp jqmas-core $(PREFIX)/bin/jqmas-core

test:
	pytest -v

clean:
	rm -f jqmas-core
