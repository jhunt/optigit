all: optigit
optigit:
	go build -o optigit .

assets:
	go build ./util/embed
	./embed static/assets.go assets/

clean:
	rm -f embed optigit
	rm -f static/assets.go

.PHONY: all optigit assets clean
