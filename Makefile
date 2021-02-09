all: optigit
optigit:
	go fmt ./...
	go build -o optigit .

assets:
	go build ./util/embed
	./embed static/assets.go assets/

clean:
	rm -f embed optigit
	rm -f static/assets.go

build:
	docker build -t filefrog/optigit:latest .
push: build
	docker push filefrog/optigit:latest

.PHONY: all optigit assets clean build push
