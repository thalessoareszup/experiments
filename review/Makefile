.PHONY: build clean serve install

build:
	go build -o review ./cmd/review

clean:
	rm -f review

serve: build
	./review serve --review examples/review.json

install: build
	mkdir -p ~/.local/bin
	cp review ~/.local/bin/review
