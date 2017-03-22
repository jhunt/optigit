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

deploy:
	@make assets optigit
	cf push

.PHONY: all optigit assets clean
