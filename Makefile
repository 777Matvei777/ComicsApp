build:
	go build ./cmd/xkcd
test:
	go test -v
bench:
	go test -bench=. 
server:
	./xkcd -c config.yaml

