.PHONY: build build-lambda zip

test:
	go test -count=1 ./...

build:
	cd cmd/standalone && go build

build-lambda:
	cd cmd/lambda && GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath

zip: build-lambda
	zip -j cmd/lambda/lambda.zip cmd/lambda/lambda
