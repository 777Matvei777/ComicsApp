build:
	go build ./cmd/xkcd
bench:
	go test -bench=. 
server:
	./xkcd -c config.yaml

test:
	go test -cover ./...


lint:
	golangci-lint run --fix --tests  ./...

sec:
	trivy fs .
	govulncheck ./...
e2e:
	@./e2e_test.sh