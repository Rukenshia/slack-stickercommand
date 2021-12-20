.PHONY: build clean deploy gomodgen

build:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/lambda lambda/main.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

upload-assets:
	aws s3 sync --delete assets/lpe s3://in.fkn.space/i/lpe/
