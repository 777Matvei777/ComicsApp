build:
	go build ./cmd/xkcd
bench:
	go test -bench=. 
server:
	./xkcd -c config.yaml

test:
	go test -coverprofile=coverage.out  ./...
	echo "generating html"
	go tool cover -html="coverage.out" 

lint:
	golangci-lint run --fix --tests  ./...

sec:
	trivy fs .
	govulncheck ./...
e2e:
	@./e2e_test.sh

web:
	go run web-server/main.go