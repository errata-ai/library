all: build

build:
	# Build our search engine ...
	go build -ldflags="-s -w" -o bin/engine cmd/engine/main.go
	# Create the index ...
	./bin/engine read library.json INDEX
	# Flatten the directory for Netlify ...
	mv INDEX/store/* INDEX
	rm -rf INDEX/store
	zip -r INDEX.zip INDEX -x "*.DS_Store"
	rm -rf INDEX

search:
	rm -rf bin/INDEX
	go build -ldflags="-s -w" -o bin/engine cmd/engine/main.go
	./bin/engine read library.json bin/INDEX
	./bin/engine search ./bin/INDEX 'type:video'