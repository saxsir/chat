.PHONY: build

build: 
	mkdir -p builds
	go build -o builds/chat

clean:
	rm -rf builds
